package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type DaySteps struct{}

func (ds DaySteps) ParsePackage(data string) (int, time.Duration, error) {
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
	


	duration, err := time.ParseDuration(durationStr)
	durationStr := strings.TrimSpace(parts[1])
	if !strings.Contains(durationStr, "s") && !strings.Contains(durationStr, "m") && !strings.Contains(durationStr, "h") {
		return 0, 0, errors.New("неверный формат продолжительности")
	}
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		
		duration, err = time.ParseDuration(durationStr + "0s")
		if err != nil {
			return 0, 0, errors.New("неверный формат продолжительности")
		}
	}

	return steps, duration, nil
}

func (ds DaySteps) ParseTraining(data string) (int, string, time.Duration, error) {
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

type SpentCalories struct {
	StepLength        float64
	MInKm             float64
	StepLengthCoeff   float64
	WalkingCoefficient float64
	RunningCoefficient float64
}

func NewSpentCalories() *SpentCalories {
	return &SpentCalories{
		StepLength:        0.65,
		MInKm:             1000,
		StepLengthCoeff:   0.414,
		WalkingCoefficient: 0.035,
		RunningCoefficient: 0.029,
	}
}

func (sc *SpentCalories) Distance(steps int, height float64) float64 {
	return float64(steps) * height * sc.StepLengthCoeff / sc.MInKm
}

func (sc *SpentCalories) MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := sc.Distance(steps, height)
	return dist / duration.Hours()
}

func (sc *SpentCalories) RunningCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	speed := sc.MeanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return sc.RunningCoefficient * weight * speed * durationInMinutes, nil
}

func (sc *SpentCalories) WalkingCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры")
	}
	distanceKm := sc.Distance(steps, height)
	durationInHours := duration.Hours()
	return sc.WalkingCoefficient * weight * distanceKm / durationInHours, nil
}

var (
	ds = DaySteps{}
	sc = NewSpentCalories()
)

func DayActionInfo(data string, weight, height float64) (string, error) {
	steps, duration, err := ds.ParsePackage(data)
	if err != nil {
		return "", fmt.Errorf("ошибка разбора данных: %w", err)
	}

	distanceKm := sc.Distance(steps, height)
	
		var calories float64
	if duration.Minutes() < 30 { // Если меньше 30 минут - считаем бег
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

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := ds.ParseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var dist, speed, calories float64
	var errCal error

	switch activity {
	case "Бег":
		calories, errCal = sc.RunningCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, errCal = sc.WalkingCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if errCal != nil {
		return "", errCal
	}

	dist = sc.Distance(steps, height)
	speed = sc.MeanSpeed(steps, height, duration)

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

func main() {
	weight := 84.6
	height := 1.87

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
		dayActionsInfo, err := DayActionInfo(v, weight, height)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			continue
		}
		fmt.Println(dayActionsInfo)
		fmt.Println()
	}

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


