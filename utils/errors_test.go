package utils_test

import (
	"encoding/json"
	"errors"
	"go-app/utils"
	"net/http/httptest"
	"testing"
)

func decodeMessage(rec *httptest.ResponseRecorder, t *testing.T) map[string]interface{} {
	var out map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	return out
}

func TestHandleBadRequest(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.HandleBadRequest(rec, errors.New("bad request"))

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	out := decodeMessage(rec, t)
	if out["message"] != "bad request" {
		t.Fatalf("expected 'bad request', got %v", out["message"])
	}
	if status, ok := out["status"].(bool); !ok || status {
		t.Fatalf("expected status false, got %v", out["status"])
	}
}

func TestHandleNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.HandleNotFound(rec)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	out := decodeMessage(rec, t)
	if out["message"] != "not found" {
		t.Fatalf("expected 'not found', got %v", out["message"])
	}
}

func TestHandleInternalError(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.HandleInternalError(rec, errors.New("boom"))

	if rec.Code != 500 {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	out := decodeMessage(rec, t)
	if out["message"] != "boom" {
		t.Fatalf("expected 'boom', got %v", out["message"])
	}
}

func TestHandleUnauthorized(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.HandleUnauthorized(rec, errors.New("no auth"))

	if rec.Code != 401 {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	out := decodeMessage(rec, t)
	if out["message"] != "no auth" {
		t.Fatalf("expected 'no auth', got %v", out["message"])
	}
}

func TestHandleForbidden(t *testing.T) {
	rec := httptest.NewRecorder()
	utils.HandleForbidden(rec, errors.New("forbidden"))

	if rec.Code != 403 {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
	out := decodeMessage(rec, t)
	if out["message"] != "forbidden" {
		t.Fatalf("expected 'forbidden', got %v", out["message"])
	}
}
