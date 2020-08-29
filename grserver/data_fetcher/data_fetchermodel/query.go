package data_fetchermodel

import (
	"encoding/json"
	//proto "github.com/golang/protobuf/proto"
)

/*
service data_fetcher {
    rpc FetchData(FetchDataReq) returns (FetchDataRsp) {}
}

message FetchDataReq {
	string data_name=1;
	int32 page_no=2;
	int32 page_size=3;
	map<string,Condition> condition=4;
}




message Condition {
	map<string, Range> range=1;
}


message Range {
	repeated string value_list=1;
}


message FetchDataRsp {
    Data data = 1;
}


message Data {

	int32 total_page_num=1; // how many pages there are
	int32 record_num=2; // how many records return this time
	int32 page_size=3;
	int32 page_no=4;
	int32 total_record_num=5; // how many records totaly
	repeated Record record_list=6; // record

}

message Record {
	map<string,string> record=1;
}
*/

type PlaceHolder struct {
	ReplacedStatement string
	Params            []interface{}
}

type MysqlCount struct {
	Count int `json:"count"`
}
type FetchDataReq struct {
	DataName  string                         `protobuf:"bytes,1,opt,name=data_name,proto3" json:"data_name"`
	PageNo    int                            `protobuf:"varint,2,opt,name=page_no,proto3" json:"page_no"`
	PageSize  int                            `protobuf:"varint,3,opt,name=page_size,proto3" json:"page_size"`
	Condition map[string]map[string][]string `protobuf:"bytes,4,opt,name=condition,proto3" json:"condition"`
}

func (m *FetchDataReq) Reset() { *m = FetchDataReq{} }
func (m *FetchDataReq) String() string {

	b, _ := json.Marshal(m)
	return string(b) //proto.CompactTextString(m)
}
func (*FetchDataReq) ProtoMessage() {}

type FetchDataRsp struct {
	TotalPageNum   int `protobuf:"varint,1,opt,name=total_page_num,proto3" json:"total_page_num,omitempty"`
	RecordNum      int `protobuf:"varint,2,opt,name=record_num,proto3" json:"record_num,omitempty"`
	TotalRecordNum int `protobuf:"varint,3,opt,name=total_record_num,proto3" json:"total_record_num,omitempty"`
	PageSize       int `protobuf:"varint,4,opt,name=page_size,proto3" json:"page_size,omitempty"`
	PageNo         int `protobuf:"varint,5,opt,name=page_no,proto3" json:"page_no,omitempty"`

	RecordList []map[string]string `protobuf:"bytes,6,opt,name=record_list,proto3" json:"record_list,omitempty"`
}

func (m *FetchDataRsp) Reset() { *m = FetchDataRsp{} }
func (m *FetchDataRsp) String() string {

	b, _ := json.Marshal(m)
	return string(b) //proto.CompactTextString(m)
}

func (*FetchDataRsp) ProtoMessage() {}

var RuleTable = map[string]string{
	"gte": ">=",
	"lte": "<=",
	"lt":  "<",
	"eq":  "=",
	"in":  "in",
}
