package exchange

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base":"USD","target":"EUR","rate":0.85}`))
	}))
	defer server.Close() 

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")
	assert.NoError(t, err)
	assert.Equal(t, 0.85, rate)
}

func TestGetRate_BusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid currency pair"}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid currency pair")
}

func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json: true`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
}

func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Millisecond) 
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	service.Client.Timeout = 1 * time.Millisecond 
	
	_, err := service.GetRate("USD", "EUR")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
}

func TestGetRate_ServerPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status: 500")
}

func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	_, err := service.GetRate("USD", "EUR")
	assert.Error(t, err)
}