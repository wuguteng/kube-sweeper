package unittests

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/kevinyjn/gocom/testingutil"
	"wuguteng.com/kube-sweeper/sweeper"
)

func TestMatchCleanLogFileName(t *testing.T) {
	logPath := "./tests"
	kubeAPIServer := ""
	kubeAPICaPath := ""
	kubeAPIToken := ""
	r, err := regexp.Compile(sweeper.LogPathExtractPattern)
	testingutil.AssertNil(t, err, "regexp.Compile err")
	testingutil.AssertNotNil(t, r, "regexp.Compile regex")
	r2, err := regexp.Compile(sweeper.LogPathPodNameExtractPattern)
	testingutil.AssertNil(t, err, "regexp.Compile err")
	testingutil.AssertNotNil(t, r2, "regexp.Compile regex")
	fis, err := ioutil.ReadDir(logPath)
	testingutil.AssertNil(t, err, "ioutil.ReadDir err")
	testingutil.AssertNotNil(t, fis, "ioutil.ReadDir files")
	for _, fi := range fis {
		if fi.IsDir() {
			podInfo, err := sweeper.ExtractPodInfoFromFileName(r, r2, fi.Name())
			testingutil.AssertNil(t, err, "ExtractPodInfoFromFileName err")
			fmt.Printf("match result for %s is: %+v\n", fi.Name(), podInfo)
			if "" != kubeAPIServer {
				status, err := sweeper.GetKubernetesPodState(podInfo, kubeAPIServer, kubeAPICaPath, kubeAPIToken)
				testingutil.AssertNil(t, err, "GetKubernetesPodState err")
				fmt.Printf("status of pod:%s is:%+v\n", podInfo.PodName, status)
			}
		}
	}
}
