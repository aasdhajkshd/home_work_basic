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

func mapCount(mapCounter *sync.Map, wg *sync.WaitGroup) {
	for i := 0; i <= countTimes; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mapCounter.Store(0, i)
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

	// Проверка по заполнению map
	fmt.Println("Running map:")
	for i := countIter; i > 0; i-- {
		mc := sync.Map{}
		startTimer := time.Now()
		mapCount(&mc, wg)
		wg.Wait()
		endTimer := time.Now()
		v, _ := mc.Load(0)
		fmt.Printf("time: %v, iter: %d, counter: %d\n", endTimer.Sub(startTimer), i, v)
	}
}
