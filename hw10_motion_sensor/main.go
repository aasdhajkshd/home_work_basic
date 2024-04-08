package main

import (
	"fmt"
	"math/rand"
	"time"
)

func readData(s string) map[string]float64 {
	d := make(map[string]float64, 1)
	d[s] = rand.Float64() * 100 //nolint:gosec
	return d
}

func readSensorData(s string, d time.Duration) chan float64 {
	c := make(chan float64)
	go func() {
		t := time.NewTimer(d)
		defer close(c)
		for {
			data := readData(s)
			c <- data[s]
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

func displaySensorData(s string, c <-chan float64) {
	go func() {
		fmt.Printf("Main goroutine started, receiving data... on channel: %v\n", c)
		for data := range c {
			fmt.Println()
			fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
			fmt.Printf("Average sum sensor data %s: %.2f°C\n", s, data)
		}
	}()
}

func averageSensorData(c <-chan float64) chan float64 {
	o := make(chan float64)
	go func() {
		defer close(o)
		m := []float64{}
		for data := range c {
			m = append(m, data)
			if len(m)%10 == 0 {
				var s float64
				for _, v := range m {
					s += v
				}
				m = []float64{}
				s /= 10
				// fmt.Printf("average: %.2f\n", s)
				o <- s
			}
		}
	}()
	return o
}

func main() {
	startTimer := time.Now()
	t := readSensorData("temperature", 60*time.Second)
	a := averageSensorData(t)
	displaySensorData("temperature", a)
	for range t {
		fmt.Printf(".")
	}
	endTimer := time.Now()
	fmt.Println("Total execution time:", endTimer.Sub(startTimer))
}
