package models

type ResponseComic struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type DbComic struct {
	Id       int      `json:"id"`
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}
