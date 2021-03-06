package model

type UpdateDataReq struct {
	DataName string `json:"data_name"`
	//Delay      bool            `toml:"delay"`
	Condition map[string]map[string][]string `protobuf:"bytes,4,opt,name=condition,proto3" json:"condition"`

	//RecordList [][]interface{} `json:"record_list"`
}

type UpdateDataRsp struct {
	IdList      []int64 `json:"id_list"`
	RowAffected int64   `json:"row_affected"`
}
