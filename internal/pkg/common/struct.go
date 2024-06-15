package commonx

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

func StructToJson(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func JsonToStruct(str string, value interface{}) error {
	return json.Unmarshal([]byte(str), value)
}

