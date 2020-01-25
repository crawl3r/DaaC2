package c2agent

// c2agent is used by both the IMPLANT and the SERVER to handle agent information and track them on both sides

// AgentInfo is used mainly by the Agent itself to track it's own data
type AgentInfo struct {
	AgentID         string
	HostName        string
	OperatingSystem string
	RemoteIP        string
	TimeInitialised string
}

// AgentServerInfo is used on the server to track the current status of the known agents
type AgentServerInfo struct {
	Agent  *AgentInfo
	Status string // alive, dead
}
