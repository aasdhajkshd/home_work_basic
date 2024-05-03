package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// Пример структуры API для клиента БД из следующих уроков.
type User struct {
	ID       int    `json:"id"`    // уникальный идентификатор
	Name     string `json:"name"`  // имя пользователя
	Email    string `json:"email"` // электронный адрес
	Password string `json:"-"`     // пароль
}

//nolint:tagliatelle
type Order struct {
	ID          int       `json:"id"`           // идентификатор заказа
	UserID      int       `json:"user_id"`      // идентификатор пользователя
	OrderDate   time.Time `json:"order_date"`   // дата заказа
	TotalAmount float32   `json:"total_amount"` // общяя стоимость заказов
	Products    []int     `json:"products"`
}

type Product struct {
	ID    int     `json:"id"`    // идентификатор товара
	Name  string  `json:"name"`  // название
	Price float32 `json:"price"` // цена
}

type Data struct {
	Users    []User    `json:"users"`
	Orders   []Order   `json:"orders"`
	Products []Product `json:"products"`
}

func ReadJSON(f string) Data {
	var jsonData Data
	data, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Read file failed: %v", err)
	}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatalf("JSON failed to parse: %v", err)
	}
	return jsonData
}

// Пример использования данных.
func PrintData(f string) {
	data := ReadJSON(f)

	fmt.Println("Пользователи:")
	for _, user := range data.Users {
		fmt.Printf("ID: %d, Имя: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}

	fmt.Println("\nЗаказы:")
	for _, order := range data.Orders {
		fmt.Printf("ID заказа: %d, ID пользователя: %d, Дата заказа: %s, Сумма заказа: %.2f"+
			"\n", order.ID, order.UserID, order.OrderDate, order.TotalAmount)
		fmt.Printf("Товары в заказе: ")
		for _, productID := range order.Products {
			fmt.Printf("%d ", productID)
		}
		fmt.Println()
	}

	fmt.Println("\nТовары:")
	for _, product := range data.Products {
		fmt.Printf("ID товара: %d, Название: %s, Цена: %.2f\n", product.ID, product.Name, product.Price)
	}
}
