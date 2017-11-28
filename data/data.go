package data

import (
	"database/sql"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// NewDatabase returns a new database connection for
func NewDatabase(url string) (*gorm.DB, error) {

	dbconf, err := parseDatabaseURL(url)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("mysql", dbconf.FormatDSN())
	if err != nil {
		return nil, err
	}

	{
		cdb := db.CommonDB().(*sql.DB)
		cdb.SetMaxIdleConns(2)
		cdb.SetMaxOpenConns(5)
		cdb.SetConnMaxLifetime(10 * time.Minute)
	}

	registerCallbacks(db)
	err = migrate(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func parseDatabaseURL(urlString string) (*mysql.Config, error) {
	urlString = strings.Trim(urlString, "'")

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	conf := &mysql.Config{
		Strict:    true,
		Net:       "tcp",
		Addr:      u.Host,
		DBName:    strings.TrimPrefix(u.Path, "/"),
		ParseTime: true,
	}

	if u.User != nil {
		if u.User.Username() != "" {
			conf.User = u.User.Username()
		}
		if pwd, ok := u.User.Password(); ok {
			conf.Passwd = pwd
		}
	}

	return conf, nil
}

func registerCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("nullable_assoc:run_before_create", func(s *gorm.Scope) {
		for _, field := range s.Fields() {
			if field.IsForeignKey && field.IsBlank {
				field.IsBlank = false
				s.SetColumn(field, nil)
				field.Field = reflect.ValueOf(gorm.Expr("NULL"))
			}
		}
	})

	db.Callback().Update().Before("gorm:save_before_associations").Register("nullable_assoc:run_before_update", func(s *gorm.Scope) {
		for _, field := range s.Fields() {
			if field.IsForeignKey && field.IsBlank {
				field.IsBlank = false
				s.SetColumn(field, nil)
				field.Field = reflect.ValueOf(gorm.Expr("NULL"))
			}
		}
	})
}

// Migrate the database
func migrate(db *gorm.DB) error {
	db = db.AutoMigrate(
		Player{},
		GamePlay{},
		Game{},
		CodeEntry{},
	)

	db.Model(Game{}).
		AddForeignKey("won_by_player_id", "players(id)", "CASCADE", "RESTRICT")

	db.Model(GamePlay{}).
		AddForeignKey("code_maker_id", "players(id)", "CASCADE", "RESTRICT").
		AddForeignKey("code_breaker_id", "players(id)", "CASCADE", "RESTRICT").
		AddForeignKey("game_id", "games(id)", "CASCADE", "RESTRICT")

	db.Model(CodeEntry{}).
		AddForeignKey("game_id", "games(id)", "CASCADE", "RESTRICT")

	return db.Error
}
