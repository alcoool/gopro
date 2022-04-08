package controller

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"example.com/mod/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestController_Seed(t *testing.T) {
	//not implemented due to proof of skill in testing Login/Logout
}

func TestController_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		database := db.DataBaseMock{}
		handler := Router(gin.New(), &database)

		data := url.Values{}
		data.Set("email", "jon.doe@company.com")
		data.Set("password", "1234")

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		assert.Nil(t, err, "No error should be given")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody := rr.Body.String()

		assert.Equal(t, "{\"message\":\"success\"}", responseBody)
	})

	t.Run("wrong email", func(t *testing.T) {
		database := db.DataBaseMock{}
		handler := Router(gin.New(), &database)

		data := url.Values{}
		data.Set("email", "bad.email@company.com")
		data.Set("password", "1234")

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		assert.Nil(t, err, "No error should be given")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody := rr.Body.String()

		assert.Equal(t, "{\"message\":\"not allowed\"}", responseBody)
	})

	t.Run("wrong password", func(t *testing.T) {
		database := db.DataBaseMock{}
		handler := Router(gin.New(), &database)

		data := url.Values{}
		data.Set("email", "jon.doe@company.com")
		data.Set("password", "12345")

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		assert.Nil(t, err, "No error should be given")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody := rr.Body.String()

		assert.Equal(t, "{\"message\":\"not allowed\"}", responseBody)
	})

	t.Run("db error", func(t *testing.T) {
		database := db.DataBaseMock{}
		handler := Router(gin.New(), &database)

		data := url.Values{}
		data.Set("email", "error")
		data.Set("password", "error")

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		assert.Nil(t, err, "No error should be given")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody := rr.Body.String()

		assert.Equal(t, "{\"message\":\"could not login\"}", responseBody)
	})
}

func TestController_Logout(t *testing.T) {
	t.Run("successful logout", func(t *testing.T) {
		database := db.DataBaseMock{}
		handler := Router(gin.New(), &database)

		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/logout", strings.NewReader(""))

		assert.Nil(t, err, "No error should be given")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		responseBody := rr.Body.String()

		assert.Equal(t, "{\"message\":\"success\"}", responseBody)
	})
}

func TestController_AddToCart(t *testing.T) {
	//not implemented due to proof of skill in testing Login/Logout
}

func TestController_Checkout(t *testing.T) {
	//not implemented due to proof of skill in testing Login/Logout
}