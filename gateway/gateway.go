package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const GatewayURL = "wss://gateway.discord.gg/?v=10&encoding=json"

type Gateway struct {
	conn   *websocket.Conn
	mu     sync.Mutex
	closed bool

	OnEvent func(evt *Payload)

	Token   string
	Intents int
}

func New() *Gateway {
	return &Gateway{}
}

func (g *Gateway) Connect(ctx context.Context) error {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, GatewayURL, nil)
	if err != nil {
		return fmt.Errorf("gateway connect failed: %w", err)
	}

	g.conn = conn
	log.Println("[Gateway] Connected to Discord")

	go g.listen(ctx)
	return nil
}

func (g *Gateway) listen(ctx context.Context) {
	defer g.Close()

	for {
		_, msg, err := g.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Printf("[Gateway] Connection closed: %v", err)
				return
			}
			log.Printf("[Gateway] Read error: %v", err)
			return
		}

		var payload Payload
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Printf("[Gateway] Unmarshal error: %v", err)
			continue
		}

		switch payload.Op {
		case 10:
			g.handleHello(payload.D)
		case 11:
			log.Println("[Gateway] ⇐ Heartbeat ACK")
		default:
			if g.OnEvent != nil {
				g.OnEvent(&payload)
			}
		}
	}
}

func (g *Gateway) Send(payload any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.closed {
		return fmt.Errorf("gateway is closed")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Gateway) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.closed {
		return
	}
	g.closed = true

	if g.conn != nil {
		_ = g.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		_ = g.conn.Close()
	}

	log.Println("[Gateway] Closed")
}
