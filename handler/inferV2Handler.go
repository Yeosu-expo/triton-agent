package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ahr-i/triton-agent/schedulerCommunicator/healthPinger"
	"github.com/ahr-i/triton-agent/setting"
	"github.com/ahr-i/triton-agent/src/httpController"
	"github.com/ahr-i/triton-agent/src/logCtrlr"
	"github.com/ahr-i/triton-agent/tritonCommunicator"
	"github.com/ahr-i/triton-agent/tritonController"
	"github.com/gorilla/mux"
)

var mutex sync.Mutex

func (h *Handler) inferV2Handler(w http.ResponseWriter, r *http.Request) {
	// Extract model information from the URL.
	vars := mux.Vars(r)
	provider := vars["provider"]
	model := vars["model"]
	version := vars["version"]

	var burstTime float64
	healthPinger.UpdateTaskInfo_start(provider, model, version)
	defer func() {
		//callback.Callback(burstTime, provider, model, version)
		healthPinger.UpdateTaskInfo_end(provider, model, version, burstTime)
	}()
	mutex.Lock()
	defer mutex.Unlock()

	// Extract the request from the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logCtrlr.Error(err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}
	printModelInfo(provider, model, version, string(body))

	// Check triton ready
	ready, err := tritonCommunicator.Ready(model, version)
	if err != nil {
		logCtrlr.Error(errors.New("triton server is not working"))
	}
	// Set model repository
	if !ready {
		if err := tritonController.ChangeModelRepository(provider, model, version); err != nil {
			logCtrlr.Error(err)
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}
	}

	// Request to tritons
	log.Println("Inference 요청 보내기!!!!")
	startTime := time.Now()
	response, err := tritonCommunicator.Inference(model, version, body)
	if err != nil {
		logCtrlr.Error(err)
		rend.JSON(w, http.StatusBadRequest, nil)
		return
	}
	endTime := time.Now()

	log.Println(string(response))

	// Send burst time to scheduler
	burstTime = float64(endTime.Sub(startTime).Milliseconds()) / 1000
	log.Printf("* (SYSTEM) Burst time: %f\n", burstTime)

	httpController.JSON(w, http.StatusOK, response)
	//httpController.JSON(w, http.StatusOK, nil)

}

func printModelInfo(provider string, model string, version string, request string) {
	logCtrlr.Log("Request: ▽▽▽▽▽▽▽▽▽▽")
	log.Printf("%19s: %s", "Provider", provider)
	log.Printf("%19s: %s", "Model:", model)
	log.Printf("%19s: %s", "Version:", version)
	//log.Printf("%28s : %s", "Request:", request)

	if setting.LoadedModel != "" {
		atSplit := strings.Split(setting.LoadedModel, "@")
		if len(atSplit) != 2 {
			fmt.Println("Error: Input does not contain a valid '@' separator")
			return
		}
		Lprovider := atSplit[0]

		// '#' 기호를 기준으로 두 번째 분할
		hashSplit := strings.Split(atSplit[1], "#")
		if len(hashSplit) != 2 {
			fmt.Println("Error: Input does not contain a valid '#' separator")
			return
		}
		Lmodel := hashSplit[0]
		Lversion := hashSplit[1]

		logCtrlr.Log("Loaded Model: ▽▽▽▽▽▽▽▽▽▽")
		log.Printf("%23s: %s", "Provider", Lprovider)
		log.Printf("%23s: %s", "Model:", Lmodel)
		log.Printf("%23s: %s", "Version:", Lversion)
	}

	logCtrlr.Log("Not Load Model: ▽▽▽▽▽▽▽▽▽▽")
	notLoadedCnt := 1
	for modelkey, value := range healthPinger.Model_info {
		for version := range value {
			log.Printf("%26d: %s#%s\n", notLoadedCnt, modelkey, version)
			notLoadedCnt++
		}
	}
}

func (h *Handler) testInferV2Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]
	model := vars["model"]
	version := vars["version"]

	log.Println("provider :", provider)
	log.Println("model :", model)
	log.Println("version :", version)

	var burstTime float64
	healthPinger.UpdateTaskInfo_start(provider, model, version)
	defer func() {
		//callback.Callback(burstTime, provider, model, version)
		healthPinger.UpdateTaskInfo_end(provider, model, version, burstTime)
	}()

	mutex.Lock()
	defer mutex.Unlock()

	// Check triton ready
	ready, err := tritonCommunicator.Ready(model, version)
	if err != nil {
		logCtrlr.Error(errors.New("triton server is not working"))
	}
	// Set model repository
	if !ready {
		if err := tritonController.ChangeModelRepository(provider, model, version); err != nil {
			logCtrlr.Error(err)
			rend.JSON(w, http.StatusBadRequest, nil)
			return
		}
	}

	//요청받으면 랜덤한 인퍼런스 타임으로 결과값 돌려주기.
	inferTime := getRandNum(500, 10000)
	time.Sleep(time.Millisecond * time.Duration(inferTime))

	//랜덤한 확률로 인퍼런스 중 fault상황
	randRate := rand.Float64()
	if randRate < 0.1 {
		log.Println("fault 상황, 종료")
		os.Exit(1)
	}

	//정상 수행 후 응답 상황
	burstTime = float64(inferTime) / 1000
	httpController.JSON(w, http.StatusOK, nil)
}

func getRandNum(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
