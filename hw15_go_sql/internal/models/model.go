package models

import (
	"time"
)

// Структура записей пользователей.
type User struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type Users struct {
	List []User
}

// Конструктор для создания нового пользователя.
func NewUser() *User {
	return &User{}
}

// Структура записей продуктов.
type Product struct {
	ID    int     `db:"id" json:"id"`
	Name  string  `db:"name" json:"name"`
	Price float64 `db:"price" json:"price"`
}

type Products struct {
	List []Product
}

// Конструктор для создания новой продукта с ценой.
func NewProduct() *Product {
	return &Product{}
}

// Структура записей заказов.
type Order struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"userId"`
	OrderDate   time.Time `db:"order_date" json:"orderDate"`
	TotalAmount float64   `db:"total_amount" json:"totalAmount"`
	Products    Products  `db:"products" json:"products"`
}

type Orders struct {
	List []Order
}

// Конструктор для создания новой продукта с ценой.
// Конструктор для создания нового заказа.
func NewOrder() *Order {
	return &Order{
		OrderDate:   time.Now(),
		TotalAmount: 0,
		Products:    Products{List: []Product{}},
	}
}
