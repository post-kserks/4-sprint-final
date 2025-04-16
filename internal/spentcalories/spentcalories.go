package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    float64 = 0.65 // средняя длина шага.
	mInKm                      float64 = 1000 // количество метров в километре.
	minInH                     float64 = 60   // количество минут в часе.
	stepLengthCoefficient      float64 = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient float64 = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		log.Printf("неверный формат данных: %s", data)
		return 0, "", 0, errors.New("неверный формат данных")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("неверный формат количества шагов: %s", parts[0])
		return 0, "", 0, errors.New("неверный формат количества шагов")
	}

	// Проверить, что количество шагов больше 0
	if steps <= 0 {
		log.Printf("количество шагов должно быть больше 0: %d", steps)
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	// Получить вид активности
	activityType := parts[1]

	// Преобразовать третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		log.Printf("неверный формат продолжительности: %s", parts[2])
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	// Проверить, что продолжительность больше 0
	if duration <= 0 {
		log.Printf("продолжительность должна быть больше 0: %s", duration)
		return 0, "", 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитать длину шага
	stepLength := height * stepLengthCoefficient

	// Умножить пройденное количество шагов на длину шага
	distanceMeters := float64(steps) * stepLength

	// Разделить полученное значение на число метров в километре
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверить, что продолжительность duration больше 0
	if duration <= 0 {
		return 0
	}

	// Вычислить дистанцию с помощью distance()
	dist := distance(steps, height)

	// Вычислить и вернуть среднюю скорость
	// Для этого разделите дистанцию на продолжительность в часах
	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}

	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Рассчитать среднюю скорость с помощью meanSpeed()
	speed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитать и вернуть количество калорий
	// (weight * meanSpeed * durationInMinutes) / minInH
	return (weight * speed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("вес должен быть больше 0")
	}

	if height <= 0 {
		return 0, errors.New("рост должен быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше 0")
	}

	// Рассчитать среднюю скорость с помощью meanSpeed()
	speed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитать количество калорий
	calories := (weight * speed * durationInMinutes) / minInH

	// Умножить полученное число калорий на корректирующий коэффициент
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получить значения из строки данных с помощью функции parseTraining()
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	// Проверить, какой вид тренировки был передан в строке
	var calories float64
	var caloriesErr error

	switch activityType {
	case "Бег":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if caloriesErr != nil {
		return "", caloriesErr
	}

	// Рассчитать дистанцию и среднюю скорость
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Сформировать и вернуть строку с информацией
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, duration.Hours(), dist, speed, calories), nil
}
