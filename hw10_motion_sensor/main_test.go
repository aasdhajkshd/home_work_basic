package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadSensorData(t *testing.T) {
	c := make(chan map[string]float64, 10)
	v := []float64{1.0, 2.5, 3.7, 4.2, 5.3, 6.8, 7.1, 8.6, 9.0, 10.5}
	r := []float64{}
	go func() {
		defer close(c)
		for _, j := range v {
			s := make(map[string]float64, 1)
			s["test"] = j
			c <- s
		}
	}()
	for i := range c {
		r = append(r, i["test"])
	}
	assert.Equal(t, v, r)
}

func TestAverageSensorData(t *testing.T) {
	c := make(chan map[string]float64)
	v := []float64{1.0, 2.5, 3.7, 4.2, 5.3, 6.8, 7.1, 8.6, 9.0, 10.5}
	r := 5.87
	go func() {
		defer close(c)
		for _, j := range v {
			s := make(map[string]float64, 1)
			s["average"] = j
			// fmt.Println(s)
			c <- s
		}
	}()
	a := averageSensorData(c)
	for i := range a {
		assert.Equal(t, r, i["average"])
	}
	// time.Sleep(time.Minute)
}
