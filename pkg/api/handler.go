package api

import (
	"net/http"
	"sprint13-14/pkg/nextdate"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(nextdate.DateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	res, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(res))
}
