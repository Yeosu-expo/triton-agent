package healthPinger

import (
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/ahr-i/triton-agent/setting"
)

var port string
var gpuName string

// var model_info map[string]map[string]TaskInfo = make(map[string]map[string]TaskInfo)
var Model_info map[string]map[string]TaskInfo

type TaskInfo struct {
	LoadedAmount         int     `json:"loaded_amount"`
	AverageInferenceTime float32 `json:"average_inference_time"`
}

func Enter() {
	Model_info = make(map[string]map[string]TaskInfo)

	port = setting.ServerPort

	cmd := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("failed to get GPU name: ", err)
	}

	gpuName = strings.TrimSpace(string(output))

	log.Println("GPU NAME :", gpuName)

	alivePoster()
}

func alivePoster() {
	var cnt int = 0

	log.Println("tcp 오픈")
	ln, err := net.Listen("tcp", ":6934")
	if err != nil {
		log.Fatal("헬스체커용 tcp 오픈 실패", err)
	}

	go func() {
		log.Println("승인 중")
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Acppea", err)
		}

		log.Println("헬스체킹용 tcp 연결 성공", conn.RemoteAddr().String())
		select {}
	}()

	transport := &http.Transport{
		MaxIdleConns:        2,
		MaxIdleConnsPerHost: 2,
		DisableKeepAlives:   false,
	}
	client := &http.Client{
		Transport: transport,
	}

	for {
		cnt++
		log.Printf("* (System) Send information to the Manager. (It is the %dth request)\n", cnt)

		postAlive(client)

		time.Sleep(1 * time.Second)
	}
}
