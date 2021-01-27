package model

type InsertDataReq struct {
	DataName   string          `json:"data_name"`
	Delay      bool            `toml:"delay"`
	RecordList [][]interface{} `json:"record_list"`
}

type InsertDataRsp struct {
	IdList []int64 `json:"id_list"`
}
