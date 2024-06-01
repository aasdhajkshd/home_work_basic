package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

// Создаёт и проверяет подключение к БД.
func NewPgxPool(dsn string) (*PgxPool, error) {
	ctx := context.Background()

	pgCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("We couldn't find any correct DSN")
	}

	conn, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Cannot connect to database")
	}

	return &PgxPool{ctx: ctx, pool: conn}, err
}

// Закрывает соединение с базой данных.
// Должен вызывается один раз при завершении работы приложения.
func (s *PgxPool) Close() {
	s.pool.Close()
}

func roundToDecimal(n float64, r int) float64 {
	return math.Round(n*math.Pow10(r)) / math.Pow10(r)
}

// Инициализация структуры БД.
func (s *PgxPool) InitDB() error {
	// Создание таблиц
	stmt := `
	SET timezone TO 'UTC';
	CREATE TABLE IF NOT EXISTS users
	(
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(60),
		email VARCHAR(30) NOT NULL,
		password VARCHAR(30) NOT NULL,
		UNIQUE(email)
	);
	CREATE TABLE IF NOT EXISTS orders
	(
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INT NOT NULL,
		order_date TIMESTAMP WITHOUT TIME ZONE NULL,
		total_amount REAL NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(200),
		price REAL
	);
	CREATE TABLE IF NOT EXISTS orderproducts
	(
		id SERIAL NOT NULL PRIMARY KEY,
		order_id INT NOT NULL,
		product_id INT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE NO ACTION
	);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
	CREATE INDEX IF NOT EXISTS idx_products_name ON products (name);
	CREATE INDEX IF NOT EXISTS idx_orderproducts_order_id ON orderproducts (order_id);
	CREATE INDEX IF NOT EXISTS idx_orderproducts_product_id ON orderproducts (product_id);`

	_, err := s.pool.Exec(s.ctx, stmt)
	return err
}

// Инициализация структуры БД.
func (s *PgxPool) DropTables() error {
	// Создание таблиц
	stmt := `
	DROP TABLE IF EXISTS orderproducts CASCADE;
	DROP TABLE IF EXISTS orders CASCADE;
	DROP TABLE IF EXISTS products CASCADE;
	DROP TABLE IF EXISTS users CASCADE;`

	_, err := s.pool.Exec(s.ctx, stmt)
	return err
}

// Добавление пользователя.
func (s *PgxPool) AddUser(user *models.User) (int, error) {
	var id int
	stmt := `
	INSERT INTO users (name, email, password) 
	VALUES 
	($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING id;`

	err := s.pool.QueryRow(s.ctx, stmt, user.Name, user.Email, user.Password).Scan(&id)
	return id, err
}

// Получение данных пользователя.
func (s *PgxPool) FetchUser(id int) (*models.User, error) {
	stmt := `
	SELECT id, name, email, password
	FROM users
	WHERE id = $1;`
	user := models.User{}
	err := s.pool.QueryRow(s.ctx, stmt, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	return &user, err
}

// Получение списка пользователей.
func (s *PgxPool) ListUsers() (*models.Users, error) {
	stmt := `
	SELECT id, name, email, password
	FROM users;`

	rows, err := s.pool.Query(s.ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users models.Users

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, err
		}
		users.List = append(users.List, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, err
}

// Обновление данных пользователя.
func (s *PgxPool) UpdateUser(user *models.User) error {
	// Начало транзации
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}

	// В случае паники, выполнить откат транзакции
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(s.ctx)
		}
	}()

	// Запрос на выборку с блокировкой для обновления
	fUser := models.User{}
	args := []interface{}{user.ID}
	stmt := `
	SELECT id, name, email, password
	FROM users 
	WHERE id = $1
	FOR UPDATE`

	err = tx.QueryRow(s.ctx, stmt, args...).Scan(
		&fUser.ID,
		&fUser.Name,
		&fUser.Email,
		&fUser.Password,
	)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	stmt = `
	UPDATE users SET`
	args = []interface{}{fUser.ID}
	i := 2
	// Обновлений только для тех полей, которые изменились или не являются пустыми
	if user.Name == "" {
		user.Name = fUser.Name
	}
	if user.Name != fUser.Name {
		stmt += fmt.Sprintf(" name = $%d", i)
		args = append(args, user.Name)
		i++
	}
	if user.Email == "" {
		user.Email = fUser.Email
	}
	if user.Email != fUser.Email {
		stmt += fmt.Sprintf(", email = $%d", i)
		args = append(args, user.Email)
		i++
	}
	if user.Password != fUser.Password {
		stmt += fmt.Sprintf(", password = $%d", i)
		args = append(args, user.Password)
	}

	// Обновление данных пользователя
	stmt += ` WHERE id = $1;`
	_, err = tx.Exec(s.ctx, stmt, args...)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	// Завершение транзакции
	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Удаление пользователя.
func (s *PgxPool) DeleteUser(user *models.User) error {
	stmt := `
	DELETE FROM users WHERE id = $1;`

	commandTag, err := s.pool.Exec(s.ctx, stmt, user.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}
	return nil
}

// Добавление продукта.
func (s *PgxPool) AddProduct(product *models.Product) (int, error) {
	var id int
	stmt := `
	INSERT INTO products (name, price) 
	VALUES ($1, $2)
	RETURNING id;`

	err := s.pool.QueryRow(s.ctx, stmt, product.Name, product.Price).Scan(&id)
	return id, err
}

// Получение данных по продукту.
func (s *PgxPool) FetchProduct(id int) (*models.Product, error) {
	var price float64
	stmt := `
	SELECT id, name, price
	FROM products
	WHERE id = $1`
	product := models.Product{}
	err := s.pool.QueryRow(s.ctx, stmt, id).Scan(
		&product.ID,
		&product.Name,
		&price,
	)
	if err != nil {
		return nil, err
	}
	product.Price = roundToDecimal(price, 2)
	return &product, nil
}

// Обновление данных по продукту.
func (s *PgxPool) UpdateProduct(product *models.Product) error {
	// Начало транзации
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}

	// В случае паники, выполнить откат транзакции
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(s.ctx)
		}
	}()

	// Запрос на выборку с блокировкой для обновления
	fProduct := models.Product{}
	args := []interface{}{product.ID}
	stmt := `
	SELECT id, name, price
	FROM products 
	WHERE id = $1
	FOR UPDATE`

	err = tx.QueryRow(s.ctx, stmt, args...).Scan(
		&fProduct.ID,
		&fProduct.Name,
		&fProduct.Price,
	)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	stmt = `
	UPDATE products SET`
	args = []interface{}{fProduct.ID}
	i := 2
	// Обновлений только для тех полей, которые изменились или не являются пустыми
	if product.Name == "" {
		product.Name = fProduct.Name
	}
	if product.Name != fProduct.Name {
		stmt += fmt.Sprintf(" name = $%d", i)
		args = append(args, product.Name)
		i++
	}
	if product.Price != fProduct.Price {
		stmt += fmt.Sprintf(", price = $%d", i)
		product.Price = roundToDecimal(product.Price, 2)
		args = append(args, product.Price)
	}

	// Обновление данных пользователя
	stmt += ` WHERE id = $1;`
	_, err = tx.Exec(s.ctx, stmt, args...)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	// Завершение транзакции
	err = tx.Commit(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Удаление продукта.
func (s *PgxPool) DeleteProduct(product *models.Product) error {
	stmt := `
	DELETE FROM products WHERE id = $1;`

	_, err := s.pool.Exec(s.ctx, stmt, product.ID)
	return err
}

// Получение списка продуктов.
func (s *PgxPool) ListProducts() (*models.Products, error) {
	stmt := `
	SELECT id, name, price
	FROM products;`

	rows, err := s.pool.Query(s.ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products models.Products

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		products.List = append(products.List, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &products, err
}

// Добавление заказа.
func (s *PgxPool) AddOrder(order *models.Order) (int, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return 0, err
	}

	// В случае паники, выполнить откат транзакции
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(s.ctx)
		}
	}()

	var (
		id   int
		stmt = `
		INSERT INTO orders (user_id, order_date, total_amount)
		VALUES ($1, $2, $3)
		RETURNING id;`
		args []interface{}
	)

	// Фиксируется дата заказа
	if order.OrderDate.IsZero() {
		order.OrderDate = time.Now()
	}

	args = []interface{}{order.UserID, order.OrderDate, order.TotalAmount}

	err = tx.QueryRow(s.ctx, stmt, args...).Scan(&id)

	if err != nil {
		tx.Rollback(s.ctx)
		return 0, err
	}

	// Добавление продукта к заказу
	stmt = `
	INSERT INTO orderproducts (order_id, product_id) 
	VALUES ($1, $2);`

	for _, j := range order.Products.List {
		args = []interface{}{id, j.ID}

		_, err = tx.Exec(s.ctx, stmt, args...)
		if err != nil {
			tx.Rollback(s.ctx)
			return 0, err
		}
	}

	// Подтверждаем транзакцию
	err = tx.Commit(s.ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Выборка информации по заказу.
func (s *PgxPool) FetchOrder(id int) (*models.Order, error) {
	// Запрос на выборку записей заметок
	stmt := `
	SELECT 
    o.id AS order_id, 
    COALESCE(o.order_date, '0001-01-01T00:00:00Z00:00') AS order_date,
    o.total_amount AS total_amount,
    o.user_id,
    u.name AS user_name, 
    p.id AS product_id,
    p.name AS product_name,
    COALESCE(p.price, 0.0) AS product_price
	FROM orders o
	LEFT JOIN users u ON u.id = o.user_id
	LEFT JOIN orderproducts ON order_id = o.id
	LEFT JOIN products p ON product_id = p.id
	WHERE o.id = $1;`
	args := []interface{}{id}

	rows, err := s.pool.Query(s.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		order    models.Order
		user     models.User
		products models.Products
	)

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&order.ID,
			&order.OrderDate,
			&order.TotalAmount,
			&user.ID,
			&user.Name,
			&product.ID,
			&product.Name,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}
		products.List = append(products.List, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	order.Products.List = append(order.Products.List, products.List...)
	return &order, err
}

// Выборка информации по заказу.
func (s *PgxPool) ListOrdersByUserID(userID int) (*models.Orders, error) {
	// Запрос на выборку записей заметок
	stmt := `
	SELECT 
    id, user_id, COALESCE(order_date, '0001-01-01T00:00:00Z00:00') AS order_date, total_amount
	FROM orders
	WHERE user_id = $1;`
	args := []interface{}{userID}

	rows, err := s.pool.Query(s.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders models.Orders

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderDate,
			&order.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		orders.List = append(orders.List, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

// Удаление выбранного заказа.
func (s *PgxPool) DeleteOrder(order *models.Order) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}

	// В случае паники, выполнить откат транзакции
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(s.ctx)
		}
	}()
	stmt := `
	DELETE FROM orderproducts WHERE order_id = $1;`

	_, err = tx.Exec(s.ctx, stmt, order.ID)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	stmt = `
	DELETE FROM orders WHERE id = $1;`

	_, err = tx.Exec(s.ctx, stmt, order.ID)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	err = tx.Commit(s.ctx)
	if err != nil {
		tx.Rollback(s.ctx)
		return err
	}

	return nil
}

// Выборка статистики по пользователю (общая сумма заказов).
func (s *PgxPool) FetchStatsTotalSumOfOrders(id int) (float64, error) {
	var (
		name string
		sum  float64
	)

	stmt := `
	SELECT u.name, 
	SUM(p.price * o.total_amount) AS total_price_orders FROM orderproducts op 
	LEFT JOIN products p ON op.product_id=p.id 
	RIGHT JOIN orders o ON op.order_id=o.id, users u 
	WHERE o.user_id = u.id AND u.id = $1 
	GROUP BY u.name;`

	err := s.pool.QueryRow(s.ctx, stmt, id).Scan(&name, &sum)
	return sum, err
}

// Выборка статистики по пользователю (средняя цена товара).
func (s *PgxPool) FetchStatsAveragePrice(id int) (float64, error) {
	var (
		name string
		avg  float64
	)

	stmt := `
	SELECT u.name, 
	ROUND(AVG(p.price)::numeric, 2) AS avg_price_per_orders FROM orderproducts op 
	LEFT JOIN products p ON op.product_id=p.id 
	RIGHT JOIN orders o ON op.order_id=o.id, users u 
	WHERE o.user_id = u.id AND u.id = $1 
	GROUP BY u.name;`

	err := s.pool.QueryRow(s.ctx, stmt, id).Scan(&name, &avg)
	return avg, err
}
