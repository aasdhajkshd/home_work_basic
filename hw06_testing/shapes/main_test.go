package shapes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestCalculateAreaCircle(t *testing.T) {
	t.Log("Проводим тестирование для круга...")
	radiusValues := []struct {
		radius, result float64
	}{
		{5.6, 98.52},
		{10.0, 314.16},
	}

	for i, v := range radiusValues {
		result, err := CalculateAreaCircle(v.radius)
		if result != v.result && err == nil {
			t.Errorf("Тест #%d -\t%v, должно быть %.2f, но не %.2f", i+1, failed, v.result, result)
		} else {
			t.Logf("Тест #%d -\t%v", i+1, success)
		}
	}
}

func TestCalculateAreaRectangle(t *testing.T) {
	t.Log("Проводим тестирование для прямоугольника...")
	rectangleValues := []struct {
		width, height, result float64
	}{
		{5.6, 6.5, 36.4},
		{10.0, 5.12, 51.2},
		{7.123, 7.12, 50.72}, // для проверки достаточно изменить на 3 знака после запятой
	}

	for _, v := range rectangleValues {
		result, err := CalculateAreaRectangle(v.width, v.height)
		assert.Equal(t, v.result, result, "Ожидается, что %d * %d равно %d, но получено %d", v.width, v.height, result, v.result)
		assert.NoError(t, err)
	}
}

func TestCalculateAreaTriangle(t *testing.T) {
	t.Log("Проводим тестирование для треуголника...")
	triangleleValues := []struct {
		a, b, c, result float64
	}{
		{5.6, 6.5, 10.1, 16.76},
		{10.0, 5.55, 5.55, 12.04},
		{10.0, 5.55, 55.55, 0.0},
	}

	for _, v := range triangleleValues {
		result, err := CalculateAreaTriangle(v.a, v.b, v.c)
		assert.Equal(t, v.result, result, "Ожидается площадь равной %d, но получено %d", result, v.result)
		if err != nil { // просто и навсегда в реке Нил жил крокодил... Данди
			t.Log(err)
		}
	}
}
