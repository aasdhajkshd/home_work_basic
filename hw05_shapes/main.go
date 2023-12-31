package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"github.com/aasdhajkshd/home_work_basic/hw05_shapes/shapes"
)

func calculateArea(s any) (float64, error) {
	if shape, ok := s.(shapes.Shape); ok {
		return shape.Area(), nil
	} else {
		return 0, errors.New("переданный объект не реализует интерфейс Shape")
	}
}

func recoverFromPanic() {
    if r := recover(); r != nil {
        fmt.Println("Продолжение:", r) // test panic/recover condition
    }
}

func main() {
	radius := 5.6
	circle := shapes.Circle{radius} // Circle struct Radius variable is exportable
	if circleArea, err := calculateArea(circle); err == nil {
		fmt.Printf("Круг: радиус %.2f\nПлощадь: %.2f\n", radius, circleArea)
	}
	fmt.Println()

	width, height := 5.2, 6.4
	rectangle := shapes.Rectangle{width, height}
	if rectangleArea, err := calculateArea(rectangle); err == nil {
		fmt.Printf("Прямоугольник: ширина %.2f, высота %.2f\nПлощадь: %.2f\n", width, height, rectangleArea)
	}
	fmt.Println()

	defer recoverFromPanic() // call recover to proceed the program

	var a, b, c float64 = 5, 7, 12.1 // to check triangle validation, change 10.1 to 12.1
	var s string
	triangle := shapes.Triangle{a, b, c} // we can change the third value from 10.1 to 12.1 to make traingle invalid
	if shapes.ValidateTriangle(a, b, c) {
		if a == b && b == c {
			s = "равносторонний"
		} else if a == b || b == c || a == c {
			s = "равнобедренный"
		} else {
			s = "неравносторонний"
		}
		if triangleArea, err := calculateArea(triangle); err == nil {
			fmt.Printf("Треугольник: длины сторон %.2f, %.2f, %.2f, %v \n", a, b, c, s)
			fmt.Printf("Площадь: %.2f\n", triangleArea)
		}
	} else {
		fmt.Printf("Треугольник: длины сторон %.2f, %.2f, %.2f \n\n" +
		           "Для того чтобы треугольник существовал, сумма длин любых \n" +
			       "двух его сторон должна быть больше длины третьей стороны.\n" +
			       "Если это условие не выполняется, то треугольник невозможен.\n", a, b, c)
		fmt.Println()
		panic("Нужно уменьшить длину стороны с 12.1 на 11.1 и запустить снова...") // to test abnormally terminate the program as critical test condition encountered!
	}
	fmt.Println()
	
	// Play around with exportable/none-exportable struct values
	// See shapes/pkg.go and change line 25 from Area() to CalculateSquareArea()
	square := &shapes.Square{} // assign exportable Square struct
	square.SetSquareSide(10.0) // calling setter of Square
	if squareArea, err := calculateArea(square); err != nil { // assign with check
		log.Fatal("ошибка: ", err)
	} else {
		fmt.Println("Квадрат: сторона", square.SquareSide()) // get via getter a side value of Square
		fmt.Println("Площадь:", squareArea)
	}

	os.Exit(0)
}
