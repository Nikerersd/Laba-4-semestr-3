package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	Name      string
	Position  string
	Age       int
	Salary    float64
}

func calculateAverageAge(workers []Worker, position string) float64 {
	var totalAge, count int
	for _, worker := range workers {
		if worker.Position == position {
			totalAge += worker.Age
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return float64(totalAge) / float64(count)
}

func findMaxSalary(workers []Worker, position string, avgAge float64) float64 {
	var maxSalary float64
	for _, worker := range workers {
		if worker.Position == position && abs(float64(worker.Age)-avgAge) <= 2 {
			if worker.Salary > maxSalary {
				maxSalary = worker.Salary
			}
		}
	}
	return maxSalary
}

// Абсолютная разница
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func processWithoutConcurrency(workers []Worker, position string) {
	start := time.Now()

	var avgAge float64
	var maxSalary float64

	countSize := 3
	subsetsSize := len(workers) / countSize

	avgAgeResults := make([]float64, countSize)
	maxSalaryResults := make([]float64, countSize)

	//Поиск среднего возраста
	for i := 0; i < countSize; i++ {
		func(i int) {
			startIndex := i * subsetsSize
			endIndex := (i + 1) * subsetsSize
			if i == countSize-1 { // Для последнего сегмента
				endIndex = len(workers)
			}
			avgAgeResults[i] = calculateAverageAge(workers[startIndex:endIndex], position)
		}(i)
	}


	// Объединяем результаты среднего возраста
	var totalAge float64
	var count int
	for _, avg := range avgAgeResults {
		if avg > 0 {
			totalAge += avg
			count++
		}
	}
	if count > 0 {
		avgAge = totalAge / float64(count)
	}
	
	//Поиск максимальной зарплаты
	for i := 0; i < countSize; i++ {
		func(i int) {
			startIndex := i * subsetsSize
			endIndex := (i + 1) * subsetsSize
			if i == countSize-1 { // Для последнего сегмента
				endIndex = len(workers)
			}
			maxSalaryResults[i] = findMaxSalary(workers[startIndex:endIndex], position, avgAge)
		}(i)
	}

	// Объединяем результаты максимальной зарплаты
	for _, max := range maxSalaryResults {
		if max > maxSalary {
			maxSalary = max
		}
	}

	duration := time.Since(start)

	fmt.Printf("Без многозадачности:\n")
	fmt.Printf("Средний возраст: %.2f\n", avgAge)
	fmt.Printf("Максимальная зарплата: %.2f\n", maxSalary)
	fmt.Printf("Время обработки: %v\n\n", duration)
}

func processWithConcurrency(workers []Worker, position string) {
	start := time.Now()

	var wg sync.WaitGroup
	var avgAge float64
	var maxSalary float64

	// Разбиваем данные на части для многозадачности
	numGoroutines := 3 // Количество горутин
	chunkSize := len(workers) / numGoroutines

	// Каналы для хранения результатов
	avgAgeResults := make([]float64, numGoroutines)
	maxSalaryResults := make([]float64, numGoroutines)

	wg.Add(3)
	// Горутин для вычисления среднего возраста
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			startIndex := i * chunkSize
			endIndex := (i + 1) * chunkSize
			if i == numGoroutines-1 { // Для последнего сегмента
				endIndex = len(workers)
			}
			avgAgeResults[i] = calculateAverageAge(workers[startIndex:endIndex], position)
		}(i)
	}

	// Ждем завершения всех горутин для среднего возраста
	wg.Wait()

	// Объединяем результаты среднего возраста
	var totalAge float64
	var count int
	for _, avg := range avgAgeResults {
		if avg > 0 {
			totalAge += avg
			count++
		}
	}
	if count > 0 {
		avgAge = totalAge / float64(count)
	}
	
	wg.Add(3)
	// Горутин для поиска максимальной зарплаты
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			startIndex := i * chunkSize
			endIndex := (i + 1) * chunkSize
			if i == numGoroutines-1 { // Для последнего сегмента
				endIndex = len(workers)
			}
			maxSalaryResults[i] = findMaxSalary(workers[startIndex:endIndex], position, avgAge)
		}(i)
	}

	// Ждем завершения всех горутин для максимальной зарплаты
	wg.Wait()

	// Объединяем результаты максимальной зарплаты
	for _, max := range maxSalaryResults {
		if max > maxSalary {
			maxSalary = max
		}
	}

	duration := time.Since(start)

	fmt.Printf("С многозадачностью (с несколькими горутинами):\n")
	fmt.Printf("Средний возраст: %.2f\n", avgAge)
	fmt.Printf("Максимальная зарплата: %.2f\n", maxSalary)
	fmt.Printf("Время обработки: %v\n\n", duration)
}

func main() {
	// Массив работников
	workers := []Worker{
		{"Иванов Иван", "Д", 30, 50000},
		{"Петров Петр", "Д", 32, 60000},
		{"Сидоров Сидор", "Д", 28, 55000},
		{"Кузнецова Ольга", "С", 40, 70000},
		{"Морозов Алексей", "С", 42, 75000},
	}

	// Позиция для обработки
	position := "Д"

	// Обработка без многозадачности
	processWithoutConcurrency(workers, position)

	// Обработка с многозадачностью
	processWithConcurrency(workers, position)
}