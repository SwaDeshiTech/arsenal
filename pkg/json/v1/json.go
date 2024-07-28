package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ParseJSON(jsonObject []byte, object interface{}) error {
	err := json.Unmarshal([]byte(jsonObject), &object)
	if err != nil {
		log.Println("failed to parse object", err)
		return err
	}
	return nil
}

func ConvertIntoJSON(object interface{}) ([]byte, error) {
	jsonResp, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}

func ReadJSONFile(filePath string) ([]byte, error) {
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("error reading json file", err)
		return nil, err
	}
	return jsonFile, nil
}
