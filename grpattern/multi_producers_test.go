package grpattern

import (
	//"encoding/json"
	//"fmt"
	"sync"
	"testing"
	"time"
)

//==========

type MultiProducerParams struct {
	//	Wg    *sync.WaitGroup //must be pointer or can not be transfer into here

	Begin int
	End   int

	//SouceData *[]map[string]string
}

func (mpp *MultiProducerParams) Produce() (interface{}, error) {

	time.Sleep(time.Millisecond * 1)

	return mpp.Begin, nil

}

var endNum = 1000
var m sync.Mutex

var countNum = 0

func Consume(ps []Product) {
	time.Sleep(time.Millisecond * 1)
	/*
		fmt.Println("")
		for _, p := range ps {
			fmt.Printf("%+v,", p)
		}*/
	m.Lock()
	defer m.Unlock()

	countNum += len(ps)

}

/*
go test --bench=BenchmarkWorkPool  -benchmem

go test --bench="BenchmarkWithArray*"  -benchmem
*/
func BenchmarkWithArrayCap(b *testing.B) {

	bulkNum := 10

	ps := make([]Product, 0, bulkNum)

	nowNum := 0

	for i := 0; i < b.N; i++ {

		if nowNum >= bulkNum {
			ps = ps[0:0]
			nowNum = 0

			//b.Log(cap(ps), len(ps))
		}

		p := 1000
		ps = append(ps, p)
		nowNum = nowNum + 1

	}

}
func BenchmarkWithArrayNoCap(b *testing.B) {

	bulkNum := 10

	ps := make([]Product, 0, bulkNum)

	nowNum := 0

	for i := 0; i < b.N; i++ {

		if nowNum >= bulkNum {
			ps = ps[bulkNum:bulkNum]
			nowNum = 0

			//b.Log(cap(ps), len(ps))
		}

		p := 1000
		ps = append(ps, p)
		nowNum = nowNum + 1

	}

}

func BenchmarkWorkPoolWithOne(b *testing.B) {

	for i := 0; i < b.N; i++ {

		wp := NewWorkPool(1, Consume, 20)

		delta := 2

		for i := 1; i <= endNum; i = i + 1 {
			begin := i * delta
			end := i*delta + delta

			var mpp MultiProducerParams

			mpp.Begin = begin
			mpp.End = end

			wp.Invoke(&mpp)

		}

		wp.Close()

	}
}
func BenchmarkWorkPoolWith100(b *testing.B) {

	for i := 0; i < b.N; i++ {

		wp := NewWorkPool(100, Consume, 20)

		delta := 2

		for i := 1; i <= endNum; i = i + 1 {
			begin := i * delta
			end := i*delta + delta

			var mpp MultiProducerParams

			mpp.Begin = begin
			mpp.End = end

			wp.Invoke(&mpp)

		}

		wp.Close()

	}
}
func Test_WorkPool(t *testing.T) {

	wp := NewWorkPool(11, Consume, 20)

	delta := 2

	for i := 1; i <= endNum; i = i + 1 {
		begin := i * delta
		end := i*delta + delta

		var mpp MultiProducerParams

		mpp.Begin = begin
		mpp.End = end

		wp.Invoke(&mpp)

	}

	wp.Close()
	if countNum != endNum {
		t.Fatalf("not equal: countNum: %d, endNum: %d", countNum, endNum)
	}
}
