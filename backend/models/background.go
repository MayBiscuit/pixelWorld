package models

type Background struct {
	BID      int         `json:"bid"`
	BName    string      `json:"bname"`
	BType    string      `json:"btype"`
	BPicture interface{} `json:"bpicture"`
}
