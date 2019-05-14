package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	m *sync.RWMutex
)

func foo() {
	time.Sleep(time.Second)
	fmt.Println("Hello World!")
}

func bar() {
	defer func() {
		fmt.Println("bar")
		m.RUnlock()
	}()
	m.RLock()
	go foo()
}

func baz() {
	defer m.Unlock()
	m.Lock()
	fmt.Println("baz")
}

func main() {
	m = new(sync.RWMutex)
	bar()
	baz()
	time.Sleep(time.Second * 3)
}
