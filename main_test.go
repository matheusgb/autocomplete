package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestPopulateIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/populate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(populateIndex)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, status)
	}

	expected := "data populated"
	if rr.Body.String() != expected {
		t.Errorf("expected %v, got %v", expected, rr.Body.String())
	}
}

func TestSendWord(t *testing.T) {
	body := []byte(`{"word": "TestWord"}`)
	req, err := http.NewRequest("POST", "/send", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sendWord)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, status)
	}

	expected := "data populated"
	if rr.Body.String() != expected {
		t.Errorf("expected %v, got %v", expected, rr.Body.String())
	}
}

func TestSendWordInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/send", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(sendWord)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("expected status %v, got %v", http.StatusMethodNotAllowed, status)
	}
}

func TestHandleWebSocket(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("failed to connect to websocket: %v", err)
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("TestWord"))
	if err != nil {
		t.Fatalf("failed to write message: %v", err)
	}

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	var response struct {
		Suggestions    []string `json:"suggestions"`
		FrequentValues []struct {
			Key   string `json:"key"`
			Count int    `json:"doc_count"`
		} `json:"frequent_values"`
	}

	err = json.Unmarshal(message, &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	ws.Close()
}
