package framework

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/antonholmquist/jason"
)

var (
	RobotName                                     string
	RobotCount, RetryCount, ReqCount, PkgInterval int
	retryTimes                                    []struct{}
	sceneID                                       int
)

func readFile(filename string) error {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("ReadFile:", err.Error())
		return err
	}

	jsonRet, err := jason.NewObjectFromBytes(bytes)

	if err != nil {
		return err
	}

	RobotName, err = jsonRet.GetString("robotName")
	v, err := jsonRet.GetInt64("robotNum")
	if err == nil {
		RobotCount = int(v)
	}
	v, err = jsonRet.GetInt64("retryNum")
	if err == nil {
		RetryCount = int(v)
	}
	v, err = jsonRet.GetInt64("reqN")
	if err == nil {
		ReqCount = int(v)
	}
	v, err = jsonRet.GetInt64("pkgInterval")
	if err == nil {
		PkgInterval = int(v)
	}
	v, err = jsonRet.GetInt64("sceneId")
	if err == nil {
		sceneID = int(v)
	}

	return nil
}

func InitConfig(configFile string) {
	log.Println("got configile file ", configFile)
	readFile(configFile)
}
