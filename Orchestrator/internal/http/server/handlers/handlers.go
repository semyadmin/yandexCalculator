package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
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
			slog.Info(string(data))
			if err != nil {
				slog.Error("Проблема с чтением данных:", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			exp, err := arithmetic.NewPolandNotation(string(data), config)
			if err != nil {
				slog.Error("Ошибка создания выражения:", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			storage.Set(exp, "new")
			slog.Info("Выражение в польской нотации", "result", exp.Expression)
			dataInfo, err := storage.GeByExpression(exp.GetExpression())
			if err != nil {
				slog.Error("Проблема с DataInfo:", "error", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slog.Info("Данные в репозитории:", "dataInfo", dataInfo)
			return
		}
	}
}
