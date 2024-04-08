package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Sensor struct {
	Time  time.Time
	Value map[string]float64
}

func readData() float64 {
	r := rand.Float64() * 100 //nolint:gosec
	return r
}

func readSensorData(s *Sensor, d time.Duration) chan Sensor {
	c := make(chan Sensor)
	go func() {
		t := time.NewTimer(d)
		defer close(c)
		for {
			s.Value["temperature"] = readData()
			s.Time = time.Now()
			c <- *s
			select {
			case <-t.C:
				fmt.Println("Timer for one minute expired")
				return
			default:
				time.Sleep(1 * time.Second) // для исключения "спама"
			}
		}
	}()
	return c
}

func displaySensorData(c <-chan Sensor) {
	go func() {
		fmt.Printf("Main goroutine started, receiving data... on channel: %v\n", c)
		for data := range c {
			for i, j := range data.Value {
				fmt.Printf("Time: %s\n", data.Time.Format("2006-01-02 15:04:05"))
				fmt.Printf("Sensor data %s: %.2f°C\n", i, j)
			}
		}
	}()
}

func averageSensorData(c <-chan Sensor) <-chan Sensor {
	runtime.GOMAXPROCS(1)
	o := make(chan Sensor)
	go func() {
		defer close(o)
		m := []float64{}
		for data := range c {
			s := Sensor{
				Time:  time.Now(),
				Value: make(map[string]float64, 10),
			}
			for i, j := range data.Value {
				m = append(m, j)
				for _, v := range m {
					s.Value[i] += v
				}
				// fmt.Println(m, s.Value[i])
				if len(m)%10 == 0 {
					m = []float64{}
					s.Value[i] /= 10
					// fmt.Printf("average: %.2f\n", s.Value[i])
					o <- s
				}
			}
		}
	}()
	return o
}

func main() {
	runtime.GOMAXPROCS(3)
	var input *Sensor = &Sensor{ //nolint:revive,stylecheck
		Value: make(map[string]float64, 100),
	}
	startTimer := time.Now()
	t := readSensorData(input, 60*time.Second)
	a := averageSensorData(t)
	displaySensorData(a)
wait:
	for { //nolint:gosimple
		select {
		case _, ok := <-t:
			if !ok {
				break wait
			}
		}
	}
	endTimer := time.Now()
	fmt.Println("...\n", endTimer.Sub(startTimer))
}
