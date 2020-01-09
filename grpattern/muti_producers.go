package grpattern

import (
	"fmt"
	"github.com/panjf2000/ants"
	"runtime/debug"
	"sync"
)

type Product interface{}

type ProducerParams interface {
	Produce() (interface{}, error)
	//Done()
}

type WorkPool struct {
	Pool        *ants.PoolWithFunc
	Queue       chan Product
	ConsumeFunc func([]Product)
	Wg          *sync.WaitGroup
}

func (wp *WorkPool) Invoke(pp ProducerParams) {

	wp.Wg.Add(1)
	wp.Pool.Invoke(pp)

}

func (wp *WorkPool) AddQueue(p Product) {
	wp.Queue <- p

	wp.Wg.Done()

}

func (wp *WorkPool) Close() {

	wp.Wg.Wait()
	wp.Wg.Add(1)

	close(wp.Queue)
	wp.Wg.Wait()
}

func NewWorkPool(producerNum int, cf func([]Product), bulkNum int) WorkPool {
	wp := WorkPool{}

	wp.Wg = &sync.WaitGroup{}

	var err error
	wp.Pool, err = ants.NewPoolWithFunc(producerNum, func(i interface{}) {

		defer func() {
			if r := recover(); r != nil {
				errString := fmt.Sprintf("%s stack: %s", r, debug.Stack())

				panic(errString)
			}
		}()

		fqp, _ := i.(ProducerParams)
		//defer fqp.Done()

		p, err := fqp.Produce()

		if err != nil {
			return
		}

		wp.AddQueue(p)

	})

	if err != nil {
		panic("ants NewPoolWithFunc init fail err: " + err.Error())
	}

	wp.Queue = make(chan Product, 65535)

	wp.ConsumeFunc = cf

	//consume goroutine

	go func() {
		defer wp.Wg.Done()
		var ps []Product

		nowNum := 0

		//bulkNum := 10

		ps = make([]Product, 0, bulkNum)

		for {

			p, ok := <-wp.Queue
			if !ok {

				break
			}
			if nowNum >= bulkNum {
				wp.ConsumeFunc(ps)

				nowNum = 0
				ps = ps[0:0] // this can reuse the space we allocate before rather than new or make one.

			}

			ps = append(ps, p)
			nowNum = nowNum + 1

		}
		if nowNum > 0 {
			wp.ConsumeFunc(ps)

		}

	}()

	return wp
}
