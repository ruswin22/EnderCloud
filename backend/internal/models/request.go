package models

type ServerRequest struct {
	WorldName string `json:"world_name"`
	Seed      string `json:"seed"`
}
