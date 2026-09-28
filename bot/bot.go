package bot

import (
	"encoding/json"
	"log"

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

func (b *Bot) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(b.Master, nil)
	if err != nil {
		return err
	}

	b.Conn = conn

	return nil
}

func (b *Bot) Register() error {
	return b.Conn.WriteJSON(Message{Payload: "R", From: b.ID})
}

func (b *Bot) Listen() error {
	defer b.Conn.Close()
	log.Printf("[+] Connected to Master as %s", b.ID)

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

		log.Printf("[C2] Received: %s", msg.Payload)
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
