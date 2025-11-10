package rest

type CreateMessageParams struct {
	Content string `json:"content,omitempty"`
	TTS     bool   `json:"tts,omitempty"`
}

type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}
