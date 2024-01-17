//nolint:revive,stylecheck
package comparator

const (
	Id CompareType = iota
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

func (t CompareType) String() string {
	switch t {
	case Id:
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

//nolint:exhaustive
func (c Comparator) Compare(bookOne, bookTwo Book) bool {
	switch c.Type {
	case Id:
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

func CompareBooks(bookOne, bookTwo Book, bookValue int) bool {
	publishingDetails := []CompareType{Id, Year, Size, Rate}
	comparator := selectBookType(publishingDetails[bookValue])
	resultOfCompare := comparator.Compare
	return resultOfCompare(bookOne, bookTwo)
}
