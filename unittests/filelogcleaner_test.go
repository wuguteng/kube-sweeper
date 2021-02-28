package unittests

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/kevinyjn/gocom/unittests"
	"wuguteng.com/kube-sweeper/sweeper"
)

func TestMatchCleanLogFileName(t *testing.T) {
	logPath := "./tests"
	kubeAPIServer := ""
	kubeAPICaPath := ""
	kubeAPIToken := ""
	r, err := regexp.Compile(sweeper.LogPathExtractPattern)
	unittests.AssertNil(t, err, "regexp.Compile err")
	unittests.AssertNotNil(t, r, "regexp.Compile regex")
	r2, err := regexp.Compile(sweeper.LogPathPodNameExtractPattern)
	unittests.AssertNil(t, err, "regexp.Compile err")
	unittests.AssertNotNil(t, r2, "regexp.Compile regex")
	fis, err := ioutil.ReadDir(logPath)
	unittests.AssertNil(t, err, "ioutil.ReadDir err")
	unittests.AssertNotNil(t, fis, "ioutil.ReadDir files")
	for _, fi := range fis {
		if fi.IsDir() {
			podInfo, err := sweeper.ExtractPodInfoFromFileName(r, r2, fi.Name())
			unittests.AssertNil(t, err, "ExtractPodInfoFromFileName err")
			fmt.Printf("match result for %s is: %+v\n", fi.Name(), podInfo)
			if "" != kubeAPIServer {
				status, err := sweeper.GetKubernetesPodState(podInfo, kubeAPIServer, kubeAPICaPath, kubeAPIToken)
				unittests.AssertNil(t, err, "GetKubernetesPodState err")
				fmt.Printf("status of pod:%s is:%+v\n", podInfo.PodName, status)
			}
		}
	}
}
