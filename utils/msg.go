package utils

import "encoding/json"

var response = map[string]interface{}{
	"code":    "0",
	"message": "",
}

func Msg(code int, message string) string {
	response["code"] = code
	response["message"] = message
	responseJson, _ := json.Marshal(response)
	return string(responseJson)
}
