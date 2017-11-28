package mastermind

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/remimenke/mastermind/data"
)

func (c *controller) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		player := &data.Player{}

		err := json.NewDecoder(r.Body).Decode(player)
		if err != nil {
			err := json.NewEncoder(w).Encode(data.Error{Error: err.Error()})
			if err != nil {
				log.Println(err)
			}
			return
		}

		if player.Email == "" {
			err := json.NewEncoder(w).Encode(data.Error{Error: "missing email"})
			if err != nil {
				log.Println(err)
			}
			return
		}

		err = c.db.
			Model(data.Player{}).
			Where("email = ?", player.Email).
			Preload("GamePlays").
			Preload("GamePlays.Game").
			FirstOrCreate(player).
			Error
		if err != nil {
			err := json.NewEncoder(w).Encode(data.Error{Error: err.Error()})
			if err != nil {
				log.Println(err)
			}
			return
		}

	})
}
