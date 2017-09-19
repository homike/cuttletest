package framework

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type RobotConfig struct {
	RobotCount  int `xml:"robot_num"`
	RetryCount  int `xml:"retry_num"`
	ReqCount    int `xml:"req_num"`
	PkgInterval int `xml:"pkginterval"`
	SceneID     int `xml:"scene_id"`
}

var RobotCfg *RobotConfig

func readFile(filename string) (*RobotConfig, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("ReadFile:", err.Error())
		return nil, err
	}

	v := RobotConfig{}
	err = xml.Unmarshal(bytes, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}

	return &v, nil
}

func InitConfig(configFile string) {
	var err error
	RobotCfg, err = readFile(configFile)
	fmt.Println(RobotCfg.PkgInterval, RobotCfg.ReqCount, RobotCfg.RetryCount, RobotCfg.RobotCount, RobotCfg.SceneID)
	_ = err
}

// var (
// 	RobotCount, RetryCount, ReqCount, PkgInterval int
// 	//retryTimes                                    []struct{}
// 	sceneID int
// )

// func readFile(filename string) error {
// 	bytes, err := ioutil.ReadFile(filename)

// 	if err != nil {
// 		fmt.Println("ReadFile:", err.Error())
// 		return err
// 	}

// 	jsonRet, err := jason.NewObjectFromBytes(bytes)

// 	if err != nil {
// 		return err
// 	}

// 	//RobotName, err = jsonRet.GetString("robotName")
// 	v, err := jsonRet.GetInt64("robotNum")
// 	if err == nil {
// 		RobotCount = int(v)
// 	}
// 	v, err = jsonRet.GetInt64("retryNum")
// 	if err == nil {
// 		RetryCount = int(v)
// 	}
// 	v, err = jsonRet.GetInt64("reqN")
// 	if err == nil {
// 		ReqCount = int(v)
// 	}
// 	v, err = jsonRet.GetInt64("pkgInterval")
// 	if err == nil {
// 		PkgInterval = int(v)
// 	}
// 	v, err = jsonRet.GetInt64("sceneId")
// 	if err == nil {
// 		sceneID = int(v)
// 	}

// 	return nil
// }

// func InitConfig(configFile string) {
// 	readFile(configFile)
// }
