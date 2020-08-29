package data_fetchermodel

type InsertDataReq struct {
	DataName   string          `json:"data_name"`
	Delay      bool            `toml:"delay"`
	RecordList [][]interface{} `json:"record_list"`
}

type InsertDataRsp struct {
}
