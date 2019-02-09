package registryserver

import (
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/notifications"
	"io/ioutil"
	"net/http"
)

func Hook(response http.ResponseWriter, request *http.Request) {
	bytes, e := ioutil.ReadAll(request.Body)
	if e != nil {
		event := &notifications.Event{}
		err := json.Unmarshal(bytes, event)
		if err != nil {
			fmt.Println("HOOK ----------------------")
			fmt.Println("", event.Action, string(bytes))
			fmt.Println("END  ----------------------\n")
		}

	}
	_, _ = response.Write([]byte("ok"))
}
