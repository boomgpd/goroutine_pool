package goroutine_pool

import (
	"fmt"
	"runtime"
	"testing"
)

func printNum(i any) {
	fmt.Println(i)
}

func panicNum(i any) {
	if i.(int) > 5 {
		panic(i)
	} else {
		printNum(i)
	}
}

func errorHandle(err error) {
	fmt.Println("custom handle err: ", err.Error())
}

func TestPool(t *testing.T) {
	numCPUs := runtime.NumCPU()
	//正常流程测试
	t.Run("test pool", func(t *testing.T) {
		var pool = NewPool(numCPUs, printNum)
		for i := 0; i < 10; i++ {
			pool.Process(i)
		}
		pool.Wait()
		if len(pool.Errs) > 0 {
			t.Error(pool.Errs)
		}
	})

	//panic流程测试
	t.Run("test pool panic", func(t *testing.T) {
		var pool = NewPool(numCPUs, panicNum)
		for i := 0; i < 10; i++ {
			fmt.Println("num:", i)
			pool.Process(i)
		}
		pool.Wait()
		if len(pool.Errs) > 0 && len(pool.Errs) != 4 {
			t.Error(pool.Errs)
		}
	})

	//panic-callback-handle测试
	t.Run("test pool panic callback", func(t *testing.T) {
		var pool = NewPool(numCPUs, panicNum)
		for i := 0; i < 10; i++ {
			pool.ProcessErrCallback(i, errorHandle)
		}
		pool.Wait()
	})

}

func TestCallbackPool(t *testing.T) {
	numCPUs := runtime.NumCPU()
	//回调处理测试
	t.Run("test callback pool", func(t *testing.T) {
		var pool = NewCallbackPool(numCPUs)
		for i := 0; i < 10; i++ {
			pool.Process(printNum, i)
		}
		pool.Wait()
		if len(pool.Errs) > 0 {
			t.Error(pool.Errs)
		}
	})

	//回调处理-panic测试
	t.Run("test callback pool panic", func(t *testing.T) {
		var pool = NewCallbackPool(numCPUs)
		for i := 0; i < 10; i++ {
			pool.Process(panicNum, i)
		}
		pool.Wait()
		if len(pool.Errs) > 0 && len(pool.Errs) != 4 {
			t.Error(pool.Errs)
		}
	})

	//回调处理-panic-回调处理测试
	t.Run("test callback pool panic", func(t *testing.T) {
		var pool = NewCallbackPool(numCPUs)
		for i := 0; i < 10; i++ {
			pool.ProcessErrCallback(panicNum, i, errorHandle)
		}
		pool.Wait()
	})
}
