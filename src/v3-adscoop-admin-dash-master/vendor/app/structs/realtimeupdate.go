package structs

import "encoding/json"

type RealTimeUpdate struct {
	Event string `json:"event"`
	Data  struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Message string `json:"message"`
	} `json:"data"`
}

func (r RealTimeUpdate) JSONify() (b []byte) {
	b, _ = json.Marshal(&r)
	return b
}
