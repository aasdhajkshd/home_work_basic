package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	countTimes int = 777777
	countIter  int = 3
)

func intCount(counter *int, wg *sync.WaitGroup, mutex bool) {
	mu := sync.Mutex{}
	for i := 0; i < countTimes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if mutex {
				mu.Lock()
			}
			*counter++
			d := countTimes / 100
			if i%d == 0 {
				fmt.Printf(".")
			}
			if mutex {
				mu.Unlock()
			}
		}(i)
	}
}

type MapCounters struct {
	mx sync.Mutex
	m  map[int]int
}

func NewCounters() *MapCounters {
	return &MapCounters{
		m: make(map[int]int),
	}
}

func (c *MapCounters) Store(key, value int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}

func (c *MapCounters) Load(key int) (int, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *MapCounters) Inc(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key]++
}

func mapCount(counters *MapCounters, wg *sync.WaitGroup) {
	for i := 0; i <= countTimes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			counters.Inc(0)
			d := countTimes / 100
			if i%d == 0 {
				fmt.Printf(".")
			}
		}(i)
	}
}

func main() {
	wg := &sync.WaitGroup{}
	// Без mutex будет race
	runtime.GOMAXPROCS(1)

	fmt.Println("Running as is on single CPU:")
	for i := countIter; i > 0; i-- {
		c := 0
		startTimer := time.Now()
		intCount(&c, wg, false)
		wg.Wait()
		endTimer := time.Now()
		fmt.Printf("time: %v, iter: %d, counter: %d\n", endTimer.Sub(startTimer), i, c)
	}

	fmt.Println("Setting CPUs:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Running as is again:")
	for i := countIter; i > 0; i-- {
		c := 0
		startTimer := time.Now()
		intCount(&c, wg, false)
		wg.Wait()
		endTimer := time.Now()
		fmt.Printf("time: %v, iter: %d, counter: %d\n", endTimer.Sub(startTimer), i, c)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Running with mutex:")
	for i := countIter; i > 0; i-- {
		c := 0
		startTimer := time.Now()
		intCount(&c, wg, true)
		wg.Wait()
		endTimer := time.Now()
		fmt.Printf("time: %v, iter: %d, counter: %d\n", endTimer.Sub(startTimer), i, c)
	}

	fmt.Println("Running map:")
	for i := countIter; i > 0; i-- {
		mapCounter := NewCounters()
		startTimer := time.Now()
		// mapCount(mapCounter, wg)
		mapCount(mapCounter, wg)
		wg.Wait()
		endTimer := time.Now()
		v, _ := mapCounter.Load(0)
		fmt.Printf("time: %v, iter: %d, counter: %d\n", endTimer.Sub(startTimer), i, v)
	}
}
