package base

type Image struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Size     string `json:"size"`
	Desc     string `json:"describe"`
	IsDelete bool   `json:"isDelete"`
}
