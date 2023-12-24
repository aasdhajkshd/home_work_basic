package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ID CompareType = iota
	Year
	Size
	Rate
	Title
	Author
)

type Book struct {
	id            uint64
	year, size    rune
	rate          float32
	title, author string
}

type CompareType uint8

type Comparator struct {
	Type CompareType
}

var path = "books.json"

func (b *Book) Id() uint64 {
	return b.id
}

func (b *Book) Year() rune {
	return b.year
}

func (b *Book) Size() rune {
	return b.size
}

func (b *Book) Rate() float32 {
	return b.rate
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) Author() string {
	return b.author
}

func (b *Book) SetId(id uint64) {
	b.id = id
}

func (b *Book) SetYear(year rune) {
	b.year = year
}

func (b *Book) SetSize(size rune) {
	b.year = size
}

func (b *Book) SetRate(rate float32) {
	b.rate = rate
}

func (b *Book) SetTitle(title string) {
	b.title = title
}

func (b *Book) SetAuthor(author string) {
	b.author = author
}

func flushBuffers() {
	bufio.NewScanner(os.Stdin).Scan() // flush input buffer in case of errored fmt.Fscanf
}

func YesNo(question string) bool {
	fmt.Printf("%s [y/n]: ", question)
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	fmt.Println(s.Text())
	return strings.ToLower(strings.TrimSpace(s.Text())) == "y"
}

func (t CompareType) String() string {
	switch t {
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

func selectBookType(t CompareType) *Comparator {
	return &Comparator{
		Type: t,
	}
}

func (c Comparator) Compare(bookOne, bookTwo Book) bool {
	switch c.Type {
	case ID:
		return bookOne.Id() > bookTwo.Id()
	case Year:
		return bookOne.Year() > bookTwo.Year()
	case Size:
		return bookOne.Size() > bookTwo.Size()
	case Rate:
		return bookOne.Rate() > bookTwo.Rate()
	default:
		return false
	}
}

func (b *Book) populateBook() error {
	if YesNo("Желаете добавить книги одной строкой:") {
		fmt.Println("Укажите информацию через запятую,")
		fmt.Print("(Номер ISBN, Год издания, Кол-во страниц, Рейтинг, \"Название\", \"Автор\"): ")
		_, e := fmt.Fscanf(os.Stdin, "%d, %d, %d, %f, %q, %q", &b.id, &b.year, &b.size, &b.rate, &b.title, &b.author)
		flushBuffers()
		fmt.Println(e)
		if e != nil {
			fmt.Println("Ошибка обработки данных")
			return e
		}
	} else {
		fmt.Print("Введите информацию о книге:\n")
		fmt.Print("Номер ISBN: ")
		fmt.Scanln(&b.id)
		fmt.Print("Год издания: ")
		fmt.Scanln(&b.year)
		fmt.Print("Количество страниц: ")
		fmt.Scanln(&b.size)
		fmt.Print("Рейтинг: ")
		fmt.Scanln(&b.rate)
		fmt.Print("Название: ")
		fmt.Scanln(&b.title)
		fmt.Print("Автор: ")
		fmt.Scanln(&b.author)
	}
	return nil
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
	bookOne, bookTwo := Book{}, Book{}
	if bookOne.populateBook() == nil && bookTwo.populateBook() == nil {
		fmt.Println("Вывод книг для проверки:")
		fmt.Printf("bookOne: %+v\n", bookOne)
		fmt.Printf("bookTwo: %+v\n", bookTwo)
	} else {
		fmt.Println("Выпал шанс ввести вручную, произошла ошибка, придётся устанавливать предопределенные значения:")
		bookOne = Book{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden"}
		bookTwo = Book{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz"}
		books := []Book{bookOne, bookTwo}
		for _, book := range books {
			fmt.Printf("%+v\n", book)
		}
	}
	var bookValue uint8
	fmt.Println("\nПо какому полю выполнить сравнение, доступные варианты: ")
	publishingDetails := []CompareType{ID, Year, Size, Rate}
	for bookIndex, bookType := range publishingDetails {
		fmt.Print("[", bookIndex, "] - ", bookType.String(), "\n")
	}
	fmt.Fscanln(os.Stdin, &bookValue)
	fmt.Println("Вы выбрали:", publishingDetails[bookValue])
	comparator := selectBookType(publishingDetails[bookValue])
	resultOfCompare := comparator.Compare
	fmt.Printf("Сравнивая \"%s\", у книги \"%s\" больше, чем у \"%s\": \n", comparator.Type.String(), bookOne.title, bookTwo.title)
	fmt.Println("===================")
	fmt.Println("Результат: ", resultOfCompare(bookOne, bookTwo))
	fmt.Println("===================")
	os.Exit(0)
	saveJSON(path)
}
