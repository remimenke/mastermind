package data

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

type Player struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Email     string      `json:"email,omitempty"`
	GamePlays []*GamePlay `json:"game_plays,omitempty"`
}

// BeforeCreate implements a gorm hook
func (s *Player) BeforeCreate(scope *gorm.Scope) error {
	if s.ID == "" {
		scope.SetColumn("id", uuid.New())
	}
	return nil
}

type Game struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Open          bool         `json:"open,omitempty"`
	Code          string       `json:"code,omitempty"`
	MaxTries      string       `json:"max_tries,omitempty"`
	WonByPlayer   *Player      `json:"won_by_player,omitempty"`
	WonByPlayerID string       `json:"won_by_player_id,omitempty"`
	GamePlays     []*GamePlay  `json:"game_plays,omitempty"`
	CodeEntries   []*CodeEntry `json:"code_entries,omitempty"`
}

// BeforeCreate implements a gorm hook
func (s *Game) BeforeCreate(scope *gorm.Scope) error {
	if s.ID == "" {
		scope.SetColumn("id", uuid.New())
	}
	return nil
}

type GamePlay struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Player   *Player `json:"player,omitempty"`
	PlayerID string  `json:"player_id,omitempty"`
	Game     *Game   `json:"game,omitempty"`
	GameID   string  `json:"game_id,omitempty"`
	Host     bool    `json:"host,omitempty"`
}

// BeforeCreate implements a gorm hook
func (s *GamePlay) BeforeCreate(scope *gorm.Scope) error {
	if s.ID == "" {
		scope.SetColumn("id", uuid.New())
	}
	return nil
}

type CodeEntry struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Code     string `json:"code,omitempty"`
	Feedback string `json:"feedback,omitempty"`
	Game     *Game  `json:"game,omitempty"`
	GameID   string `json:"game_id,omitempty"`
}

// BeforeCreate implements a gorm hook
func (s *CodeEntry) BeforeCreate(scope *gorm.Scope) error {
	if s.ID == "" {
		scope.SetColumn("id", uuid.New())
	}
	return nil
}

type Error struct {
	Error string `json:"error"`
}
