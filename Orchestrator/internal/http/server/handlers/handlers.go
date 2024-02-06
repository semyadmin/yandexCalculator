package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
)

type Handlers struct{}

func NewServeMux() (http.Handler, error) {
	// Создам маршрутизатор
	serveMux := http.NewServeMux()
	// Регистрируем обработчики событий
	serveMux.HandleFunc("/hello", helloHandler)
	serveMux.HandleFunc("/expression", expressionHandler)
	return serveMux, nil
}

func Decorate(next http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	decorated := next

	for i := len(middleware) - 1; i >= 0; i-- {
		decorated = middleware[i](decorated)
	}

	return decorated
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func expressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		type Expression interface {
			Result() []string
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		exp, err := arithmetic.NewPolandNotation(string(data))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Info("Result of expression by poland notation", "result", exp.Result())
		return
	}
}
