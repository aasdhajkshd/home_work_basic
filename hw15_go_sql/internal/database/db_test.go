package database

//nolint:gofumpt,gci,nolintlint
import (
	"testing"
	"time"

	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	dsnURI          = "postgres://manager:password@postgresql:5432/onlinestore_test?sslmode=disable&pool_max_conns=100"
	userEmail       = "johndoe@example.com"
	userNewEmail    = "johnydoe@example.ru"
	userName        = "John Doe"
	userNewName     = "Johny Doe"
	userPass        = "Password"
	userNewPass     = "Drowssap"
	productName     = "Telephone"
	productNewName  = "Reduziertes Telefon"
	productPrice    = 10.11
	productNewPrice = 20.02
)

func TestConnection(t *testing.T) {
	// Подключение к тестовой базе данных (можно использовать временную или тестовую базу данных)
	db, err := NewPgxPool(dsnURI)
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	// Простой запрос
	_, err = db.pool.Exec(db.ctx, "SELECT 1 AS id;")
	assert.NoError(t, err, "Failed to make SELECT statement")
}

func TestInitDB(t *testing.T) {
	// Подключение к тестовой базе данных (можно использовать временную или тестовую базу данных)
	db, err := NewPgxPool(dsnURI)
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}
	defer db.pool.Close()

	err = db.InitDB()
	assert.NoError(t, err, "Failed to init test DB")
}

func TestUser(t *testing.T) {
	// Открытие соединения с базой
	db, err := NewPgxPool(dsnURI)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	defer db.pool.Close()

	// Создание временной таблицы пользователей в рамках теста
	stmt := `
	CREATE TABLE IF NOT EXISTS users
	(
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(60),
		email VARCHAR(30) NOT NULL,
		password VARCHAR(30) NOT NULL,
		UNIQUE(email)
	);
	DELETE FROM users;`

	_, err = db.pool.Exec(db.ctx, stmt)
	if err != nil {
		t.Fatalf("Error setting up test user table: %v", err)
	}

	testUser := models.NewUser()
	testUser.Name = userName
	testUser.Email = userEmail
	testUser.Password = userPass

	t.Run("TestAddUser", func(t *testing.T) {
		id, err := db.AddUser(testUser)
		testUser.ID = id
		assert.NoError(t, err, "Failed to add user")

		// Проверяем, что ID пользователя больше нуля (успешное добавление)
		assert.NotZero(t, testUser.ID, "Invalid user ID")
	})

	t.Run("TestFetchUser", func(t *testing.T) {
		fetchedUser, err := db.FetchUser(testUser.ID)
		assert.NoError(t, err, "Failed to fetch user")
		assert.Equal(t, testUser.Name, fetchedUser.Name, "User name mismatch")
		assert.Equal(t, testUser.Email, fetchedUser.Email, "User email mismatch")
	})

	t.Run("TestListUsers", func(t *testing.T) {
		fetchedUsers, err := db.ListUsers()
		assert.NoError(t, err, "Failed to list users")
		assert.Len(t, fetchedUsers.List, 1, "No users fetched")
	})

	t.Run("TestUpdateUser", func(t *testing.T) {
		testUser.Name = userNewName
		testUser.Email = userNewEmail
		err := db.UpdateUser(testUser)
		assert.NoError(t, err, "Failed to update a user info")
		fetchedUser, _ := db.FetchUser(testUser.ID)
		assert.Equal(t, userNewName, fetchedUser.Name, "User name mismatch")
		assert.Equal(t, userNewEmail, fetchedUser.Email, "User email mismatch")
	})

	t.Run("TestDeleteUser", func(t *testing.T) {
		err := db.DeleteUser(testUser)
		assert.NoError(t, err, "Failed to delete user")
	})
}

func TestProduct(t *testing.T) {
	// Открытие соединения с базой
	db, err := NewPgxPool(dsnURI)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	defer db.pool.Close()

	// Создание временной таблицы пользователей в рамках теста
	stmt := `
	CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL NOT NULL PRIMARY KEY,
		name VARCHAR(200),
		price REAL
	);
	DELETE FROM products;`

	_, err = db.pool.Exec(db.ctx, stmt)
	if err != nil {
		t.Fatalf("Error setting up test user table: %v", err)
	}

	testProduct := models.NewProduct()
	testProduct.Name = productName
	testProduct.Price = productPrice

	t.Run("TestAddProduct", func(t *testing.T) {
		id, err := db.AddProduct(testProduct)
		testProduct.ID = id
		assert.NoError(t, err, "Failed to add a product")
		assert.NotZero(t, testProduct.ID, "Invalid product ID")
	})

	t.Run("TestFetchProduct", func(t *testing.T) {
		fetchedProduct, err := db.FetchProduct(testProduct.ID)
		assert.NoError(t, err, "Failed to fetch product")
		assert.Equal(t, testProduct.Name, fetchedProduct.Name, "product name mismatch")
		assert.Equal(t, testProduct.Price, fetchedProduct.Price, "Product price mismatch")
	})

	t.Run("TestUpdateProduct", func(t *testing.T) {
		testProduct.Name = productNewName
		testProduct.Price = productNewPrice
		err := db.UpdateProduct(testProduct)
		assert.NoError(t, err, "Failed to update a product")
		fetchedProduct, _ := db.FetchProduct(testProduct.ID)
		assert.Equal(t, productNewName, fetchedProduct.Name, "product name mismatch")
		assert.Equal(t, productNewPrice, fetchedProduct.Price, "Product price mismatch")
	})

	t.Run("TestListProducts", func(t *testing.T) {
		fetchedProducts, err := db.ListProducts()
		assert.NoError(t, err, "Failed to list users")
		assert.Len(t, fetchedProducts.List, 1, "No products fetched")
	})

	t.Run("TestDeleteProduct", func(t *testing.T) {
		err := db.DeleteProduct(testProduct)
		assert.NoError(t, err, "Failed to delete product")
	})
}

func TestOrder(t *testing.T) {
	// Открытие соединения с базой
	db, err := NewPgxPool(dsnURI)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	defer db.pool.Close()

	// Создание временной таблицы пользователей в рамках теста
	stmt := `
	CREATE TABLE IF NOT EXISTS orders
	(
		id SERIAL NOT NULL PRIMARY KEY,
		user_id INT NOT NULL,
		order_date TIMESTAMP WITHOUT TIME ZONE NULL,
		total_amount REAL NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS orderproducts
	(
		id SERIAL NOT NULL PRIMARY KEY,
		order_id INT NOT NULL,
		product_id INT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE NO ACTION
	);
	DELETE FROM orderproducts;
	DELETE FROM orders;
	DELETE FROM products;
	DELETE FROM users;`

	_, err = db.pool.Exec(db.ctx, stmt)
	if err != nil {
		t.Fatalf("Error setting up test user table: %v", err)
	}

	testUser := *models.NewUser()
	testUser.Name = userName
	testUser.Email = userEmail
	testUser.Password = userPass

	t.Run("TestAddUser", func(t *testing.T) {
		id, _ := db.AddUser(&testUser)
		testUser.ID = id
		assert.NotZero(t, testUser.ID, "Invalid user ID")
	})

	testProduct := *models.NewProduct()
	testProduct.Name = productName
	testProduct.Price = 123.10

	t.Run("TestAddProduct", func(t *testing.T) {
		id, _ := db.AddProduct(&testProduct)
		testProduct.ID = id
		assert.NotZero(t, testProduct.ID, "Invalid product ID")
	})

	testOrder := *models.NewOrder()
	testOrder.UserID = testUser.ID
	testOrder.OrderDate = time.Now()
	testOrder.TotalAmount = 123.21

	testOrder.Products.List = append(testOrder.Products.List, testProduct)

	t.Run("TestAddOrder", func(t *testing.T) {
		id, err := db.AddOrder(&testOrder)
		testOrder.ID = id
		assert.NoError(t, err, "Failed to create an order")
		assert.NotZero(t, testOrder.ID, "Invalid order ID")
	})

	t.Run("TestFetchStatsTotalSumOfOrders", func(t *testing.T) {
		sum, err := db.FetchStatsTotalSumOfOrders(testUser.ID)
		assert.NoError(t, err, "Failed to get stats of total sum")
		assert.NotZero(t, sum, "Sum cannot be zero")
	})
	// Смысла в этом запросе нет, так как нужно насыщение различными данными
	t.Run("TestFetchStatsAveragePrice", func(t *testing.T) {
		_, err := db.FetchStatsAveragePrice(testUser.ID)
		assert.NoError(t, err, "Failed to get average price")
	})

	t.Run("TestDeleteOrder", func(t *testing.T) {
		err := db.DeleteOrder(&testOrder)
		assert.NoError(t, err, "Failed to delete the order")
	})
}
