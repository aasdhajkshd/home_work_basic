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
				fmt.Println("Timer one minute expired")
				return
			default:
				time.Sleep(1000 * time.Millisecond) // для исключения "спама"
				fmt.Printf(".")
			}
		}
	}()
	return c
}

func displaySensorData(s string, c <-chan float64) chan struct{} {
	o := make(chan struct{})
	go func() {
		defer close(o)
		for data := range c {
			fmt.Println()
			fmt.Printf("Average sum sensor data %s: %.2f°C\n", s, data)
			o <- struct{}{}
		}
	}()
	return o
}

func averageSensorData(c <-chan float64) chan float64 {
	o := make(chan float64)
	go func() {
		defer close(o)
		i := 0
		var m float64
		for data := range c {
			i++
			m += data
			if i%10 == 0 {
				m /= 10
				// fmt.Printf("average: %.2f\n", s)
				o <- m
				m = 0.0
			}
		}
	}()
	return o
}

func main() {
	startTimer := time.Now()
	c := readSensorData("temperature", 60*time.Second)
	a := averageSensorData(c)
	d := displaySensorData("temperature", a)
	for range d {
		fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}
	endTimer := time.Now()
	fmt.Println("Total execution time:", endTimer.Sub(startTimer))
}
