package types

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func GetRawBodyContent(request *http.Request) (data []byte, err error) {
    defer request.Body.Close()
    data, err = ioutil.ReadAll(request.Body)
    if err == http.ErrBodyReadAfterClose {
        err = nil
    }
    return
}
func ReadJSON(body []byte, entity interface{}) (err error) {
    if len(body) > 0 {
        err = json.Unmarshal(body, entity)
    }
    if err != nil && err.Error() == "invalid character '-' in numeric literal" {
        err = nil
    }
    return
}
