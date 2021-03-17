package sweeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/kevinyjn/gocom/definations"
	"github.com/kevinyjn/gocom/httpclient"
	"github.com/kevinyjn/gocom/logger"
	"github.com/kevinyjn/gocom/utils"
	"wuguteng.com/kube-sweeper/models"
)

type fileLogCleaner struct {
	logpath             string
	kubernetesApiServer string
	kubernetesToken     string
	kubernetesCaPath    string
	timer               *utils.Timer
	pathpattern         *regexp.Regexp
	podnamepattern      *regexp.Regexp
}

// Constants
const (
	LogPathExtractPattern        = `(?P<lan>\w+)-(?P<namespace>\w+)-(?P<pod>\w+)-(?P<version>[v]?\d+(?:\.\d+)+(?:-(?:alpha|beta|rc|hotfix)\d*)?)-(?P<pod_suffix>.*)$`
	LogPathPodNameExtractPattern = `\w+-\w+-(.*)`
	KubernetesCAFilePath         = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	KubernetesTokenPath          = "/var/run/secrets/kubernetes.io/serviceaccount/token"

	KubernetesPodAPI = "/api/v1/namespaces/%s/pods/%s"
)

// Variables
var (
	c = fileLogCleaner{}
)

// StartFileLogCleaner cleans kubernetes app file log
func StartFileLogCleaner(filePath string) error {
	return c.start(filePath)
}

func (c *fileLogCleaner) start(filePath string) error {
	if "" == filePath {
		logger.Error.Printf("start file log cleaner with log path:%s were not valid", filePath)
		return fmt.Errorf("file log path:%s were not valid", filePath)
	}

	r, err := regexp.Compile(LogPathExtractPattern)
	if nil != err {
		logger.Error.Printf("compile log path extract full info pattern:%s failed with error:%v", LogPathExtractPattern, err)
		return err
	}
	r2, err := regexp.Compile(LogPathPodNameExtractPattern)
	if nil != err {
		logger.Error.Printf("compile log path extract pod name pattern:%s failed with error:%v", LogPathPodNameExtractPattern, err)
		return err
	}

	c.logpath = filePath
	c.pathpattern = r
	c.podnamepattern = r2
	c.kubernetesApiServer = "https://kubernetes.default"
	c.kubernetesCaPath = KubernetesCAFilePath

	if "" != os.Getenv("KUBERNETES_SERVICE_HOST") && "" != os.Getenv("KUBERNETES_SERVICE_PORT") {
		c.kubernetesApiServer = fmt.Sprintf("https://%s:%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"))
	}

	token, err := ioutil.ReadFile(KubernetesTokenPath)
	if nil != err {
		logger.Error.Printf("read kubernetes token:%s failed with error:%v", KubernetesTokenPath, err)
	} else {
		c.kubernetesToken = string(token)
	}

	if nil != c.timer {
		c.timer.Stop()
	}

	c.timer, err = utils.NewTimer(20000, 300000, c.doClean, &c)
	if nil != err {
		logger.Error.Printf("start file log cleaner with log path:%s failed with error:%v", filePath, err)
	}
	return err
}

func (c *fileLogCleaner) doClean(t *utils.Timer, tim time.Time, delegate interface{}) {
	fis, err := ioutil.ReadDir(c.logpath)
	if nil != err {
		logger.Error.Printf("list files in directory:%s failed with error:%v", c.logpath, err)
		return
	}
	for _, fi := range fis {
		if fi.IsDir() {
			podInfo, err := ExtractPodInfoFromFileName(c.pathpattern, c.podnamepattern, fi.Name())
			if nil != err {
				continue
			}

			// check pod exists
			podStatus, err := GetKubernetesPodState(podInfo, c.kubernetesApiServer, c.kubernetesCaPath, c.kubernetesToken)
			if "NotFound" == podStatus.GetReason() {
				// remove the path
				logger.Info.Printf("found that the pod belongs to log path:%s were not exists, clean it.", fi.Name())
				rp := path.Join(c.logpath, fi.Name())
				err = os.RemoveAll(rp)
				if nil != err {
					logger.Error.Printf("removing path:%s failed with error:%v", rp, err)
				}
			}
		}
	}
	// force run gc to recycle the memory usage
	runtime.GC()
}

// ExtractPodInfoFromFileName from path name
func ExtractPodInfoFromFileName(pathpattern *regexp.Regexp, podnamepattern *regexp.Regexp, fileName string) (models.PodInfo, error) {
	mr2 := podnamepattern.FindAllStringSubmatch(fileName, -1)
	podInfo := models.PodInfo{Labels: map[string]string{}}
	if len(mr2) < 1 {
		return podInfo, fmt.Errorf("path:%s not valid for pod path name", fileName)
	}
	if len(mr2[0]) > 1 {
		podInfo.PodName = mr2[0][1]
	}
	mr := pathpattern.FindAllStringSubmatch(fileName, -1)
	if len(mr) < 1 {
		return podInfo, fmt.Errorf("path:%s not valid for pod path name", fileName)
	}
	parts := mr[0]
	pl := len(parts)
	if pl > 1 {
		podInfo.Labels["lan"] = parts[1]
	}
	if pl > 2 {
		podInfo.Namespace = parts[2]
	}
	if pl > 3 {
		podInfo.Name = parts[3]
	}
	if pl > 4 {
		podInfo.Version = parts[4]
	}
	if pl > 5 {
		podInfo.PosSuffix = parts[5]
	}
	return podInfo, nil
}

// GetKubernetesPodState by kube api
func GetKubernetesPodState(podInfo models.PodInfo, apiServer string, caPath string, token string) (models.KubeAPIStatus, error) {
	var result models.KubeAPIStatus
	tlsOption := definations.TLSOptions{
		Enabled: true,
		CaFile:  caPath,
	}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}
	api := fmt.Sprintf(KubernetesPodAPI, podInfo.Namespace, podInfo.PodName)
	resp, err := httpclient.HTTPQuery("GET", apiServer+api, nil, httpclient.WithHTTPTLSOptions(&tlsOption), httpclient.WithHTTPHeaders(headers))
	if nil != err {
		logger.Warning.Printf("query api:%s failed with error:%+v", api, err)
		result = &models.KubeFailureResult{}
	} else {
		result = &models.KubeAPIResult{}
	}
	if nil != resp {
		err = json.Unmarshal(resp, result)
		if nil != err {
			logger.Error.Printf("parse kube api:%s result:%s failed with error:%v", api, string(resp), err)
		}
	}
	return result, err
}
