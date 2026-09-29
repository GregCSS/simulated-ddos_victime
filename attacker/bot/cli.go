package bot

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	return net.IP(conn.LocalAddr().(*net.UDPAddr).IP).String()
}

func NewCLI() *cobra.Command {
	var (
		id     string
		master string
	)

	fmt.Println(`
$$$$$$$\   $$$$$$\ $$$$$$$$\ 
$$  __$$\ $$  __$$\\__$$  __|
$$ |  $$ |$$ /  $$ |  $$ |   
$$$$$$$\ |$$ |  $$ |  $$ |   
$$  __$$\ $$ |  $$ |  $$ |   
$$ |  $$ |$$ |  $$ |  $$ |   
$$$$$$$  | $$$$$$  |  $$ |   
\_______/  \______/   \__|   
	`)

	fmt.Println("DDoS Simulation: BOT CLI\n")

	rootCmd := &cobra.Command{
		Use:   "bot",
		Short: "DDoS Simulation: Bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			bot := NewBot(id, master)

			if err := bot.Run(); err != nil {
				return fmt.Errorf("bot stopped: %w", err)
			}

			return nil
		},
	}

	// Setting flags which creates bot automatically
	rootCmd.Flags().StringVarP(&id, "id", "i", getOutboundIP(), "Bot ID")
	rootCmd.Flags().StringVarP(&master, "master", "m", "ws://localhost:8080/connect", "Botmaster WebSocket URL")

	return rootCmd
}
