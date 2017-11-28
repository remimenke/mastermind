package mastermind

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

type Controller interface {
	Register() http.Handler
	CreateGame() http.Handler
	JoinGame() http.Handler
	TryCode() http.Handler
	ValidateCode() http.Handler
}

func New(db *gorm.DB) Controller {
	return &controller{
		db: db,
	}
}

type controller struct {
	db *gorm.DB
}
