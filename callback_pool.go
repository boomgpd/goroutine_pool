package goroutine_pool

import (
	"fmt"
	"sync"
)

type callbackPool struct {
	basePool
}

func NewCallbackPool(size int) *callbackPool {
	return &callbackPool{
		basePool: basePool{
			queue: make(chan struct{}, size),
			wg:    sync.WaitGroup{},
		},
	}
}

//Process 运行协程func
func (p *callbackPool) Process(handle func(any), args any) {
	p.Add(1)
	go func(args any) {
		defer p.Done()
		//先收集错误日志，再done
		defer func() {
			if err := recover(); err != nil {
				p.collectErr(fmt.Errorf("process panic error:%v", err))
			}
		}()
		handle(args)
	}(args)
}

//ProcessErrCallback 运行协程func,错误回调处理
func (p *callbackPool) ProcessErrCallback(handle func(any), args any, errCall func(err error)) {
	p.Add(1)
	go func(args any) {
		defer func() {
			if err := recover(); err != nil {
				p.collectErr(fmt.Errorf("process panic error:%v", err))
				errCall(fmt.Errorf("process panic error:%v", err))
			}
		}()
		//先done，再执行callback，以防外部callback时panic，协程未终止
		defer p.Done()
		handle(args)
	}(args)
}
