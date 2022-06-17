package goroutine_pool

import (
	"sync"
)

type basePool struct {
	//数量控制通道
	queue chan struct{}
	//协程组
	wg sync.WaitGroup
	//panic错误收集
	Errs []error
	//错误采集锁
	errMutex sync.Mutex
}

func (p *basePool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- struct{}{}
	}
	p.wg.Add(delta)
}

func (p *basePool) Done() {
	<-p.queue
	p.wg.Done()
}

func (p *basePool) Wait() {
	p.wg.Wait()
}

//collectErr 错误采集
func (p *basePool) collectErr(err error) {
	p.errMutex.Lock()
	p.Errs = append(p.Errs, err)
	p.errMutex.Unlock()
}
