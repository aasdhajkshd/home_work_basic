package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	book "github.com/aasdhajkshd/home_work_basic/hw09_serialize/book"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Book struct {
	Id      uint64    `json:"id"` //nolint:revive,stylecheck
	Year    rune      `json:"year"`
	Size    rune      `json:"size"`
	Rate    float32   `json:"rate,omitempty"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Updated time.Time `json:"updated,omitempty"`
}

type Message interface {
	MarshalProto() ([]byte, error)
}

func (b *Book) MarshalProto() ([]byte, error) {
	bookProto := &book.Message{
		Id:      b.Id,
		Year:    b.Year,
		Size:    b.Size,
		Rate:    b.Rate,
		Title:   b.Title,
		Author:  b.Author,
		Updated: timestamppb.New(b.Updated),
	}
	return proto.Marshal(bookProto)
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

func (b *Book) MarshalJSON() ([]byte, error) {
	type bookList Book
	return json.Marshal(&struct {
		Updated int64 `json:"updated"`
		*bookList
	}{
		Updated:  b.Updated.Unix(),
		bookList: (*bookList)(b),
	})
}

func (b *Book) UnmarshalJSON(data []byte) error {
	type bookList Book
	aux := &struct {
		Updated int64 `json:"updated"`
		*bookList
	}{
		bookList: (*bookList)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	b.Updated = time.Unix(aux.Updated, 0)
	return nil
}

func main() {
	nowTime := time.Now()

	books := []Book{
		{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden", nowTime},
		{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz", nowTime},
	}

	fmt.Println("*** Marshaling with json package ***")
	data, err := json.MarshalIndent(&books, "", "  ")
	if err != nil {
		fmt.Printf("JSON marshalling error: %v\n", err)
		log.Fatalln()
	}
	fmt.Printf("Books in json:\n%s\n", data)

	var bookData []Book
	err = json.Unmarshal(data, &bookData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unmarshalled json:")
	for i := range bookData {
		fmt.Printf("%v\n", bookData[i])
	}

	// Protobuf
	bookProto := &Book{
		Id:      1484296656,
		Year:    2023,
		Size:    384,
		Rate:    2.0,
		Title:   "Go Crazy",
		Author:  "Nicolas Modrzyk",
		Updated: time.Now(),
	}

	fmt.Println()
	fmt.Println("*** Marshaling with proto package ***")
	res, err := bookProto.MarshalProto()
	if err != nil {
		fmt.Printf("ProtoBuf marshalling error: %v\n", err)
		log.Fatalln()
	}

	filename := "book.data"

	if err = os.WriteFile(filename, res, 0644); err != nil { //nolint:gofumpt,gosec
		log.Fatalln("Failed to write book:", err)
	}
	fmt.Printf("Book in bytes %v written to file \"%s\"\n\n", len(res), filename)

	res, err = os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error: %v'n", err)
	}

	bookBlob := &book.Message{}
	err = proto.Unmarshal(res, bookBlob)

	if err != nil {
		fmt.Printf("error: %v'n", err)
	}
	fmt.Printf("Unmarshalled protobuf:\n%+v\n", bookBlob)
}
