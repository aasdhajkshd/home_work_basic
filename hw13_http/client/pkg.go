package client

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// для тестирования отправки POST запроса
func AddNewUser(id int, name, email, password string) ([]byte, error) {
	user := struct {
		ID       int    `json:"id"`       // уникальный идентификатор
		Name     string `json:"name"`     // имя пользователя
		Email    string `json:"email"`    // электронный адрес
		Password string `json:"password"` // пароль
	}{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	return json.Marshal(user)
}

var (
	method, url, path, content string
	timeout                    time.Duration
)

func init() {
	flag.StringVar(&method, "method", "GET", "GET or POST HTTP methods")
	flag.StringVar(&url, "url", "https://phet-dev.colorado.edu/", "HTTP URL address")
	flag.StringVar(&path, "path", "html/build-an-atom/0.0.0-3/simple-text-only-test-page.html", "HTTP URL path arguments")
	flag.StringVar(&content, "content", "application/json", "HTTP POST content")
	flag.DurationVar(&timeout, "timeout", 5 * time.Second, "connection timeout")
}

func RunClient() {
	flag.Parse()
	args := flag.Args()

	// для тестирования (упрощения) указания адреса и/или порта как аргумент
	if len(args) > 0 {
		url = args[0]
	}

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	if path != "" {
		url = url + "/" + path
	}

	resp, _ := RequestURL(method, url, content, timeout)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read body error:", err)
	}

	fmt.Printf("%s:%d\n", resp.Header, len(body))
	fmt.Println(string(body))
}

func RequestURL(method, url, content string, timeout time.Duration) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
		user []byte
	)

	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	switch strings.ToUpper(method) {
	case http.MethodGet:
		resp, err = client.Get(url)
		if err != nil {
			log.Println("Request GET error:", err)
		}
	case http.MethodPost:
		user, err = AddNewUser(4, "Dmitry Dimovich", "dd@mail.ru", "11111111")
		if err != nil {
			log.Fatalf("Impossible to marshall user data: %s", err)
		}
		resp, err = client.Post(url, content, bytes.NewReader(user))
		if err != nil {
			log.Println("Request POST error:", err)
		}
	default:
		log.Fatalf("Method %s is not supported", method)
	}
	return resp, err
}
