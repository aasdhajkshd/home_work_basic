package server

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

const jsonFile = "data/data.json"

type User struct {
	ID       int    `json:"id"`    // уникальный идентификатор
	Name     string `json:"name"`  // имя пользователя
	Email    string `json:"email"` // электронный адрес
	Password string `json:"-"`     // пароль
}

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

type State string

const (
	StateRunning State = "running"
	StateStopped State = "stopped"
)

const (
	timeout time.Duration = 60 * time.Second
)

type Web struct {
	web   http.Server
	state State
}

func NewWeb(ip, port string, handler http.Handler) *Web {
	return &Web{
		web: http.Server{
			Addr:        ip + ":" + port,
			IdleTimeout: timeout,
			Handler:     handler,
		},
		state: StateStopped,
	}
}

func (s *Web) StartWeb() error {
	if s.state == StateRunning {
		return errors.New("web server already running")
	}
	go func() {
		err := s.web.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("web server serving failure: %v", err)
		}
	}()

	s.state = StateRunning

	return nil
}

func HandlerHello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("Привет, пишете письмо?"))
	case http.MethodGet:
		w.Write([]byte("Привет, мой друг!"))
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func handlerUsers(w http.ResponseWriter, r *http.Request) {
	usersJSON, err := json.Marshal(readJSON(jsonFile).Users)
	if err != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("impossible to read all POST body: %s", err)
		}
		// пример для получения данных для создания нового пользователя
		var newUser User
		fmt.Printf("%v", string(data))
		w.Header().Set("Content-Type", "text/plain")
		if err := json.Unmarshal(data, &newUser); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
		} else {
			fmt.Printf("Reveived JSON data, decoded: %v\n", newUser)
			http.Error(w, "OK", http.StatusOK)
		}
		// по-умолчанию выдаём запрощенные данные пользователей без пароля
	default:
		w.Header().Set("Content-Type", "application/json")
		w.Write(usersJSON)
	}
}

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

func readJSON(f string) Data {
	var j Data
	d, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("Read file failed: %v", err)
	}
	if err := json.Unmarshal(d, &j); err != nil {
		log.Fatalf("JSON failed to parse: %v", err)
	}

	return j
}

var (
	address = flag.String("address", "localhost", "HTTP web service address")
	port    = flag.String("port", "8081", "HTTP web service port")
)

func RunServer() {
	flag.Parse()
	log.Println("for interruption, press Ctrl+C...")
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	web := NewWeb(*address, *port, nil)
	http.HandleFunc("/hello", HandlerHello)
	http.HandleFunc("/users", handlerUsers)
	err := web.StartWeb()

	args := flag.Args()
	// для проверки другого варианта использования handler
	if len(args) > 0 {
		addr := strings.Split(args[0], ":")
		r := mux.NewRouter()
		r.HandleFunc("/ehlo", HandlerHello).Methods("GET")
		r.HandleFunc("/ehlo", HandlerHello).Methods("POST")
		web2 := NewWeb(addr[0], addr[1], r)
		err = web2.StartWeb()
	}
	if err != nil {
		panic(err)
	}
	sig := <-done

	fmt.Printf("Received signal: %v\n", sig)

	err = web.StopWeb()
	if err != nil {
		panic(err)
	}

	fmt.Println("Exiting...")
}
