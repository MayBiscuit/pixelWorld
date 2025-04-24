package models

type Sticker struct {
	SID      int         `json:"sid"`
	SName    string      `json:"sname"`
	SType    string      `json:"stype"`
	SPicture interface{} `json:"spicture"`
}
