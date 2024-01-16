package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

func countWords(x string) (map[string]int, error) { // по условаию задания принимаем строку
	wordsMap := make(map[string]int)
	s := bufio.NewScanner(strings.NewReader(x))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		w := strings.ToLower(strings.Trim(s.Text(), "!?,. "))
		wordsMap[w]++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return wordsMap, nil
}

func printDups(m map[string]int) {
	for w, c := range m {
		if c > 1 { // здесь можем указать кол-во повторяющихся слов
			switch utf8.RuneCountInString(w) {
			case 1:
				fmt.Printf("Буква")
			default:
				fmt.Printf("Слово")
			}
			fmt.Printf(" \"%s\" встречается - %d\n", w, c)
		}
	}
}

func main() {
	var inputText bytes.Buffer
	s := bufio.NewScanner(os.Stdin)
	const exampleText string = "\"Раз, два, три, даю пробу... Костя, как слышно... три, два, один, прием.\""
	fmt.Println("Введите тест для подсчёта количества одинаковых слов,")
	fmt.Println("завершая ввод новой пустой строкой или \".\"")
	for s.Scan() {
		readLine := s.Text()
		if readLine == "" || readLine == "." {
			break
		}
		inputText.WriteString(readLine + "\n")
	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка чтения:", err)
	}
	if inputText.Len() == 0 {
		fmt.Println("Вы ввели ни слова, " +
			"для демонстрации проводится подсчёт слов в тексте:")
		fmt.Println(exampleText)
		fmt.Fprint(&inputText, exampleText)
	}

	if arr, err := countWords(inputText.String()); err == nil {
		printDups(arr)
	}
}
