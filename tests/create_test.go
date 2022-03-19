package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CristianCurteanu/url-shortener/pkg"
	"github.com/CristianCurteanu/url-shortener/pkg/mappings"
	"github.com/CristianCurteanu/url-shortener/pkg/mappings/handlers"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	expect "github.com/stretchr/testify/require"
)

type UrlMappingDAOMock struct {
	mock.Mock
}

func (m *UrlMappingDAOMock) Add(ctx context.Context, mapping mappings.UrlMapping) error {
	args := m.Called(ctx, mapping)
	return args.Error(0)
}

func (m *UrlMappingDAOMock) SearchById(ctx context.Context, key string) (mappings.UrlMapping, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(mappings.UrlMapping), args.Error(1)
}

func (m *UrlMappingDAOMock) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *UrlMappingDAOMock) IncrementCounter(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	app := pkg.NewApp(&pkg.Config{
		Host:     ":8080",
		MongoURL: "mongodb://localhost:4343",
		Redis: pkg.RedisConfig{
			URL: mr.Addr(),
		},
	})
	require.NoError(t, app.Init())

	dao := new(UrlMappingDAOMock)
	dao.On("Add", mock.Anything, mock.Anything).Return(nil)

	app.UrlsDAO = dao
	url := "https://google.com"

	body, err := json.Marshal(handlers.URLCreationRequest{url})
	expect.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/url", bytes.NewReader(body))

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusCreated, w.Code)
	expect.Contains(t, string(w.Body.Bytes()), mappings.CreateKey(url))
}

func TestCreate_FailIfNotURL(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	app := pkg.NewApp(&pkg.Config{
		Host:     ":8080",
		MongoURL: "mongodb://localhost:4343",
		Redis: pkg.RedisConfig{
			URL: mr.Addr(),
		},
	})
	require.NoError(t, app.Init())

	app.UrlsDAO = new(UrlMappingDAOMock)
	url := "randomstring"

	body, err := json.Marshal(handlers.URLCreationRequest{url})
	expect.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/url", bytes.NewReader(body))

	app.SetupRouter().Router().ServeHTTP(w, req)
	expect.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreate_FailIfWrongRequestBody(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	app := pkg.NewApp(&pkg.Config{
		Host:     ":8080",
		MongoURL: "mongodb://localhost:4343",
		Redis: pkg.RedisConfig{
			URL: mr.Addr(),
		},
	})
	require.NoError(t, app.Init())

	app.UrlsDAO = new(UrlMappingDAOMock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/url", bytes.NewReader([]byte("wrong-input")))

	app.SetupRouter().Router().ServeHTTP(w, req)
	expect.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestCreate_FailAddRequest(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	app := pkg.NewApp(&pkg.Config{
		Host:     ":8080",
		MongoURL: "mongodb://localhost:4343",
		Redis: pkg.RedisConfig{
			URL: mr.Addr(),
		},
	})
	require.NoError(t, app.Init())

	dao := new(UrlMappingDAOMock)
	dao.On("Add", mock.Anything, mock.Anything).Return(errors.New("failed DB connection"))

	app.UrlsDAO = dao
	url := "https://google.com"

	body, err := json.Marshal(handlers.URLCreationRequest{url})
	expect.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/url", bytes.NewReader(body))

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusBadGateway, w.Code)
}
