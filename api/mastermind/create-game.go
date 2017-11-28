package mastermind

import "net/http"

func (c *controller) CreateGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
