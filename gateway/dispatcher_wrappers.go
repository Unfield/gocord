package gateway

func (d *TypedDispatcher) OnReady(fn func(ReadyEvent)) {
	On[ReadyEvent](d, EventReady, fn)
}

func (d *TypedDispatcher) OnMessageCreate(fn func(MessageCreateEvent)) {
	On[MessageCreateEvent](d, EventMessageCreate, fn)
}
