package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Book struct {
	ID            uint64
	Year, Size    rune
	Rate          float32
	Title, Author string
}

// var book1 = Book{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden"}
// var book2 = Book{1594206015, 2020, 340, 4.5, "Dark Mirror", "Edward Snowden"}
// var book3 = Book{1593275900, 2014, 192, 4.6, "Black Hat Python", "Justin Seitz"}
// var book4 = Book{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz"}
var (
	books = []Book{}
	path  = "books.json"
)

type BookData uint8

const (
	ID BookData = iota
	Year
	Size
	Rate
	Title
	Author
)

func flushBuffers() {
	bufio.NewScanner(os.Stdin).Scan() // flush input buffer in case of errored fmt.Fscanf
}

func yesno(question string) bool {
	fmt.Printf("%s [y/n]: ", question)
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	fmt.Println(s.Text())
	return strings.ToLower(strings.TrimSpace(s.Text())) == "y"
}

func (b BookData) String() string {
	switch b {
	case ID:
		return "Номер ISBN"
	case Year:
		return "Год издания"
	case Size:
		return "Количество страниц"
	case Rate:
		return "Рейтинг"
	case Title:
		return "Название"
	case Author:
		return "Автор"
	default:
		return "Неизвестное поле"
	}
}

func (b BookData) StringRU() string {
	BookInfo := [...]string{"Номер ISBN", "Год издания", "Количество страниц", "Рейтинг", "Название", "Автор"}[b]
	return BookInfo
}

func (b BookData) EnumIndex() uint8 {
	return uint8(b)
}

func getBookDetail() uint8 {
	var bookValue uint8
	fmt.Println("\nПо какому полю выполнить сравнение, доступные варианты: ")
	publishingDetails := []BookData{ID, Year, Size, Rate}
	for bookIndex, bookDetail := range publishingDetails {
		fmt.Print("[", bookIndex, "] - ", bookDetail.String(), "\n")
	}
	fmt.Fscanln(os.Stdin, &bookValue)
	fmt.Println("Вы выбрали:", bookValue)
	return bookValue
}

func getBooks() ([]Book, error) {
	var book Book
	fmt.Print("(Номер ISBN, Год издания, Кол-во страниц, Рейтинг, \"Название\", \"Автор\"): ")
	_, e := fmt.Fscanf(os.Stdin, "%d, %d, %d, %f, %q, %q", &book.ID, &book.Year, &book.Size, &book.Rate, &book.Title, &book.Author)
	if e != nil {
		fmt.Println("Ошибка при сканировании:", e)
	} else {
		books = append(books, book)
		fmt.Println(books)
	}
	flushBuffers()
	return books, e
}

func fileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return !os.IsNotExist(e)
}

func saveJSON(filePath string) {
	if !fileExists(path) {
		fmt.Println("Файл", path, "не существует, можем записать... ")
	}
}

func main() {
	for i := 0; i < 6; i++ {
		fmt.Println("Укажите информацию по книге через запятую:")
		books, e := getBooks()
		if e != nil {
			fmt.Println("Ошибка при получении книг")
			if yesno("Попробовать еще раз:") {
				fmt.Println("continue")
				continue
			} else {
				fmt.Println("break")
				break
			}
		}
		if len(books) == 1 {
			fmt.Println("Для сравнения нужна вторая книга...")
		}
		if len(books) == 2 {
			fmt.Println("\nСписок книг для сравнения:")
			for _, book := range books {
				fmt.Printf("%+v\n", book)
			}
			selectedBookIndex := getBookDetail()
			fmt.Print("Сравниваем \"", BookData(selectedBookIndex), "\": ")
			switch selectedBookIndex {
			case 0:
				fmt.Println(books[0].Title, ":", books[0].ID, "больше", books[1].Title, ":", books[1].ID, "=", books[0].ID > books[1].ID)
			case 1:
				fmt.Println(books[0].Title, ":", books[0].Year, "больше", books[1].Title, ":", books[1].Year, "=", books[0].Year > books[1].Year)
			case 2:
				fmt.Println(books[0].Title, ":", books[0].Size, "больше", books[1].Title, ":", books[1].Size, "=", books[0].Size > books[1].Size)
			case 3:
				fmt.Println(books[0].Title, ":", books[0].Rate, "больше", books[1].Title, ":", books[1].Rate, "=", books[0].Rate > books[1].Rate)
			}
			break
		}
	}
	saveJSON(path)
}
