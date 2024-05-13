package healthPinger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ahr-i/triton-agent/setting"
	"github.com/ahr-i/triton-agent/src/logCtrlr"
)

type RequestData struct {
	Port       string                         `json:"port"`
	Gpuname    string                         `json:"gpuname"`
	Model_info map[string]map[string]TaskInfo `json:"model_info"`
}

func postAlive(client *http.Client) {
	jsonData, err := json.Marshal(RequestData{
		Port:       port,
		Gpuname:    gpuName,
		Model_info: Model_info,
	})
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("http://%s/alive", setting.ManagerUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("<NewRequest>", err)
	}

	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("<client.Do>", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logCtrlr.Error(errors.New("there is no manager"))

		//os.Exit(1)
	}
}
