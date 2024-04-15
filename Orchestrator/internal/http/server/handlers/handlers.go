package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/duration"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/expression"
	newexpression "github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/expression"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/services/user"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

func NewServeMux(config *config.Config,
	queue *queue.MapQueue,
	storage *memory.Storage,
	userStorage *memory.UserStorage,
) (http.Handler, error) {
	// Создам маршрутизатор
	serveMux := http.NewServeMux()
	// Регистрируем обработчики событий
	patchToFront := "./frontend/build"
	// Страница статики для фронтенда
	serveMux.Handle("/", http.FileServer(http.Dir(patchToFront)))
	// Аутентификация
	serveMux.HandleFunc("/api/v1/auth", authHandler(userStorage, config))
	// Установка продолжительности работы выражений
	serveMux.HandleFunc("/duration", authMiddleware(durationHandler(userStorage)))
	// Получение выражения
	serveMux.HandleFunc("/expression", authMiddleware(expressionHandler(config, queue, storage, userStorage)))
	// Отдаем все сохраненные выражения
	serveMux.HandleFunc("/getexpressions", authMiddleware(getExpressionsHandler(storage)))
	// Получение выражения по ID
	serveMux.HandleFunc("/id/", authMiddleware(getById(storage)))
	// Регистрируем обработчики для воркеров
	serveMux.HandleFunc("/workers", getWorkers(config, queue))
	// Регистрируем обработчики WebSocket для выражений
	serveMux.HandleFunc("/ws", serveWS(config))
	return serveMux, nil
}

func authHandler(userStorage *memory.UserStorage, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Проблема с чтением данных:", "ОШИБКА:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := user.User(userStorage, data, conf)
		if err != nil {
			slog.Error("Невозможно добавить пользователя:", "ОШИБКА:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(token))
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.Split(auth, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			slog.Error("Неверные данные для аутентификации", "Токен:", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Выполняем все middleware на все запросы
func Decorate(next http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	decorated := next
	for i := len(middleware) - 1; i >= 0; i-- {
		decorated = middleware[i](decorated)
	}

	return decorated
}

// Возвращаем данные по ID
func getById(storage *memory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		patch := r.URL.Path
		auth := r.Header.Get("Authorization")
		token := strings.Split(auth, " ")
		data, err := expression.GetById(storage, patch[len("/id/"):], token[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write(data)
	}
}

// Возвращаем данные по воркерам
func getWorkers(conf *config.Config, q *queue.MapQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			newWorkers := config.Workers{
				Agents:      conf.AgentsAll.Load(),
				Workers:     conf.WorkersAll.Load(),
				WorkersBusy: conf.WorkersBusy.Load(),
			}
			array := q.GetQueue()
			newWorkers.Expressions = array
			data, err := json.Marshal(newWorkers)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(data)
		}
	}
}

// Читаем входящее выражение, валидируем его,сохраняем в память и возвращаем результат
func expressionHandler(config *config.Config,
	queue *queue.MapQueue,
	storage *memory.Storage,
	userStorage *memory.UserStorage,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Проблема с чтением данных:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slog.Info("Полученное выражение от пользователя:", "выражение:", string(data))
			auth := r.Header.Get("Authorization")
			token := strings.Split(auth, " ")
			answer, err := newexpression.NewExpression(
				config,
				storage,
				queue,
				string(data),
				token[1],
				userStorage,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(answer)
			slog.Info("Выражение добавлено в базу", "ответ:", string(answer))
		}
	}
}

func getExpressionsHandler(storage *memory.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			auth := r.Header.Get("Authorization")
			token := strings.Split(auth, " ")
			data, err := json.Marshal(expression.GetAllExpressions(storage, token[1]))
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ОШИБКА:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(data)
		}
	}
}

// Обрабатываем входящее время и возвращаем
func durationHandler(userStorage *memory.UserStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем входящее время
		auth := r.Header.Get("Authorization")
		token := strings.Split(auth, " ")
		if r.Method == http.MethodPost {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Проблема с чтением данных:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slog.Info("Полученное время для операций от пользователя:", "данные:", string(data))
			data, err = duration.ChangeDuration(data, token[1], userStorage)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			slog.Info("Время для операций обновлено и отправлено", "новое время:", string(data))
			w.Write(data)
		}
		// Возвращаем текущее установки времени
		if r.Method == http.MethodGet {
			data, err := duration.GetDuration(token[1], userStorage)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(data)
		}
	}
}
