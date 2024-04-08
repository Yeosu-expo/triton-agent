package tritonController

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ahr-i/triton-agent/modelStoreCommunicator"
	"github.com/ahr-i/triton-agent/setting"
)

type requestData struct {
	Provider  string `json:"id"`
	ModelName string `json:"model_name"`
	Version   string `json:"version"`
	Address   string `json:"addr"`
}

func SetModel(body []byte) error {
	// torrent로 요청 전달
	url := fmt.Sprintf("http://%s/serving", setting.TorrentUrl)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	var request requestData
	if err := json.Unmarshal(body, &request); err != nil {
		return err
	}

	err = modelStoreCommunicator.GetScript(request.Provider, request.ModelName, request.Version)
	if err != nil {
		return err
	}

	//triton server의 ssh로 접속, script 파일 실행
	err = RunScriptOnTritonServer(request.Provider, request.ModelName, request.Version)
	if err != nil {
		return err
	}

	return nil
	// filePath := fmt.Sprintf("%s/%s", setting.ModelsPath, provider)
	// fileName := fmt.Sprintf("%s@%s<%s>.torrent", provider, model, version)

	// // Creating the provider folder.
	// // If the provider folder already exists, it will not be created.
	// logCtrlr.Log("Create provider folder.")
	// if err := makeFolder(filePath); err != nil {
	// 	return err
	// }

	// // downloading the model from the Model Store.
	// logCtrlr.Log("Request model download to the Model Store.")
	// modelFile, err := modelStoreCommunicator.GetModel(provider, model, version, filename)
	// if err != nil {
	// 	return err
	// }
	// logCtrlr.Log("Successfully completed the model download.")

	// // Create a directory corresponding to the path.
	// if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
	// 	return err
	// }

	// // Create and write to a file.
	// file, err := os.Create(filePath + "/" + fileName)
	// if err != nil {
	// 	return err
	// }
	// file.Write(modelFile)

	// // File Download by Torrent
	// // 토렌트 클라이언트 설정
	// cfg := torrent.NewDefaultClientConfig()
	// cfg.Seed = true // 시딩 활성화

	// // 토렌트 클라이언트 생성
	// cl, err := torrent.NewClient(cfg)
	// if err != nil {
	// 	log.Fatalf("error creating client: %s", err)
	// }
	// defer cl.Close()

	// // 토렌트 파일 추가
	// torrentPath := filePath + "/" + fileName // 여기에 토렌트 파일 경로 입력
	// metaInfo, err := metainfo.LoadFromFile(torrentPath)
	// if err != nil {
	// 	log.Fatalf("error loading torrent file: %s", err)
	// }
	// t, err := cl.AddTorrent(metaInfo)
	// if err != nil {
	// 	log.Fatalf("error adding torrent: %s", err)
	// }

	// <-t.GotInfo() // 토렌트 정보를 받을 때까지 대기
	// t.DownloadAll()

	// // 파일 다운로드 대기
	// if cl.WaitAll() {

	// 	log.Printf("Downloaded %s", t.Name())
	// 	*channel <- fileName

	// }
	// return nil
}
