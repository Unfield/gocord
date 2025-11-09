package gateway

import (
	"log"
	"runtime"
)

const OpIdentify = 2

type identifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type identifyData struct {
	Token      string             `json:"token"`
	Intents    int                `json:"intents"`
	Properties identifyProperties `json:"properties"`
}

func (g *Gateway) Identify(token string, intents int) error {
	payload := Payload{
		Op: OpIdentify,
		D: identifyData{
			Token:   token,
			Intents: intents,
			Properties: identifyProperties{
				OS:      runtime.GOOS,
				Browser: "gocord",
				Device:  "gocord",
			},
		},
	}

	log.Println("[Gateway] Sending Identify")
	return g.Send(payload)
}
