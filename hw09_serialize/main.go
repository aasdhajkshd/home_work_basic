package main

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	_ "io"
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
	bookToMessage() *book.Message
}

func (b Book) bookToMessage() *book.Message {
	return &book.Message{
		Id:      b.Id,
		Year:    b.Year,
		Size:    b.Size,
		Rate:    b.Rate,
		Title:   b.Title,
		Author:  b.Author,
		Updated: timestamppb.New(b.Updated),
	}
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

func (b *Book) MarshalJSON() ([]byte, error) {
	type bookAlias Book
	return json.Marshal(&struct {
		Updated int64 `json:"updated"`
		*bookAlias
	}{
		Updated:   b.Updated.Unix(),
		bookAlias: (*bookAlias)(b),
	})
}

func (b *Book) UnmarshalJSON(data []byte) error {
	type bookAlias Book
	aux := &struct {
		Updated int64 `json:"updated"`
		*bookAlias
	}{
		bookAlias: (*bookAlias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	b.Updated = time.Unix(aux.Updated, 0)
	return nil
}

func SerializeBooksJSON(books []Book) ([]byte, error) {
	return json.Marshal(books)
}

func DeserializeBooksJSON(data []byte) ([]Book, error) {
	var books []Book
	err := json.Unmarshal(data, &books)
	return books, err
}

func bookToFile(filename string, data []byte) error {
	if err := os.WriteFile(filename, data, 0644); err != nil { //nolint:gofumpt,gosec
		log.Fatalln("Failed to write data to book:", err)
		return err
	}
	return nil
}

func bookFromFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func serializeBooksProto(books []Book) ([]byte, error) {
	protoBooks := make([]*book.Message, 0, len(books))
	for _, j := range books {
		protoBooks = append(protoBooks, j.bookToMessage())
	}
	return proto.Marshal(&book.List{Books: protoBooks})
}

func main() {
	nowTime := time.Now()

	books := []Book{
		{1250237238, 2019, 263, 4.7, "Permanent Record", "Edward Snowden", nowTime},
		{1718501129, 2021, 216, 4.6, "Black Hat Python 2", "Justin Seitz", nowTime},
	}
	// Marshalling with json
	fmt.Println("*** Marshaling with json package ***")
	data, err := json.Marshal(&books)
	if err != nil {
		fmt.Printf("JSON marshalling error: %v\n", err)
		log.Fatalln()
	}
	fmt.Printf("%s\n", data)

	var bookData []Book
	err = json.Unmarshal(data, &bookData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Unmarshalled json:")
	for i := range bookData {
		fmt.Printf("%v\n", bookData[i])
	}

	// Serialization with json package via functions
	fmt.Println()
	fmt.Println("*** Serialization with json ***")
	serializedDataJSON, err := SerializeBooksJSON(books)
	if err != nil {
		log.Fatalf("Error serialization: %v\n", err)
	}

	fmt.Println()
	deserializedData, err := DeserializeBooksJSON(serializedDataJSON)
	if err != nil {
		log.Fatalf("Error deserialization: %v\n", err)
	}
	for i := range deserializedData {
		fmt.Printf("JSON in slice Book %d: %+v\n", i, deserializedData[i])
	}

	// Protobuf

	bookProto := Book{
		Id:      1484296656,
		Year:    2023,
		Size:    384,
		Rate:    2.0,
		Title:   "Go Crazy",
		Author:  "Nicolas Modrzyk",
		Updated: time.Now(),
	}.bookToMessage()

	fmt.Println()
	fmt.Println("*** Marshalling with protobuf ***")
	data, err = proto.Marshal(bookProto)
	if err != nil {
		fmt.Printf("ProtoBuf marshalling error: %v\n", err)
		log.Fatalln()
	}
	fmt.Printf("%s\n", data)

	// Writing to file bytes
	filename := "books.data"
	bookToFile(filename, data)

	// Reading from file
	data, err = bookFromFile(filename)
	if err != nil {
		log.Fatalln("Failed to read data from file:", err)
	}
	bookMessage := &book.Message{}
	err = proto.Unmarshal(data, bookMessage)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Unmarshalled protobuf:\n%+v\n", bookMessage)

	// Serialization protobuf slice
	fmt.Println()
	fmt.Println("*** Serialization with protobuf books slice ***")
	serializedDataProto, err := serializeBooksProto(books)
	if err != nil {
		log.Fatalf("Error serialization: %v\n", err)
	}
	fmt.Printf("%s", serializedDataProto)

	// Deserialization protobuf slice
	fmt.Println("*** Deserialization of proto books slice ***")
	bookList := &book.List{}
	if err = proto.Unmarshal(serializedDataProto, bookList); err != nil {
		log.Fatalf("Error deserialization: %v\n", err)
	}

	// Output of deserialized protobuf books slice
	for _, b := range bookList.Books {
		fmt.Printf("ID: %d, Year: %d, Size: %d, Rate: %.2f, Title: %s, Author: %s, Updated: %v\n",
			b.GetId(), b.GetYear(), b.GetSize(), b.GetRate(), b.GetTitle(), b.GetAuthor(), b.GetUpdated())
	}

	fmt.Println()
	fmt.Println("*** Results ***")
	data, _ = json.MarshalIndent(&books, "", "  ")
	fmt.Printf("Books:\n%s\n", data)
	fmt.Printf("JSON books in bytes:  \"%d\"\n", len(serializedDataJSON))
	fmt.Printf("Protobuf books in bytes: \"%d\"\n", len(serializedDataProto))
}
