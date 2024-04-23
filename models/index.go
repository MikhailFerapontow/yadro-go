package models

type WeightedId struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

type KwIndex struct {
	Keyword string       `json:"keyword"`
	Ids     []WeightedId `json:"ids"`
}
