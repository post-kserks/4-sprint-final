package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength float64 = 0.65
	// Количество метров в одном километре
	mInKm float64 = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		log.Printf("неверный формат данных: %s", data)
		return 0, 0, errors.New("неверный формат данных")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("неверный формат количества шагов: %s", parts[0])
		return 0, 0, errors.New("неверный формат количества шагов")
	}

	// Проверить, что количество шагов больше 0
	if steps <= 0 {
		log.Printf("количество шагов должно быть больше 0: %d", steps)
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// Преобразовать второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		log.Printf("неверный формат продолжительности: %s", parts[1])
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	// Проверить, что продолжительность больше 0
	if duration <= 0 {
		log.Printf("продолжительность должна быть больше 0: %s", duration)
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получить данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("ошибка при парсинге данных: %v", err)
		return ""
	}

	// Проверить, чтобы количество шагов было больше 0
	if steps <= 0 {
		log.Printf("количество шагов должно быть больше 0: %d", steps)
		return ""
	}

	// Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// Вычислить количество калорий, потраченных на прогулке
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("ошибка при расчете калорий: %v", err)
		return ""
	}

	// Сформировать строку с информацией
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
