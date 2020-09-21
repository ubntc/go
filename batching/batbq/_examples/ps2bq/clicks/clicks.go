package clicks

import (
	"encoding/json"
	"log"
	"time"
)

// Click describes the context of a click on an Ad.
type Click struct {
	ID     string    `json:"id"`
	Origin string    `json:"origin"`
	Time   time.Time `json:"time"`
}

func init() {
	_, err := json.Marshal(Click{ID: "click1", Time: time.Now().UTC(), Origin: "init"})
	if err != nil {
		log.Fatal(err)
	}
}
