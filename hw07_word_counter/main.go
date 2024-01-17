package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func removePunct(s string) string {
	return strings.TrimFunc(s, unicode.IsPunct)
}

func countWords(x string) (map[string]int, error) { // по условаию задания принимаем строку
	wordsMap := make(map[string]int)
	s := bufio.NewScanner(strings.NewReader(x))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		w := removePunct(s.Text())
		wordsMap[strings.ToLower(w)]++
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
	var inputText []string
	s := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите тест для подсчёта количества одинаковых слов,")
	fmt.Println("завершая ввод новой пустой строкой или \".\"")
	for s.Scan() {
		readLine := s.Text()
		if readLine == "" || readLine == "." {
			break
		}
		inputText = append(inputText, readLine)
	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка чтения:", err)
	}
	if len(inputText) == 0 {
		fmt.Println("Вы ввели ни слова, " +
			"для демонстрации проводится подсчёт слов в тексте:")
		inputText = []string{
			"Раз, два, три, даю пробу...",
			"Костя, как слышно...",
			"Три, два, один, прием.",
		}

		fmt.Println(inputText)
	}

	if arr, err := countWords(strings.Join(inputText, "\n")); err == nil {
		printDups(arr)
	}
}
