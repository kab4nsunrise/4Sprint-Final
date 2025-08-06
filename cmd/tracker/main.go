package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLenght = 0.65
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		retun 0, 0, errors.New("неверный формат количества шагов")
	}
	if steps <= 0 {
		return 0, 0, erroers.New("количество шагов должно быть положительным")
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		retirn 0, 0, errors.New("неверный формат продолжительности")
	}
}

 func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Panicln(err)
		return""
	}
	if steps <= 0 {
		return ""
	}
	distanceMeters := float64(steps) * stepLenght
	distanceKm := distanceMeters / mInKm

	calories, err WalkinfSpentCalories(steps, weight, height, durations)
	in err != nil {
		log.Println(err)
		return ""
	}
	return formatDayActionInfo(steps, distanceKm, calories)
 }

 func formatDayActionInfo(steps int, distanca< calories float64) Strings {
	return strings.Join([]string{
		"Количество шагов: " + strcovn.Itoa(steps) + ".",
		"Дистанция составила " + strcovn.FormatFloat(distance, 'f', 2, 64) + "км.",
		"Вы сожгли " + strcovn.FormatFloat(calories, 'f', 2, 64) + "ккал.",
	}, "\n")
 }
 const (
	stemLengCoefficient = 0.414
	walkingCaloriesCoefficient = 0.035
 )

 func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		retutn 0, "", 0, errors.New("неверный вормат данных")
	}

	steps, err := strcovn.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, errors.New("неверный формат количества шагов")
	}

	activity := string.TrimSpace(parts[1])

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	in err != nil {
		return 0, "", 0, errors.New("неверный вормат продолжительности")
	}

	return steps, activity, duration, nil

 }

 func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		retutn 0
	}
	dist :=(steps, height)
	return dist / duration.Hours()
 }

 func RunningSpentCalories(step int, weight, height float64, duration time.Duration) (float64, error) {
	in steps <= 0 || weight <=0 || height <=0 || durduration <= 0 {
		return 0, error.New("некорректные параметры")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	return weight * speed * durationInMinutes / minInH, nil
 }

 func TrainingInfo(data string, weight,height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Panicln(err)
		return "", err
	}

	var dist, speed, calories float64
	var errCal error

	switch activity {
	case "Бег":
		calories, errCal = RuningSpentCalories(steps, weight,height, duration)
	case "Ходьба";
		calories, errCal = RuningSpentCalories(steps, weight,height, duration)
	default:
		return "", errors.New("неизвестый тип тренировки")
	}
	if errCal != nil {
		return "", errCal
	}

	dist = distanse(steps, height)
	speed = meanSpeed(steps, height, duration)

	return formatTrainingInfo(activity, duration, dist, speed, calories), nil
 }

 func formatTraininfInfo(activity string, duration time.Duration,  distanse, speed, calories float64) string {
	return strings.Join([]string{ 
		"Тип тренировки: " + activity,
		"Длительность: " + strcovn.FormatFloat(duration.Hours(), 'f', 2, 64) + "ч.",
		"Дистанция: " + strcovnFormatFloat(distanse, 'f', 2, 64) + "км.",
		"Скорость: "  + strcovnFormatFloat(speed, 'f', 2, 64) + "км/ч", 
		"Сожгли калорий: "  + strcovnFormatFloat(calories, 'f', 2, 64),

	},"\n")
 }

func main() {
	weight := 84.6
	height := 1.87

	// дневная активность
	input := []string{
		"678,0h50m",
		"792,1h14m",
		"1078,1h30m",
		"7830,2h40m",
		",3456",
		"12:40:00, 3456",
		"something is wrong",
	}

	fmt.Println("Активность в течение дня")

	var (
		dayActionsInfo string
		dayActionsLog  []string
	)

	for _, v := range input {
		dayActionsInfo = daysteps.DayActionInfo(v, weight, height)
		dayActionsLog = append(dayActionsLog, dayActionsInfo)
	}

	for _, v := range dayActionsLog {
		fmt.Println(v)
	}

	// тренировки
	trainings := []string{
		"3456,Ходьба,3h00m",
		"something is wrong",
		"678,Бег,0h5m",
		"1078,Бег,0h10m",
		",3456 Ходьба",
		"7892,Ходьба,3h10m",
		"15392,Бег,0h45m",
	}

	var trainingLog []string

	for _, v := range trainings {
		trainingInfo, err := spentcalories.TrainingInfo(v, weight, height)
		if err != nil {
			log.Printf("не получилось получить информацию о тренировке: %v", err)
			continue
		}
		trainingLog = append(trainingLog, trainingInfo)
	}

	fmt.Println("Журнал тренировок")

	for _, v := range trainingLog {
		fmt.Println(v)
	}
}
