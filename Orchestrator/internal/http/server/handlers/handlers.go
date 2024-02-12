package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
)

func NewServeMux(config *config.Config, queue *queue.MapQueue, storage *memory.Storage) (http.Handler, error) {
	// Создам маршрутизатор
	serveMux := http.NewServeMux()
	// Регистрируем обработчики событий
	patchToFront := "./frontend/build"
	serveMux.Handle("/", http.FileServer(http.Dir(patchToFront)))
	serveMux.HandleFunc("/expression", expressionHandler(config, queue, storage))
	return serveMux, nil
}

func Decorate(next http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	decorated := next
	for i := len(middleware) - 1; i >= 0; i-- {
		decorated = middleware[i](decorated)
	}

	return decorated
}

func expressionHandler(config *config.Config, queue *queue.MapQueue, storage *memory.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type Expression interface {
				Result() []string
			}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Проблема с чтением данных:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slog.Info("Полученное выражение от пользователя:", "выражение:", string(data))
			str, ok := validator.Validator(string(data))
			if !ok {
				slog.Error("Некорректное выражение:", "ошибка:", err)
				http.Error(w, "Ваше выражение "+str+" некорректное", http.StatusBadRequest)
				return
			}
			dataInfo, err := storage.GeByExpression(str)
			if err == nil {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("Ваш id: " + strconv.FormatInt(int64(dataInfo.Id), 10)))
			}
			exp, err := arithmetic.NewASTTree(str, config, queue)
			if err != nil {
				slog.Error("Проблема с вычислением выражения:", "выражение:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			storage.Set(exp, "new")
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("Ваш id: " + strconv.FormatInt(int64(exp.ID), 10)))
			return
		}
	}
}
