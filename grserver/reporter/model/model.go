package model

//"device_filter/reporter/model"

//"github.com/mongodb/mongo-go-driver/core/result"

type ReportData interface {
	DumpBytes() ([]byte, error)
	DumpOrderedList() ([]interface{}, error)
	SetDataTime(string)
	Prepare() (err error)
	GetDataName() string
}
