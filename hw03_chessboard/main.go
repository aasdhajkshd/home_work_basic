package main

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	for i := 0; i < 3; i++ {
		_, e := fmt.Fscanln(os.Stdin, &answer)
		if e != nil {
			fmt.Println("Ошибка:", e)
			if yesno("Попробовать еще раз:") {
				continue
			} else {
				fmt.Printf("... на нет и спроса нет, тогда будет всё %d", 8)
				return 8
			}
		}
		break
	}
	return answer
}

func draw(x, y int) {
	var hash, pipe = "#", "|"
    for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			fmt.Printf("%s", pipe)
			if j%2 == i%2 {
				fmt.Printf("%s", hash)
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println(pipe)
	}
}

func main() {
	fmt.Print("Введите количество строк: ")
	x := size()
	fmt.Print("Введите количество столбцов: ")
	y := size()
	fmt.Println(x, y)
	draw(x, y)
}
