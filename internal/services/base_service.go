package services

import "database/sql"

type BaseService struct{ db *sql.DB }

func NewBaseService(db *sql.DB) *BaseService { return &BaseService{db} }

//Loading base function is written with many assumptions which might not map to the design documents, and there's a high chance
//I'd recommend a JSONB for serialization purposes, and sharing it between the unity client is nice as well
func (s *BaseService) LoadBase(playerID int64) (map[string]any, error) {
	rows, err := s.db.Query(`SELECT b.id,b.type,b.level,b.hp,b.state
	FROM buildings b JOIN bases ba ON ba.id=b.base_id WHERE ba.player_id=$1`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []map[string]any
	for rows.Next() {
		var id int64
		var t, state string
		var lvl, hp int
		if err = rows.Scan(&id, &t, &lvl, &hp, &state); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"id": id, "type": t, "level": lvl, "hp": hp, "state": state})
	}
	return map[string]any{"buildings": list}, nil
}
