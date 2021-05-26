package common

import "encoding/json"

func DataToString(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func StringToData(s string, data interface{}) error {
	return json.Unmarshal([]byte(s), data)
}