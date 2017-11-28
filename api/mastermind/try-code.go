package mastermind

import "net/http"

func (c *controller) TryCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
