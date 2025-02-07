package main

import (
	"fmt"
	"math"
	"time"
)

const (
	// Общие константы для вычислений.
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре

	// Константы для расчета потраченных килокалорий при беге.
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости

	// Константы для расчета потраченных килокалорий при ходьбе.
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с

	// Константы для расчета потраченных килокалорий при плавании.
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса пользователя
)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	if t.Duration == 0 {
		return 0
	}
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0
}

type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo возвращает структуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

type Running struct {
	Training
}

// Calories возвращает количество потраченных килокалорий при беге.
func (r Running) Calories() float64 {
	// Вычисление калорий для бега с учетом средней скорости и веса.
	return (CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours
}

func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

type Walking struct {
	Training
	Height float64 // рост пользователя в метрах
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
func (w Walking) Calories() float64 {
	// Преобразование скорости в м/с
	speedMInS := w.meanSpeed() * KmHInMsec
	// Преобразование роста в метры
	heightInM := w.Height / CmInM
	// Вычисление калорий для ходьбы с учетом веса, скорости и роста.
	return (CaloriesWeightMultiplier*w.Weight + (math.Pow(speedMInS, 2)/heightInM)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours
}

func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

// meanSpeed возвращает среднюю скорость при плавании.
func (s Swimming) meanSpeed() float64 {
	if s.Duration == 0 {
		return 0
	}
	return float64(s.LengthPool) * float64(s.CountPool) / MInKm / float64(s.Duration.Hours())
}

// Calories возвращает количество калорий, потраченных при плавании.
func (s Swimming) Calories() float64 {
	return (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     s.distance(),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	// Вычисление количества калорий для тренировки
	calories := training.Calories()
	// Получение информации о тренировке
	info := training.TrainingInfo()
	// Добавление количества калорий в информацию о тренировке
	info.Calories = calories

	return fmt.Sprint(info)
}

func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))
}
