package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if dstart == "" {
		return "", errors.New("invalid start date")
	}

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

		var targets []int
		for _, d := range weekDays {
			val, err := strconv.Atoi(d)
			if err != nil || val < 1 || val > 7 {
				return "", errors.New("invalid weekday")
			}
			targets = append(targets, val)
		}

		next := start
		if now.Sub(next).Hours() > 24*7 {

			weeksBreak := int(now.Sub(next).Hours() / (24 * 7))

			next = next.AddDate(0, 0, weeksBreak*7)

		}
		for {
			next = next.AddDate(0, 0, 1)

			weekday := int(next.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			for _, targetDay := range targets {
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

		var targetDays []int
		for _, dStr := range daysParam {
			d, err := strconv.Atoi(dStr)
			if err != nil || d < -2 || d > 31 || d == 0 {
				return "", errors.New("invalid month day")
			}
			targetDays = append(targetDays, d)
		}

		var targetMonths []int
		if len(parts) > 2 {
			monthsParam := strings.Split(parts[2], ",")

			for _, mStr := range monthsParam {
				m, err := strconv.Atoi(mStr)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid month")
				}
				targetMonths = append(targetMonths, m)
			}
		}

		next := start
		if next.Before(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())) {

			next = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			next = next.AddDate(0, 0, -1)
		}

		for {
			next = next.AddDate(0, 0, 1)

			if len(targetMonths) > 0 {
				matchMonth := false
				currentMonth := int(next.Month())
				for _, targetMonth := range targetMonths {
					if currentMonth == targetMonth {
						matchMonth = true
						break
					}
				}
				if !matchMonth {
					continue
				}
			}

			currentDay := next.Day()
			lastDay := time.Date(next.Year(), next.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()

			matchDay := false
			for _, d := range targetDays {
				targetDay := d
				if d < 0 {
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
	default:
		return "", errors.New("unknown repeat rule")
	}
	return "", errors.New("could not calculate next date")
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	res, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
}
