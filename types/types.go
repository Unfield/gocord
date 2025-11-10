package types

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	IsBot      bool   `json:"bot"`
}

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ChannelID string `json:"channel_id"`
	Author    User   `json:"author"`
}
