package domain

type ResponseComic struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

type WeightedWord struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Comic struct {
	Id       int            `json:"id"`
	Url      string         `json:"url"`
	Keywords []WeightedWord `json:"keywords"`
}

type WeightedId struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

type KwIndex struct {
	Keyword string       `json:"keyword"`
	Ids     []WeightedId `json:"ids"`
}
