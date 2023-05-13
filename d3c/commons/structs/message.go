package structs

type Message struct {
	SessionId     string
	AgentId       string
	AgentHostName string
	AgentCWD      string
	Commands      []Command
}
