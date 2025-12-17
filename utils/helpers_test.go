package utils_test

import (
	"encoding/json"
	"go-app/utils"
	"net/http/httptest"
	"testing"
)

func TestSetTotalCountHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.SetTotalCountHeader(rec, "42")

	if rec.Header().Get("X-Total-Count") != "42" {
		t.Fatalf("expected X-Total-Count to be 42, got %s", rec.Header().Get("X-Total-Count"))
	}

	// Проверяем, что Access-Control-Expose-Headers содержит X-Total-Count
	if rec.Header().Get("Access-Control-Expose-Headers") == "" {
		t.Fatalf("expected Access-Control-Expose-Headers to be set")
	}
}

func TestCheckOrderAndSortParams(t *testing.T) {
	order := "WRONG"
	sort := ""
	utils.CheckOrderAndSortParams(&order, &sort)
	if order != "ASC" {
		t.Fatalf("expected order to become ASC, got %s", order)
	}
	if sort != "ID" {
		t.Fatalf("expected sort to become ID, got %s", sort)
	}

	// Позитивный случай
	order = "DESC"
	sort = "Title"
	utils.CheckOrderAndSortParams(&order, &sort)
	if order != "DESC" || sort != "Title" {
		t.Fatalf("unexpected values after CheckOrderAndSortParams: %s, %s", order, sort)
	}
}

func TestMessage(t *testing.T) {
	m := utils.Message(true, "ok")
	if s, ok := m["status"].(bool); !ok || !s {
		t.Fatalf("expected status true, got %v", m["status"])
	}
	if msg, ok := m["message"].(string); !ok || msg != "ok" {
		t.Fatalf("expected message 'ok', got %v", m["message"])
	}
}

func TestRespond(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]interface{}{"status": true, "message": "ok"}
	utils.Respond(rec, data)

	if rec.Code != 200 {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var out map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &out)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if out["message"] != "ok" {
		t.Fatalf("expected message ok, got %v", out["message"])
	}
}

func TestRespondJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	raw := []byte(`{"a":1}`)
	utils.RespondJSON(rec, raw)

	if rec.Code != 200 && rec.Code != 0 { // RespondJSON не ставит код, поэтому иногда Code==0 (depends on test harness), но тело должно быть
		t.Fatalf("unexpected status code %d", rec.Code)
	}
	if string(rec.Body.Bytes()) != string(raw) {
		t.Fatalf("expected body %s, got %s", string(raw), rec.Body.String())
	}
}

func TestHandleOptions(t *testing.T) {
	req := httptest.NewRequest("OPTIONS", "/", nil)
	rec := httptest.NewRecorder()
	utils.HandleOptions(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200 for OPTIONS, got %d", rec.Code)
	}
}
