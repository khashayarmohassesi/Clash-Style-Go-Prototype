package models

//The purpose of these models are for setting up the request structures
//not the database models, for purposes of setting up this prototype I took the requirements
//of anti-cheating as the main driving force for showing the approach to server-authoritative
//backend

type UpgradeRequest struct {
	BuildingID int64 `json:"building_id"`
}

type BattleRequest struct {
	//AttackerID  int64 `json:"attacker_id"` //This is omitted since this model is sent from the attacker only
	//and the attacker ID is derived from the middleware
	//setting up this attackerID with the request leaves room for cheating/fraudulent attacks
	DefenderID  int64 `json:"defender_id"`
	Stars       int   `json:"stars"`
	Destruction int   `json:"destruction"`
	LootGold    int64 `json:"loot_gold"`
}
