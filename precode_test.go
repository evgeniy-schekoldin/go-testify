package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenRequestCorrect(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow&count=3", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	res := responseRecorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEmpty(t, res.Body)
}

func TestMainHandlerCityNotSupported(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=unknown&count=3", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	res := responseRecorder.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.Nil(t, err)
	require.Equal(t, res.StatusCode, http.StatusBadRequest)
	assert.Equal(t, string(body), "wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow&count=5", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)
	res := responseRecorder.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	require.Nil(t, err)
	require.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, len(strings.Split(string(body), ",")), len(cafeList["moscow"]))
}
