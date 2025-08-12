package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат количества шагов")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, errors.New("неверный формат количества шагов")
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть положительным")
	}

	durationStr := strings.TrimSpace(parts[1])
	if !strings.Contains(durationStr, "s") && !strings.Contains(durationStr, "m") && !strings.Contains(durationStr, "h") {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration, err = time.ParseDuration(durationStr + "0s")
		if err != nil {
			return 0, 0, errors.New("неверный формат продолжительности")
		}
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) (string, error) {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return "", fmt.Errorf("ошибка разбора данных: %w", err)
	}

	sc := NewSpentCalories()
	distanceKm := sc.Distance(steps, height)
	
	var calories float64
	if duration.Minutes() < 30 {
		calories, err = sc.RunningCalories(steps, weight, height, duration)
	} else {
		calories, err = sc.WalkingCalories(steps, weight, height, duration)
	}
	
	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %w", err)
	}

	return formatDayActionInfo(steps, distanceKm, calories), nil
}

func formatDayActionInfo(steps int, distance, calories float64) string {
	return strings.Join([]string{
		"Количество шагов: " + strconv.Itoa(steps) + ".",
		"Дистанция составила " + strconv.FormatFloat(distance, 'f', 2, 64) + " км.",
		"Вы сожгли " + strconv.FormatFloat(calories, 'f', 2, 64) + " ккал.",
	}, "\n")
}
