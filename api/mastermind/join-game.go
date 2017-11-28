package mastermind

import "net/http"

func (c *controller) JoinGame() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
