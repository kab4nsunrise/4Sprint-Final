package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLength       = 0.65
	mInKm            = 1000
	minInH           = 60
	stepLengthCoeff  = 0.414
	walkingCaloriesCoefficient = 0.035
	runningCaloriesCoefficient = 0.029
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

	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return formatDayActionInfo(steps, distanceKm, calories)
}

func formatDayActionInfo(steps int, distance, calories float64) string {
	return strings.Join([]string{
		"Количество шагов: " + strconv.Itoa(steps) + ".",
		"Дистанция составила " + strconv.FormatFloat(distance, 'f', 2, 64) + " км.",
		"Вы сожгли " + strconv.FormatFloat(calories, 'f', 2, 64) + " ккал.",
	}, "\n")
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

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoeff / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	return dist / duration.Hours()
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

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
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
	return strings.Join([]string{
		"Тип тренировки: " + activity,
		"Длительность: " + fmt.Sprintf("%.2f", duration.Hours()) + " ч.",
		"Дистанция: " + fmt.Sprintf("%.2f", distance) + " км.",
		"Скорость: " + fmt.Sprintf("%.2f", speed) + " км/ч",
		"Сожгли калорий: " + fmt.Sprintf("%.2f", calories),
	}, "\n")
}

func main() {
	weight := 84.6
	height := 1.87

	// дневная активность
	input := []string{
		"678,50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		"3456,1h",
		"1234,40m",
		"5678,2h30m",
	}

	fmt.Println("Активность в течение дня")

	for _, v := range input {
		dayActionsInfo := DayActionInfo(v, weight, height)
		if dayActionsInfo != "" {
			fmt.Println(dayActionsInfo)
			fmt.Println()
		}
	}

	// тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"678,Бег,5m",
		"1078,Бег,10m",
		"7892,Ходьба,3h10m",
		"15392,Бег,45m",
	}

	fmt.Println("\nЖурнал тренировок")

	for _, v := range trainings {
		trainingInfo, err := TrainingInfo(v, weight, height)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			continue
		}
		fmt.Println(trainingInfo)
		fmt.Println()
	}
}
