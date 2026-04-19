package core

import (
	"sync"
)

// Pool 队列结构
type Pool struct {
	// 使用channel控制并发数量
	queue chan int
	// 使用sync.WaitGroup实现协程同步
	wg *sync.WaitGroup
}

// NewPool 初始化
func NewPool(size int) *Pool {
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add 队列操作
func (p *Pool) Add(n int) {
	// 如果i大于0,则向队列中添加
	for i := 0; i < n; i++ {
		p.queue <- i
	}
	for i := 0; i > n; i-- {
		<-p.queue
	}
	p.wg.Add(n)
}

// Done 出队列
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait 等待队列操作
func (p *Pool) Wait() {
	p.wg.Wait()
}
