package services

import (
	"database/sql"
	"errors"
)

type BuildingService struct{ db *sql.DB }

func NewBuildingService(db *sql.DB) *BuildingService { return &BuildingService{db} }

// this is the base server authoritative upgrade request
// the request is sent by the player, and only the request, the time is captured in the server
// the completion time is recorded by the server, all the resource requirements are checked on the server
// various different checks for the authoritativeness and correctness of data (wrong player ID, correct level, the building is not idle....)
func (s *BuildingService) StartUpgrade(playerID, buildingID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var level int
	var state string
	var owner int64
	err = tx.QueryRow(`SELECT b.level,b.state,ba.player_id
	FROM buildings b JOIN bases ba ON ba.id=b.base_id
	WHERE b.id=$1 FOR UPDATE`, buildingID).Scan(&level, &state, &owner)
	if err != nil {
		return err
	}
	if owner != playerID {
		return errors.New("not owner")
	}
	if state != "idle" {
		return errors.New("busy building")
	}

	goldCost := int64(level * 100)
	var gold int64
	err = tx.QueryRow(`SELECT gold FROM resources WHERE player_id=$1 FOR UPDATE`, playerID).Scan(&gold)
	if err != nil {
		return err
	}
	if gold < goldCost {
		return errors.New("not enough gold")
	}

	_, err = tx.Exec(`UPDATE resources SET gold=gold-$1, updated_at=NOW() WHERE player_id=$2`, goldCost, playerID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE buildings SET state='upgrading' WHERE id=$1`, buildingID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO upgrades(building_id,from_level,to_level,start_at,finish_at,status)
	VALUES($1,$2,$3,NOW(),NOW()+(($2*30)||' seconds')::interval,'pending')`, buildingID, level, level+1)
	if err != nil {
		return err
	}

	return tx.Commit()
}
