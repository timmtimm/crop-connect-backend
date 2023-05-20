package util

import "crop_connect/dto"

func FillNotAvailableMonth(object []dto.StatisticByYear) []dto.StatisticByYear {
	for i := 0; i < 12; i++ {
		isMonthAvailable := false
		for _, month := range object {
			if month.Month == i+1 {
				isMonthAvailable = true
				break
			}
		}

		if !isMonthAvailable {
			object = append(object, dto.StatisticByYear{
				Month: i + 1,
				Total: 0,
			})
		}
	}

	for i := 0; i < len(object); i++ {
		for j := i + 1; j < len(object); j++ {
			if object[i].Month > object[j].Month {
				temp := object[i]
				object[i] = object[j]
				object[j] = temp
			}
		}
	}

	return object
}
