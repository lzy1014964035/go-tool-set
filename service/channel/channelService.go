package channel

import (
	"sync"
)

var wg sync.WaitGroup

// 创建携程
func Create(callableFunction func()) {
	wg.Add(1)
	// 启动 WebSocket 服务器携程
	go func() {
		callableFunction();
		wg.Done();
	}()
}

// 等待携程运行
func Wait() {
	wg.Wait();
}