package gateway

type Payload struct {
	Op int     `json:"op"`
	D  any     `json:"d"`
	S  *int    `json:"s,omitempty"`
	T  *string `json:"t,omitempty"`
}
