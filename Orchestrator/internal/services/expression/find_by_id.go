package expression

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/entity"
	"github.com/adminsemy/yandexCalculator/Orchestrator/internal/storage/memory"
)

func GetById(storage *memory.Storage, number string) ([]byte, error) {
	id, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		slog.Error("Невозможно распарсить ID:", "ОШИБКА:", err)
		return nil, err
	}
	exp, err := storage.GeById(id)
	if err != nil {
		slog.Error("Невозможно получить данные по ID:", "ОШИБКА:", err, "ID:", id)
		return nil, err
	}
	resp := entity.NewResponseExpression(exp.ID, exp.Expression, exp.Start, exp.Duration, exp.IsCalc, exp.Result, exp.Err)
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Невозможно сериализовать данные:", "ОШИБКА:", err)
		return nil, err
	}
	return data, nil
}
