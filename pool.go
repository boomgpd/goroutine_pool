package goroutine_pool

import (
	"fmt"
	"sync"
)

type pool struct {
	basePool
	//处理函数
	handle func(any)
}

func NewPool(size int, handle func(any)) *pool {
	return &pool{
		basePool: basePool{
			queue: make(chan struct{}, size),
			wg:    sync.WaitGroup{},
		},
		handle: handle,
	}
}

//Process 运行协程func
func (p *pool) Process(args any) {
	p.Add(1)
	go func(args any) {
		//先收集错误日志，再done
		defer p.Done()
		defer func() {
			if err := recover(); err != nil {
				p.collectErr(fmt.Errorf("process panic error:%v", err))
			}
		}()
		p.handle(args)
	}(args)
}

//ProcessErrCallback 运行协程func,错误回调处理
func (p *pool) ProcessErrCallback(args any, errCall func(err error)) {
	p.Add(1)
	go func(args any) {
		defer p.Done()
		defer func() {
			if err := recover(); err != nil {
				p.collectErr(fmt.Errorf("process panic error:%v", err))
				errCall(fmt.Errorf("process panic error:%v", err))
			}
		}()
		p.handle(args)
	}(args)
}
