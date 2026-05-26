package services

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

//The player login flow is now completely based off of Guest logins, later Oauth providers could be added

type PlayerService struct{ db *sql.DB }

func NewPlayerService(db *sql.DB) *PlayerService { return &PlayerService{db} }

func (s *PlayerService) GuestLogin() (map[string]any, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	gid := uuid.NewString()
	var playerID int64
	err = tx.QueryRow(`INSERT INTO players(auth_provider,provider_user_id,name,xp,trophies,level,created_at,last_login)
	VALUES('guest',$1,'Guest',0,0,1,NOW(),NOW()) RETURNING id`, gid).Scan(&playerID)
	if err != nil {
		return nil, err
	}

	var baseID int64
	err = tx.QueryRow(`INSERT INTO bases(player_id,layout_version) VALUES($1,1) RETURNING id`, playerID).Scan(&baseID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO resources(player_id,gold,elixir,gems,updated_at) VALUES($1,1000,1000,50,NOW())`, playerID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO buildings(base_id,type,level,hp,state) VALUES($1,'townhall',1,1000,'idle')`, baseID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return map[string]any{"player_id": playerID, "provider_user_id": gid}, nil
}

func (s *PlayerService) GetProfile(id int64) (map[string]any, error) {
	var name string
	var xp, trophies, level int
	var lastLogin time.Time
	err := s.db.QueryRow(`SELECT name,xp,trophies,level,last_login FROM players WHERE id=$1`, id).
		Scan(&name, &xp, &trophies, &level, &lastLogin)
	return map[string]any{"id": id, "name": name, "xp": xp, "trophies": trophies, "level": level}, err
}
