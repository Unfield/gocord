package gateway

import (
	"encoding/json"
	"log"
	"time"
)

type helloData struct {
	HeartbeatInterval float64 `json:"heartbeat_interval"`
}

func (g *Gateway) handleHello(raw any) {
	dataMap, ok := raw.(map[string]any)
	if !ok {
		log.Println("[Gateway] Invalid Hello data")
		return
	}

	var hello helloData
	bytes, _ := json.Marshal(dataMap)
	if err := json.Unmarshal(bytes, &hello); err != nil {
		log.Printf("[Gateway] Failed to decode Hello payload: %v", err)
		return
	}

	interval := time.Duration(hello.HeartbeatInterval) * time.Millisecond
	log.Printf("[Gateway] Starting heartbeat every %v", interval)

	go g.startHeartbeat(interval)

	if g.Token != "" {
		go func() {
			time.Sleep(1 * time.Second)
			if err := g.Identify(g.Token, g.Intents); err != nil {
				log.Printf("[Gateway] Identify error: %v", err)
			}
		}()
	}
}

func (g *Gateway) startHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		packet := Payload{
			Op: 1,
			D:  nil,
		}

		if err := g.Send(packet); err != nil {
			log.Printf("[Gateway] Heartbeat send error: %v", err)
			return
		}

		log.Println("[Gateway] ⇒ Sent heartbeat")
	}
}
