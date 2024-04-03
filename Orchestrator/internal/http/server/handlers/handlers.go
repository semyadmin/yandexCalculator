package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_ast"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/postgresql/postgresql_config"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/arithmetic"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/validator"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/web_socket/client"
)

func NewServeMux(config *config.Config,
	queue *queue.MapQueue,
	storage *memory.Storage,
) (http.Handler, error) {
	// Создам маршрутизатор
	serveMux := http.NewServeMux()
	// Регистрируем обработчики событий
	patchToFront := "./frontend/build"
	serveMux.Handle("/", http.FileServer(http.Dir(patchToFront)))
	serveMux.HandleFunc("/duration", durationHandler(config))
	serveMux.HandleFunc("/expression", expressionHandler(config, queue, storage))
	serveMux.HandleFunc("/id/", getById(storage))
	serveMux.HandleFunc("/workers", getWorkers(config, queue))
	serveMux.HandleFunc("/ws", serveWS(config))
	return serveMux, nil
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
func getById(storage *memory.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		patch := r.URL.Path
		id, err := strconv.ParseUint(patch[len("/id/"):], 10, 64)
		if err != nil {
			slog.Error("Невозможно распарсить ID:", "ОШИБКА:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		exp, err := storage.GeById(id)
		if err != nil {
			slog.Error("Невозможно получить данные по ID:", "ОШИБКА:", err, "ID:", id)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp := arithmetic.NewExpression(exp.Expression)
		data, err := json.Marshal(resp)
		if err != nil {
			slog.Error("Невозможно сериализовать данные:", "ОШИБКА:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

// Возвращаем данные по воркерам
func getWorkers(conf *config.Config, q *queue.MapQueue) func(w http.ResponseWriter, r *http.Request) {
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
) func(w http.ResponseWriter, r *http.Request) {
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
			// Формируем новое выражение для вычисления
			exp, err := arithmetic.NewASTTree(string(data), config, queue, validator.Validator)
			if err != nil {
				slog.Error("Проблема с вычислением выражения:", "выражение:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Проверяем, есть ли такое выражение в базе. Если есть - отдаем
			dataInfo, err := storage.GeByExpression(exp.Expression)
			if err == nil {
				resp := arithmetic.NewExpression(dataInfo.Expression)
				data, err := json.Marshal(resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				w.WriteHeader(http.StatusAccepted)
				w.Write(data)
				slog.Info("Такое выражение уже было в базе", "ответ:", string(data))
				return
			}
			// Сохраняем в память
			storage.Set(exp, "new")
			postgresql_ast.Add(exp, config)
			resp := arithmetic.NewExpression(exp)
			answer, err := json.Marshal(resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			w.Write(answer)
			go func() {
				config.WSmanager.MessageCh <- &client.Message{
					Payload: answer,
					Type:    client.ClientExpression,
				}
			}()
			slog.Info("Выражение добавлено в базу", "ответ:", string(answer))
		}
	}
}

// Обрабатываем входящее время и возвращаем
func durationHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обрабатываем входящее время
		if r.Method == http.MethodPost {
			newDuration := config.ConfigExpression{}
			data, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Проблема с чтением данных:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			slog.Info("Полученное время для операций от пользователя:", "данные:", string(data))
			json.Unmarshal(data, &newDuration)
			err = conf.NewDuration(&newDuration)
			if err != nil {
				slog.Error("Некорректное время:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			data, err = json.Marshal(newDuration)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Сохраняем в базу
			postgresql_config.Save(conf)
			slog.Info("Время для операций обновлено и отправлено", "новое время:", newDuration)
			w.Write(data)
		}
		// Возвращаем текущее установки времени
		if r.Method == http.MethodGet {
			newDuration := config.ConfigExpression{}
			newDuration.Init(conf)
			data, err := json.Marshal(newDuration)
			if err != nil {
				slog.Error("Невозможно сериализовать данные:", "ошибка:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(data)
		}
	}
}
