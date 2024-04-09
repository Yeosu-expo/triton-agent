package tritonController

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ahr-i/triton-agent/setting"
	"github.com/ahr-i/triton-agent/src/logCtrlr"
	"github.com/ahr-i/triton-agent/tritonCommunicator"
	"golang.org/x/crypto/ssh"
)

func ChangeModelRepository(provider string, model string, version string) error {
	log.Println("========== Model Repository ==========")
	log.Println("Model repository:", provider)

	// // Checks if the currently polling folder is the provider folder.
	// if provider == modelRepository {
	// 	return nil
	// } else {
	// 	modelRepository = provider
	// }

	// Creating the provider folder.
	// If the provider folder already exists, it will not be created.
	modelRepositoryPath := fmt.Sprintf("%s/%s/%s@%s#%s", "/opt/tritonserver/models", provider, provider, model, version)
	makeFolder(modelRepositoryPath)

	/*
		// Check model.
		if !CheckModel(provider, model) {
			logCtrlr.Log("Model is not exist.")
			return errors.New("ah")
			}
	*/

	// This is the SSH configuration for connecting to Triton.
	if err := startTritonServer(modelRepositoryPath); err != nil {
		return err
	}

	// Starts polling the model repository.
	if err := polling(model, version); err != nil {
		return err
	}
	setting.LoadedModel = fmt.Sprintf("%s@%s#%s", provider, model, version)
	return nil
}

func startTritonServer(modelRepositoryPath string) error {
	// This is the SSH configuration for connecting to Triton.
	config := &ssh.ClientConfig{
		User: setting.TritonUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(setting.TritonPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connects to the Triton server.
	logCtrlr.Log("Attach Triton server.")
	client, err := ssh.Dial("tcp", setting.TritonSSH, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Executes the 'pkill' command to terminate the running Triton server.
	logCtrlr.Log("Kill Triton server.")
	if err := executeCommand(client, "pkill -f /opt/tritonserver"); err != nil {
		log.Println(err)
	}

	// Starts the Triton server, specifying the provider folder as the model repository.
	logCtrlr.Log("Change the model repository and start the Triton server.")
	startCommand := fmt.Sprintf("nohup /opt/tritonserver/bin/tritonserver --model-repository %s > /dev/null 2>&1 & exit", modelRepositoryPath)
	//startCommand := fmt.Sprintf("nohup /opt/tritonserver/bin/tritonserver --model-repository %s", modelRepositoryPath)
	if err := executeCommand(client, startCommand); err != nil {
		return err
	}

	return nil
}

func RunScriptOnTritonServer(provider, model, version string) error {
	// This is the SSH configuration for connecting to Triton.
	config := &ssh.ClientConfig{
		User: setting.TritonUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(setting.TritonPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connects to the Triton server.
	logCtrlr.Log("Attach Triton server.")
	client, err := ssh.Dial("tcp", setting.TritonSSH, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// Executes the 'bash .sh' command to terminate the running Triton server.
	fileName := fmt.Sprintf("%s@%s#%s", provider, model, version)
	scriptPath := fmt.Sprintf("/opt/tritonserver/models/%s/%s.sh", provider, fileName)
	command := "chmod 755 " + scriptPath + " && " + scriptPath
	//command := "bash " + scriptPath
	logCtrlr.Log("Run bash Script to Triton server. Command: " + command)
	if err := executeCommand(client, command); err != nil {
		log.Println(err)
	}

	return nil
}

func executeCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()

	if err != nil {
		return err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &b

	if err := session.Run(command); err != nil {
		return fmt.Errorf("failed to run command: %s, output: %s", err, b.String())
	}
	log.Println("* (SYSTEM) Command output:", b.String())

	return nil
}

func polling(model string, version string) error {
	logCtrlr.Log("Polling start - Model repository.")
	cnt := 0

	for {
		if cnt == 20 {
			return errors.New("model not found")
		}

		cnt++
		log.Printf("* (SYSTEM) Checks the model repository of the Triton server. (It is the %dth checking)\n", cnt)
		time.Sleep(5 * time.Second)

		ready, err := tritonCommunicator.Ready(model, version)
		if err != nil {
			logCtrlr.Error(errors.New("triton server is not working"))
			continue
		}

		if ready {
			break
		}
	}

	logCtrlr.Log("Triton Server Started")
	return nil
}
