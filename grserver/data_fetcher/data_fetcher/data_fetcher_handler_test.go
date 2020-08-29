package main

import (
	"github.com/gdgrc/grutils/grframework"
	model "github.com/gdgrc/grutils/grserver/data_fetcher/data_fetchermodel"
	"testing"
)

func testFetchDataConditionSecurity(t *testing.T) {
	var ctx grframework.Context

	//---------test condition, default is not permit empty------------
	req := &model.FetchDataReq{}
	req.DataName = "test_data"
	req.PageSize = 50

	rsp := &model.FetchDataRsp{}

	err := FetchData(&ctx, req, rsp)
	if err == nil {
		t.Fatal(err)
	}
	//----------default is not permit empty------------
	req = &model.FetchDataReq{}

	req.DataName = "test_data"
	req.PageSize = 50

	req.Condition = map[string]map[string][]string{
		"test_condition1": map[string][]string{"in": []string{"aa", "bb", "ccc"}},
		"test_condition2": map[string][]string{"gte": []string{"100"}}}

	err = FetchData(&ctx, req, rsp)

	if err != nil {
		t.Fatal(err)
	}
	if rsp.TotalRecordNum <= 0 {
		t.Fatal("record num should not be zero")
	}

	if len(rsp.RecordList) == 0 {
		t.Fatal("record list should not be zero")
	}

	t.Logf("%v", rsp)
}

func TestFetchData(t *testing.T) {
	testFetchDataConditionSecurity(t)
}
