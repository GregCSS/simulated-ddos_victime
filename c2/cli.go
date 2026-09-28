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
		Short: "DDoS Simulation: Bot Master CLI",
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
		Short: "Broadcast a payload to all connected bots",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.broadcast(args[0])
		},
	}

	rootCmd.AddCommand(listBotsCmd, broadcastCmd)
	return rootCmd
}

func RunCLI(m *Master) error {
	rootCmd := NewCLI(m)
	time.Sleep(25 * time.Millisecond) // Just to ease my OCD because HTTP serve log was messing up prints

	fmt.Println("DDoS Simulation: Bot Master CLI")
	fmt.Println("Type 'help' for available commands.")
	fmt.Println()

	// Interactive Console
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("c2> ")
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
