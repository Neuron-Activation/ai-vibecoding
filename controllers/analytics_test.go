package controllers

import (
	"go-app/db"
	"go-app/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	_ = db.InitDB("sqlite3", ":memory:")
	code := m.Run()
	_ = db.CloseDB()
	os.Exit(code)
}

func TestAnalyticsEndpoints(t *testing.T) {
	// Добавим пару заметок
	DB := db.GetDB()
	if DB == nil {
		t.Fatal("db is nil")
	}
	DB.Create(&models.Note{Title: "t1", Content: "hello"})
	DB.Create(&models.Note{Title: "t2", Content: "hello world"})

	// Протестируем /analytics/notes/count
	req := httptest.NewRequest("GET", "/analytics/notes/count", nil)
	rec := httptest.NewRecorder()
	AnalyticsNotesCount(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "\"total_notes\"") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}

	// Протестируем /analytics/summary (через middleware, чтобы считать request)
	req2 := httptest.NewRequest("GET", "/analytics/summary", nil)
	rec2 := httptest.NewRecorder()
	// обернём handler метриками
	handler := MetricsMiddleware(http.HandlerFunc(AnalyticsSummary))
	handler.ServeHTTP(rec2, req2)

	if rec2.Code != 200 {
		t.Fatalf("expected 200 for summary, got %d", rec2.Code)
	}
	body := rec2.Body.String()
	if !strings.Contains(body, "\"total_notes\"") || !strings.Contains(body, "\"avg_note_length\"") {
		t.Fatalf("unexpected summary body: %s", body)
	}
}
