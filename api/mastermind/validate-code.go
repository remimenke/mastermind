package mastermind

import "net/http"

func (c *controller) ValidateCode() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
