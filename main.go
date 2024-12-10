package main

import (
	"fmt"
	"math"
	"sync"
)

func main() {
	// Данные температуры
	temperatures := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	// Канал для передачи результатов сгрупиированных темепратур
	resultChanel := make(chan map[int][]float64)

	// Мьютекс результататов
	var temperatures_mutex sync.Mutex

	// Группа ожидания для горутин
	var temperatures_group sync.WaitGroup

	// Функция для обработки каждого значения температуры
	go func() {
		// Создаем map для хранения результатов
		result := make(map[int][]float64)

		// Обрабатываем каждую температуру и определяем ее к группе
		for _, temp := range temperatures {
			group := int(math.Round(temp/10) * 10)
			temperatures_mutex.Lock()
			result[group] = append(result[group], temp)
			temperatures_mutex.Unlock()
		}

		// Отправляем результат в канал
		resultChanel <- result

		// Завершаем горутину
		temperatures_group.Done()
	}()

	temperatures_group.Add(1)

	// Закрываем канал по завершении горутины
	go func() {
		temperatures_group.Wait()
		close(resultChanel)
	}()

	// Выводим результат
	for result := range resultChanel {
		for key, values := range result {
			fmt.Printf("%d: {", key)
			for i, val := range values {
				if i == len(values)-1 {
					fmt.Printf("%.1f", val)
				} else {
					fmt.Printf("%.1f, ", val)
				}
			}
			fmt.Println("}")
		}
	}
}
