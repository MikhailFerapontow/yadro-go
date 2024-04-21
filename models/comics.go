package models

type ResponseComic struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type WeightedWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type DbComic struct {
	Id       int            `json:"id"`
	Url      string         `json:"url"`
	Keywords []WeightedWord `json:"keywords"`
}
