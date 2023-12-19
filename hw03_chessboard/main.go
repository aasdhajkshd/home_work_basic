// Шахматная доска
// Программа, которая создаёт строку, содержащую решётку 8х8, а линии разделяются символами новой строки.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func yesno(question string) bool {
	fmt.Printf("%s [y/n]: ", question)

	_, e := fmt.Fscan(os.Stdin, &question)
	if e != nil {
		log.Fatal(e)
	}
	return strings.ToLower(strings.TrimSpace(question)) == "y"
}

func size() int {
	var answer int
	for i := 0; i < 3; i++ { // три попытки на неверный ввод, так как ожидается цифра
		_, e := fmt.Fscanln(os.Stdin, &answer) // сюда еще добавить ограничение по цифре по range
		if e != nil {
			fmt.Println("Ошибка:", e)
			if yesno("Попробовать еще раз:") { // здесь можно ответить и нет, тогда уходим в else
				continue
			} else {
				fmt.Println("... на нет и спроса нет, тогда 8")
				return 8
			}
		}
		break
	}
	return answer
}

func draw(x, y int) {
	hash, pipe := "#", "|"
	for i := 0; i < y; i++ { // для прохода по строкам (вертикаль)
		for j := 0; j < x; j++ { // для прохода по столбцам (горизонталь)
			time.Sleep(700 * time.Millisecond) // ... как на матричном принтере
			fmt.Printf("%s", pipe)
			if j%2 == i%2 { // вот здесь меняем местами "шашечки" по строкам, чтобы были вдоль
				fmt.Printf("%s", hash)
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println(pipe) // ставится завершающая линия и перенос строки
	}
}

func main() {
	fmt.Print("Введите количество строк: ")
	x := size()
	fmt.Print("Введите количество столбцов: ")
	y := size()
	fmt.Println("Ниже вывод доски размером", x, "x", y)
	draw(x, y)
}
