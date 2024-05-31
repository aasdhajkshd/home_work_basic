package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/database"
	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/models"
	"github.com/gorilla/mux"
)

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
)

const (
	timeout           time.Duration = 60 * time.Second
	idleTimeout       time.Duration = 15 * time.Second
	readHeaderTimeout time.Duration = 15 * time.Second
)

const (
	contentTypeJSON           = "application/json"
	contentTypeFormUrlencoded = "application/x-www-form-urlencoded"
	contentTypeFormData       = "multipart/form-data"
)

type Web struct {
	web   http.Server
	state State
}

// Конструктор для создания веб-сервиса.
func NewWeb(ip, port string, handler http.Handler) *Web {
	return &Web{
		web: http.Server{
			Addr:              ip + ":" + port,
			IdleTimeout:       idleTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			Handler:           handler,
		},
		state: StateStopped,
	}
}

// Структуру для хранения зависимостей всего web-сервера.
type Application struct {
	database          *database.PgxPool
	infoLog, errorLog *log.Logger
}

func NewApplication(db *database.PgxPool) *Application {
	return &Application{
		database: db,
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Обработчик запросов для добавления, отображение и удаления пользователей.
func (a *Application) handlerUser(w http.ResponseWriter, r *http.Request) {
	a.infoLog.Printf("Received request: %s %s %s", r.Method, r.URL.Path, r.Header.Get("Content-Type"))
	defer r.Body.Close()
	w.Header().Set("Content-Type", contentTypeJSON)
	user := models.NewUser()

	switch r.Method {
	case http.MethodGet:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id <= 0 {
			// Список пользователей, если id нулевой
			users, err := a.database.ListUsers()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(users)
		} else {
			// Отображение пользователя
			user, err := a.database.FetchUser(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(user)
		}
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case contentTypeJSON:
			err := json.NewDecoder(r.Body).Decode(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unsupported Content-Type Header", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPost {
			// Добавление нового пользователя
			id, err := a.database.AddUser(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})
		}
		if r.Method == http.MethodPut {
			// Обновление существующей заметки
			err := a.database.UpdateUser(user)
			if err != nil {
				http.Error(w, "Failed to update the user", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "User successfully updated"})
		}
		if r.Method == http.MethodDelete {
			// Удаление пользователя
			err := a.database.DeleteUser(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "User successfully deleted"})
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodDelete)
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

// Обработчик запросов на добавление, отображение и удаление продуктов.
func (a *Application) handlerProduct(w http.ResponseWriter, r *http.Request) {
	a.infoLog.Printf("Received request: %s %s %s", r.Method, r.URL.Path, r.Header.Get("Content-Type"))
	if r.Body != nil {
		defer r.Body.Close()
	}
	w.Header().Set("Content-Type", contentTypeJSON)

	switch r.Method {
	case http.MethodGet:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id == 0 {
			// Список всех продуктов
			products, err := a.database.ListProducts()
			if err != nil {
				http.Error(w, "Failed to list products", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(products)
		} else {
			// Отображение продукта
			product, err := a.database.FetchProduct(id)
			if err != nil {
				http.Error(w, "Failed to get product data", http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(product)
		}

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		contentType := r.Header.Get("Content-Type")
		product := models.NewProduct()
		switch contentType {
		case contentTypeJSON:
			err := json.NewDecoder(r.Body).Decode(product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, "Unsupported Content-Type Header", http.StatusBadRequest)
			return
		}
		if r.Method == http.MethodPost {
			// Добавление новой записи по продукту
			id, err := a.database.AddProduct(product)
			if err != nil {
				http.Error(w, "Failed to add a product", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})
		}
		if r.Method == http.MethodPut {
			// Обновление информации по продукту
			err := a.database.UpdateProduct(product)
			if err != nil {
				http.Error(w, "Failed to update note", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product successfully updated"})
		}
		if r.Method == http.MethodDelete {
			// Удаление существующей записи продукта
			err := a.database.DeleteProduct(product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotModified)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Product successfully deleted"})
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Add("Allow", http.MethodPut)
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodDelete)
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

// Обработчик для запросов на добавление и удаления заказа.
func (a *Application) handlerOrder(w http.ResponseWriter, r *http.Request) { //nolint:gocognit
	a.infoLog.Printf("Received request: %s %s %s", r.Method, r.URL.RawPath, r.Header.Get("Content-Type"))
	if r.Body != nil {
		defer r.Body.Close()
	}
	w.Header().Set("Content-Type", contentTypeJSON)

	switch r.Method {
	case http.MethodGet:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		if id > 0 {
			// Информация по заказу
			order, err := a.database.FetchOrder(id)
			if err != nil {
				http.Error(w, "Failed to fetch the order", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(order)
		}
		// Если же есть и user_id, то отобразим заодно и список заказов
		userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
		if userID > 0 { //nolint:nestif
			stat := r.URL.Query().Get("stat")
			if stat != "" {
				// Отображение статистики заказа (общая сумма заказов)
				orders, err := a.database.FetchStatsTotalSumOfOrders(userID)
				if err != nil {
					http.Error(w, "Failed to select the stats", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(orders)
			} else {
				// Отображение заказа
				orders, err := a.database.ListOrdersByUserID(userID)
				if err != nil {
					http.Error(w, "Failed to show the order", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(orders)
			}
		}
	case http.MethodPost, http.MethodDelete:
		contentType := r.Header.Get("Content-Type")
		order := models.NewOrder()
		switch contentType {
		case contentTypeJSON:
			err := json.NewDecoder(r.Body).Decode(order)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if r.Method == http.MethodPost {
				// Добавление нового заказа
				id, err := a.database.AddOrder(order)
				if err != nil {
					http.Error(w, "Failed to add note", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]int{"id": id})
			}
			if r.Method == http.MethodDelete {
				// Удаление существующей заказа
				err := a.database.DeleteOrder(order)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotModified)
					return
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Note successfully deleted"})
			}
		default:
			http.Error(w, "Unsupported Content-Type Header", http.StatusBadRequest)
			return
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodDelete)
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

// Настройки маршрутизации.
func (a *Application) Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/user", a.handlerUser).Methods("POST", "PUT", "GET", "DELETE")
	r.HandleFunc("/product", a.handlerProduct).Methods("POST", "PUT", "GET", "DELETE")
	r.HandleFunc("/order", a.handlerOrder).Methods("POST", "GET", "DELETE")
	return r
}

// Запуск веб-сервера.
func (s *Web) StartWeb() error {
	if s.state == StateRunning {
		return errors.New("web server already running")
	}
	go func() {
		err := s.web.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Web server serving failure: %v", err)
		}
	}()

	s.state = StateRunning

	return nil
}

// Остановка веб-сервера.
func (s *Web) StopWeb() error {
	if s.state == StateStopped {
		return errors.New("web server already stopped")
	}

	err := s.web.Close()
	if err != nil {
		log.Fatalf("Error shutting down the web: %v", err)
	}

	s.state = StateStopped

	return nil
}
