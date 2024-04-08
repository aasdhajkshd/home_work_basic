package main

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadSensorData(t *testing.T) {
	runtime.GOMAXPROCS(1)
	c := make(chan Sensor, 10)
	v := []float64{1.0, 2.5, 3.7, 4.2, 5.3, 6.8, 7.1, 8.6, 9.0, 10.5}
	r := []float64{}
	go func() {
		defer close(c)
		for _, j := range v {
			s := Sensor{
				Time:  time.Now(),
				Value: make(map[string]float64, 1),
			}
			s.Value["test"] = j
			s.Time = time.Now()
			c <- s
		}
	}()
    for i := range c {
        r = append(r, i.Value["test"])
	}
	assert.Equal(t, v, r)

}

func TestAverageSensorData(t *testing.T) {
	runtime.GOMAXPROCS(1)
	c := make(chan Sensor)
	v := []float64{1.0, 2.5, 3.7, 4.2, 5.3, 6.8, 7.1, 8.6, 9.0, 10.5}
	r := 5.87
	go func() {
		defer close(c)
		for _, j := range v {
			s := Sensor{
				Time:  time.Now(),
				Value: make(map[string]float64, 1),
			}
			s.Value["average"] = j
			s.Time = time.Now()
			// fmt.Println(s)
			c <- s
		}
	}()
	a := averageSensorData(c)
	for i := range a {
		assert.Equal(t, r, i.Value["average"])
	}
	// time.Sleep(time.Minute)
}
