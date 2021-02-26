package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/kataras/iris"
	"github.com/kevinyjn/gocom/application"
	"github.com/kevinyjn/gocom/config"
	"github.com/kevinyjn/gocom/healthz"
	"github.com/kevinyjn/gocom/logger"
	"wuguteng.com/kube-sweeper/sweeper"
)

type serviceStarter func() error

var NonConfiguredError = fmt.Errorf("Not configured")

func main() {
	succeeds := 0
	startModule(startFileLogCleaner, "file log cleaner", &succeeds)

	if 0 >= succeeds {
		logger.Error.Printf("non module started.")
		return
	}

	startServer()
}

func startModule(module serviceStarter, moduleName string, succeeds *int) {
	err := module()
	if nil == err {
		if nil != succeeds {
			*succeeds = *succeeds + 1
		}
		logger.Info.Printf("start module %s succeed.", moduleName)
	} else {
		if err != NonConfiguredError {
			logger.Error.Printf("start module %s failed with error:%v", moduleName, err)
		}
	}
}

func startFileLogCleaner() error {
	logPath := os.Getenv("KUBE_FILE_LOG_PATH")
	if "" == logPath {
		logPath = "/var/log/k8sapps"
	}
	return sweeper.StartFileLogCleaner(logPath)
}

func startServer() {
	serverPort := os.Getenv("LISTEN_PORT")
	bindPort := 80
	if "" != serverPort {
		bindPort, err := strconv.Atoi(serverPort)
		if nil != err {
			logger.Error.Printf("server port:%s were not valid:%v", serverPort, err)
			return
		} else if 0 >= bindPort {
			logger.Error.Printf("server port:%s were not valid", serverPort)
			return
		}
	}

	listenAddr := fmt.Sprintf("%s:%d", "0.0.0.0", bindPort)
	logger.Info.Printf("starting server on %s...", listenAddr)
	app := application.GetApplication(config.ServerMain)
	healthz.InitHealthz(app)
	app.Run(iris.Addr(listenAddr), iris.WithoutServerError(iris.ErrServerClosed))
}
