package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/database"
	"github.com/aasdhajkshd/home_work_basic/hw15_go_sql/internal/web"
)

const (
	dsn string = "postgres://manager:password@localhost:5432/onlinestore?sslmode=disable&pool_max_conns=100"
)

var (
	address = flag.String("address", "localhost", "HTTP web service address")
	port    = flag.String("port", "8888", "HTTP web service port")
	verbose = flag.Bool("verbose", false, "Verbose output")
)

// Запуск приложения, включая веб-сервис.
func RunServer() {
	flag.Parse()
	log.Println("for interruption, press Ctrl+C...")
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Создается новое подключение к ДБ
	db, err := database.NewPgxPool(dsn)
	if err != nil {
		log.Fatalf("Failed to create PgxPool: %v", err)
	}
	defer db.Close()

	// Инициируем связку приложения с БД.
	app := web.NewApplication(db)

	// Передаем handler (nil - default) в настройки веб-сервиса.
	web := web.NewWeb(*address, *port, app.Router())

	err = web.StartWeb()
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

func main() {
	flag.Parse()

	if *verbose {
		for i, j := range flag.Args() {
			print("parsed arguments:", i, j)
		}
	}

	print("Starting 'server'...\n")
	RunServer()
}
