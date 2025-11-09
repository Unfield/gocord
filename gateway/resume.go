package gateway

import "log"

const OpResume = 6

type resumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

func (g *Gateway) Resume() error {
	if g.sessionID == "" || g.lastSeq == nil {
		return ErrNoSession
	}

	payload := Payload{
		Op: OpResume,
		D: resumeData{
			Token:     g.Token,
			SessionID: g.sessionID,
			Seq:       *g.lastSeq,
		},
	}

	log.Printf("[Gateway] Trying to RESUME with seq=%d", *g.lastSeq)
	return g.Send(payload)
}
