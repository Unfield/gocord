package gateway

import "github.com/Unfield/gocord/types"

type EventName string

const (
	EventReady         EventName = "READY"
	EventMessageCreate EventName = "MESSAGE_CREATE"
)

type MessageCreateEvent struct {
	types.Message
}

type ReadyEvent struct {
	User types.User `json:"user"`
}
