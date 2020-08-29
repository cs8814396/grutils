package main

import (
	model "data_fetcher/data_fetchermodel"
	"github.com/gdgrc/grutils/grframework"
	"testing"
)

func testInsertData(t *testing.T) {
	var ctx grframework.Context

	//---------test condition, default is not permit empty------------
	req := &model.InsertDataReq{}
	req.DataName = "test_insert"
	req.RecordList = [][]interface{}{

		[]interface{}{
			"11", 22,
		},

		[]interface{}{
			"33", 44,
		},
	}
	rsp := &model.InsertDataRsp{}

	err := InsertData(&ctx, req, rsp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", rsp)
}

func TestInsertData(t *testing.T) {
	testInsertData(t)
}
