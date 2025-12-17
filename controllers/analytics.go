package controllers

import (
	"fmt"
	"go-app/db"
	"go-app/models"
	u "go-app/utils"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	totalRequests  uint64
	totalLatencyNs uint64
	appStart       = time.Now()
)

var MetricsMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		lat := time.Since(start)
		atomic.AddUint64(&totalRequests, 1)
		atomic.AddUint64(&totalLatencyNs, uint64(lat.Nanoseconds()))
	})
}

func AnalyticsSummary(w http.ResponseWriter, r *http.Request) {
	DB := db.GetDB()
	if DB == nil {
		u.HandleInternalError(w, fmt.Errorf("db not initialized"))
		return
	}

	var totalNotes int64
	if err := DB.Model(&models.Note{}).Count(&totalNotes).Error; err != nil {
		u.HandleBadRequest(w, err)
		return
	}

	// получаем avg length
	type resRow struct {
		Avg *float64
	}
	row := DB.Raw("SELECT AVG(LENGTH(content)) as avg FROM notes").Row()
	var raw interface{}
	if err := row.Scan(&raw); err != nil {
		// если ошибка — вернуть 0
		raw = nil
	}
	var avgLen float64
	if raw == nil {
		avgLen = 0
	} else {
		switch v := raw.(type) {
		case float64:
			avgLen = v
		case int64:
			avgLen = float64(v)
		case []uint8:
			// sqlite возвращает []byte
			var parsed float64
			_, _ = fmt.Sscanf(string(v), "%f", &parsed)
			avgLen = parsed
		default:
			avgLen = 0
		}
	}

	tr := atomic.LoadUint64(&totalRequests)
	var avgLatencyMs float64
	if tr > 0 {
		avgLatencyMs = float64(atomic.LoadUint64(&totalLatencyNs)) / float64(tr) / 1e6
	} else {
		avgLatencyMs = 0
	}

	resp := map[string]interface{}{
		"total_notes":            totalNotes,
		"avg_note_length":        avgLen,
		"total_requests":         tr,
		"avg_request_latency_ms": avgLatencyMs,
		"uptime_sec":             int(time.Since(appStart).Seconds()),
	}
	u.Respond(w, resp)
}

func AnalyticsNotesCount(w http.ResponseWriter, r *http.Request) {
	DB := db.GetDB()
	if DB == nil {
		u.HandleInternalError(w, fmt.Errorf("db not initialized"))
		return
	}
	var count int64
	if err := DB.Model(&models.Note{}).Count(&count).Error; err != nil {
		u.HandleBadRequest(w, err)
		return
	}
	u.Respond(w, map[string]interface{}{"total_notes": count})
}

func AnalyticsAvgNoteLength(w http.ResponseWriter, r *http.Request) {
	DB := db.GetDB()
	if DB == nil {
		u.HandleInternalError(w, fmt.Errorf("db not initialized"))
		return
	}
	row := DB.Raw("SELECT AVG(LENGTH(content)) FROM notes").Row()
	var raw interface{}
	if err := row.Scan(&raw); err != nil || raw == nil {
		u.Respond(w, map[string]interface{}{"avg_note_length": 0})
		return
	}
	var avg float64
	switch v := raw.(type) {
	case float64:
		avg = v
	case int64:
		avg = float64(v)
	case []uint8:
		var parsed float64
		_, _ = fmt.Sscanf(string(v), "%f", &parsed)
		avg = parsed
	default:
		avg = 0
	}
	u.Respond(w, map[string]interface{}{"avg_note_length": avg})
}
