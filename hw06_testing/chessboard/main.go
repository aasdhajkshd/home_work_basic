package chessboard

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

func YesNo(question string) bool {
	fmt.Printf("%s [y/n]: ", question)

	_, e := fmt.Fscan(os.Stdin, &question)
	if e != nil {
		log.Fatal(e)
	}
	return strings.ToLower(strings.TrimSpace(question)) == "y"
}

func DrawChessBoard(x, y int) string {
	var board string
	if x < 2 || y < 2 || x > 12 || y > 12 {
		msg := fmt.Sprintf("incompatiple board values: %d x %d", x, y)
		return msg
	}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			board = board + "|"
			if j%2 == i%2 {
				board = board + "#"
			} else {
				board = board + " "
			}
		}
		board = board + "|\n"
	}
	hasher := md5.New()
	hasher.Write([]byte(board))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// fmt.Printf("Наша доска:\n%v", board)
	// fmt.Println("MD5 Hash:", hash)
	return hash
}

func SizeOfBoard() int {
	answer := 8
	for i := 0; i < 3; i++ { // три попытки на неверный ввод, так как ожидается цифра
		_, e := fmt.Fscanln(os.Stdin, &answer) // сюда еще добавить ограничение по цифре по range
		if e != nil {
			fmt.Println("Ошибка:", e)
			if YesNo("Попробовать еще раз:") { // здесь можно ответить и нет, тогда уходим в else
				continue
			}
		}
		break
	}
	return answer
}

func ChessBoard() {
	fmt.Print("Введите количество ячеек: ") // строк
	x := SizeOfBoard()
	// fmt.Print("Введите количество столбцов: ")
	y := x
	fmt.Println("Ниже вывод доски размером", x, "x", y)
	fmt.Println(DrawChessBoard(x, y))
}
