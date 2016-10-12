package structs

import "encoding/json"

type PubSubMessage struct {
	Event string
	Data  struct {
		IDString string
		IDInt    uint
	}
}

func (r PubSubMessage) JSONify() string {
	b, _ := json.Marshal(&r)
	return string(b)
}
