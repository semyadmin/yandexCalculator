package sendtocalculate

import (
	"log/slog"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	grpcserver "github.com/adminsemy/yandexCalculator/Orchestrator/internal/grpc_server"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/tasks/queue"
)

type SendToCalculate struct {
	queue *queue.MapQueue
}

func NewSendToCalculate(queue *queue.MapQueue) *SendToCalculate {
	return &SendToCalculate{
		queue: queue,
	}
}

func (s *SendToCalculate) Dequeue() (grpcserver.Expression, error) {
	var bool bool
	var exp entity.Operation
	for !bool {
		exp, bool = s.queue.Dequeue()
	}
	slog.Info("Получена операция для отправки данных", "операция:", exp)
	return exp, nil
}

func (s *SendToCalculate) Done(id string, result float64, err string) {
	slog.Info("Операция выполнена", "операция:", id, "результат:", result, "ошибка:", err)
	s.queue.Done(id, result, err)
}
