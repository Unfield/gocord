package gateway

import "github.com/Unfield/Gocord/types"

const (
	MESSAGE_CREATE_EVENT = "MESSAGE_CREATE"
)

type MessageCreateEvent struct {
	types.Message
}
