package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %v", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("days interval not specified")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days > 400 || days < 1 {
			return "", errors.New("invalid days interval")
		}

		next := start
		for {
			next = next.AddDate(0, 0, days)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateLayout), nil

	case "y":
		next := start
		for {
			next = next.AddDate(1, 0, 0)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateLayout), nil

	case "w":
		if len(parts) < 2 {
			return "", errors.New("week days not specified")
		}
		weekDays := strings.Split(parts[1], ",")

		next := start
		for {
			next = next.AddDate(0, 0, 1) // Шагаем по дню

			weekday := int(next.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			for _, d := range weekDays {
				targetDay, _ := strconv.Atoi(d)
				if weekday == targetDay && next.After(now) {
					return next.Format(DateLayout), nil
				}
			}
			if next.After(now.AddDate(2, 0, 0)) {
				break
			}
		}

	case "m":
		if len(parts) < 2 {
			return "", errors.New("month days not specified")
		}

		daysParam := strings.Split(parts[1], ",")
		var monthsParam []string
		if len(parts) > 2 {
			monthsParam = strings.Split(parts[2], ",")
		}

		next := start
		for {
			next = next.AddDate(0, 0, 1)

			// 1. Проверяем месяц (если он указан)
			if len(monthsParam) > 0 {
				matchMonth := false
				currentMonth := int(next.Month())
				for _, m := range monthsParam {
					targetMonth, _ := strconv.Atoi(m)
					if currentMonth == targetMonth {
						matchMonth = true
						break
					}
				}
				if !matchMonth {
					continue
				}
			}

			// 2. Проверяем день месяца
			currentDay := next.Day()
			lastDay := time.Date(next.Year(), next.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			matchDay := false
			for _, dStr := range daysParam {
				d, _ := strconv.Atoi(dStr)

				targetDay := d
				if d < 0 {
					// Логика для -1 (последний), -2 (предпоследний)
					targetDay = lastDay + d + 1
				}

				if currentDay == targetDay && next.After(now) {
					matchDay = true
					break
				}
			}

			if matchDay {
				return next.Format(DateLayout), nil
			}

			if next.After(now.AddDate(5, 0, 0)) {
				break
			}
		}
	}
	return "", errors.New("could not calculate next date")
}
