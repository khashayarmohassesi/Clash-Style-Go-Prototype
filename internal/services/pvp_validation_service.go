package services

import (
	"database/sql"
	"errors"

	"server/internal/models"
)

type PVPService struct{ db *sql.DB }

func NewPVPService(db *sql.DB) *PVPService { return &PVPService{db} }

// The battle Mock validation
// the battle request is only valid if all the conditions are met
// the client could be thinner, since this is very dependent on the design of the game
// the client could just send an attack command and waits for the result as the player just chooses the forces sending to another player's base
// and the first layer and most important validation layer is taken at the point of starting the attack (validating what forces are being sent with what tier)
// highly dependent on the design of the game, the assumption for this prototype is, there's no live attack commands, all the data that's being sent
// from the client is getting validated through and through
func (s *PVPService) SubmitBattle(attackerID int64, req models.BattleRequest) error {
	if req.DefenderID == attackerID {
		return errors.New("self attack")
	}
	if req.Stars < 0 || req.Stars > 3 {
		return errors.New("bad stars")
	}
	if req.Destruction < 0 || req.Destruction > 100 {
		return errors.New("bad destruction")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var defenderGold int64
	err = tx.QueryRow(`SELECT gold FROM resources WHERE player_id=$1 FOR UPDATE`, req.DefenderID).Scan(&defenderGold)
	if err != nil {
		return err
	}

	maxLoot := defenderGold / 10
	loot := req.LootGold
	if loot > maxLoot {
		loot = maxLoot
	}

	_, err = tx.Exec(`UPDATE resources SET gold=gold-$1 WHERE player_id=$2`, loot, req.DefenderID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE resources SET gold=gold+$1 WHERE player_id=$2`, loot, attackerID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE players SET xp=xp+$1,trophies=trophies+$2 WHERE id=$3`, req.Stars*10, req.Stars*5, attackerID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO battles(attacker_id,defender_id,stars,destruction,loot_gold,validated,created_at)
	VALUES($1,$2,$3,$4,$5,true,NOW())`, attackerID, req.DefenderID, req.Stars, req.Destruction, loot)
	if err != nil {
		return err
	}

	return tx.Commit()
}
