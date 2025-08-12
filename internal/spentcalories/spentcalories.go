package spentcalories

import (
	"errors"
	"fmt"
	"strings"
	"time"
)
const (
	lenStep                    = 0.65 
	mInKm                      = 1000 
	minInH                     = 60 
	stepLengthCoefficient      = 0.45 
	walkingCaloriesCoefficient = 0.5 
	runningCaloriesCoefficient = 0.029 
)

type SpentCalories struct{}

func NewSpentCalories() *SpentCalories {
	return &SpentCalories{}
}

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, errors.New("неверный формат количества шагов")
	}

	activity := strings.TrimSpace(parts[1])
	if activity != "Бег" && activity != "Ходьба" {
		return 0, "", 0, errors.New("неизвестный тип тренировки")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		duration, err = time.ParseDuration(durationStr + "0s")
		if err != nil {
			return 0, "", 0, errors.New("неверный формат продолжительности")
		}
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var dist, speed, calories float64
	var errCal error

	switch activity {
	case "Бег":
		calories, errCal = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, errCal = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if errCal != nil {
		return "", errCal
	}

	dist = distance(steps, height)
	speed = meanSpeed(steps, height, duration)

	return formatTrainingInfo(activity, duration, dist, speed, calories), nil
}

func formatTrainingInfo(activity string, duration time.Duration, distance, speed, calories float64) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	
	return strings.Join([]string{
		"Тип тренировки: " + activity,
		fmt.Sprintf("Длительность: %d ч. %d мин.", hours, minutes),
		"Дистанция: " + fmt.Sprintf("%.2f", distance) + " км.",
		"Скорость: " + fmt.Sprintf("%.2f", speed) + " км/ч",
		"Сожгли калорий: " + fmt.Sprintf("%.2f", calories),
	}, "\n")
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return runningCaloriesCoefficient * weight * speed * durationInMinutes, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	distanceKm := distance(steps, height)
	durationInHours := duration.Hours()
	return walkingCaloriesCoefficient * weight * distanceKm / durationInHours, nil
}
