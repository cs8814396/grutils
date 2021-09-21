package model
import (
	"encoding/json"
)
type CommReportData struct {
	Event string `json:"event"`
	Ip    string `json:"ip"`
	Area  string `json:"area"`
	Time string `json:"time"`
	Data map[string]interface{} `json:"data"`
}

func (dd *CommReportData) GetDataName() (string) {
	return "comm_report"
}
func (dd *CommReportData) DumpOrderedList() (il []interface{}, err error) {
	
	
	il = []interface{}{
		
	}
	return
}
func (dd *CommReportData) DumpBytes() ([]byte, error) {

	return json.Marshal(dd)
}
func (dd *CommReportData) SetDataTime(t string) {


	return
}
func (dd *CommReportData) Prepare() (err error) {



	return
}
