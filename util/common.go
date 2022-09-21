package util

import (
	"encoding/json"
	"github.com/golang/glog"
)

func StructToMap(value interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	resJson, err := json.Marshal(value)
	if err != nil {
		glog.Errorf("Json Marshal failed ,msg: %s", err.Error())
		return nil
	}
	err = json.Unmarshal(resJson, &m)
	if err != nil {
		glog.Errorf("Json Unmarshal failed,msg : %s", err.Error())
		return nil
	}
	return m
}
