package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CristianCurteanu/url-shortener/pkg"
	"github.com/stretchr/testify/mock"
	expect "github.com/stretchr/testify/require"
)

func TestDelete_Success(t *testing.T) {
	app := pkg.NewApp(&pkg.Config{
		Host: ":8080",
	})

	key := "aNjR4bEV"
	dao := new(UrlMappingDAOMock)
	dao.On("Delete", mock.Anything, key).Return(nil)
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/url/"+key, nil)

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusOK, w.Code)
	expect.Contains(t, string(w.Body.Bytes()), "OK")
	expect.True(t, dao.AssertCalled(t, "Delete", mock.Anything, key))
}

func TestDelete_FailDeleteRequest(t *testing.T) {
	app := pkg.NewApp(&pkg.Config{
		Host: ":8080",
	})

	key := "aNjR4bEV"
	dao := new(UrlMappingDAOMock)
	dao.On("Delete", mock.Anything, key).Return(errors.New("not found"))
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/url/"+key, nil)

	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusBadGateway, w.Code)
	expect.Contains(t, string(w.Body.Bytes()), "mapping_not_deleted")
	expect.True(t, dao.AssertCalled(t, "Delete", mock.Anything, key))
}

func TestDelete_FailEmptyKey(t *testing.T) {
	app := pkg.NewApp(&pkg.Config{
		Host: ":8080",
	})
	key := ""
	dao := new(UrlMappingDAOMock)
	dao.On("Delete", mock.Anything, key).Return(errors.New("not found"))
	app.UrlsDAO = dao

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/url/"+key, nil)
	app.SetupRouter().Router().ServeHTTP(w, req)

	expect.Equal(t, http.StatusNotFound, w.Code)
}
