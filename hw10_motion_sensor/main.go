package main

import (
	"fmt"
	"math/rand"
	"time"
)

func readData() float64 {
	r := rand.Float64() * 100 //nolint:gosec
	// r := rand.Intn(100)
	return r
}

func readSensorData(s string, d time.Duration) chan map[string]float64 {
	c := make(chan map[string]float64)
	data := make(map[string]float64, 1)
	go func() {
		t := time.NewTimer(d)
		defer close(c)
		for {
			data[s] = readData()
			c <- data
			select {
			case <-t.C:
				fmt.Println("Timer for one minute expired")
				return
			default:
				time.Sleep(500 * time.Millisecond) // для исключения "спама"
			}
		}
	}()
	return c
}

func displaySensorData(c <-chan map[string]float64) {
	go func() {
		fmt.Printf("Main goroutine started, receiving data... on channel: %v\n", c)
		for data := range c {
			for i, j := range data {
				fmt.Println()
				fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
				fmt.Printf("Average sum sensor data %s: %.2f°C\n", i, j)
			}
		}
	}()
}

func averageSensorData(c <-chan map[string]float64) <-chan map[string]float64 {
	o := make(chan map[string]float64)
	go func() {
		defer close(o)
		m := []float64{}
		for data := range c {
			s := make(map[string]float64, 10)
			for i, j := range data {
				m = append(m, j)
				for _, v := range m {
					s[i] += v
				}
				// fmt.Println(m, s[i])
				if len(m)%10 == 0 {
					m = []float64{}
					s[i] /= 10
					// fmt.Printf("average: %.2f\n", s[i])
					o <- s
				}
			}
		}
	}()
	return o
}

func main() {
	startTimer := time.Now()
	t := readSensorData("temperature", 60*time.Second)
	a := averageSensorData(t)
	displaySensorData(a)
	for range t {
		fmt.Printf(".")
	}
	endTimer := time.Now()
	fmt.Println("Total execution time:", endTimer.Sub(startTimer))
}
