package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CristianCurteanu/url-shortener/pkg"
	"github.com/CristianCurteanu/url-shortener/pkg/urls"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	expect "github.com/stretchr/testify/require"
)

func TestRead_Success(t *testing.T) {
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

	key := "aNjR4bEV"
	dao := new(UrlMappingDAOMock)
	dao.On("SearchById", mock.Anything, key).Return(urls.UrlMapping{
		Key:     key,
		URL:     "http://google.com",
		Counter: 0,
	}, nil)
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/url/"+key, nil)

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusOK, w.Code)
	expect.Contains(t, string(w.Body.Bytes()), "google.com")
	expect.Contains(t, string(w.Body.Bytes()), key)
	expect.True(t, dao.AssertCalled(t, "SearchById", mock.Anything, key))
}

func TestRead_FailDeleteRequest(t *testing.T) {
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

	key := "aNjR4bEV"
	dao := new(UrlMappingDAOMock)
	dao.On("SearchById", mock.Anything, key).Return(urls.UrlMapping{}, errors.New("not found"))
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/url/"+key, nil)

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusNotFound, w.Code)
	expect.Contains(t, string(w.Body.Bytes()), "mapping_not_found")
	expect.True(t, dao.AssertCalled(t, "SearchById", mock.Anything, key))
}

func TestRead_FailEmptyKey(t *testing.T) {
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

	key := ""
	dao := new(UrlMappingDAOMock)
	dao.On("SearchById", mock.Anything, key).Return(urls.UrlMapping{}, errors.New("not found"))
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/url/"+key, nil)
	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusNotFound, w.Code)
}
