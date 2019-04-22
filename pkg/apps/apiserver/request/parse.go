package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseBody(r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
