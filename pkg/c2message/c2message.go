package c2message

import (
	"encoding/base64"
	"encoding/json"
)

// Message types that we can if/else on when we are reading a message
const MESSAGE_AGENT_JOIN = "new"
const MESSAGE_AGENT_LEFT = "left"
const MESSAGE_COMMAND = "cmd"
const MESSAGE_SHELLCODE = "shellcode"
const MESSAGE_AGENTS = "agents"
const MESSAGE_RESPONSE = "reply"
const MESSAGE_PING = "ping"
const MESSAGE_KILL = "kill"

// Message is used across the network as the object to hold all the data
type Message struct {
	FromServer  bool   // flip this bool to true if it's from the server
	ToID        string // This will be set to the target AgentID if we require specific communication
	FromID      string // This will be set as the AgentID that has sent the data
	MessageType string // This will be one of the above strings to quickly determine what the message actually is
	Data        string // This will be the actual contents of our message, i.e the actual command or the response
}

// CreateNewMessage builds our new message object and returns a pointer to this obj
func CreateNewMessage(messageType string, fromID string, toID string, dataText string) *Message {
	newMessageStruct := &Message{}
	newMessageStruct.MessageType = messageType
	newMessageStruct.FromID = fromID
	newMessageStruct.ToID = toID
	newMessageStruct.FromServer = fromID == "SERVER"
	newMessageStruct.Data = dataText
	return newMessageStruct
}

// EncodeMessageObject takes a struct, json's it and then encodes in Base64, returning the b64 string
func EncodeMessageObject(d *Message) string {
	jsonData, _ := json.Marshal(d)
	data := []byte(string(jsonData)) // string reparesentation of the marshalled struct above
	return base64.StdEncoding.EncodeToString(data)
}
