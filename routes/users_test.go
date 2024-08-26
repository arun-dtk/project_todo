package routes

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"project_todo/models"
	"project_todo/utils"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignup_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock utils.HashPassword to return a dummy hashed password
	monkey.Patch(utils.HashPassword, func(password string) (string, error) {
		return "hashedPassword", nil
	})
	defer monkey.Unpatch(utils.HashPassword)

	// Patch the Save method on *models.User
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.User{}), "Save", func(u *models.User) error {
		// Simulate a successful save
		u.ID = 1
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.User{}), "Save")

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/signup", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"firstName": "John",
		"lastName": "Doe",
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	signup(c)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"User created successfully"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestSignup_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the request with invalid JSON body
	c.Request = httptest.NewRequest("POST", "/signup", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"firstName": "John",  // Missing closing brace
		"lastName": "Doe"
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	signup(c)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Could not parse the request body"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestSignup_SaveError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock utils.HashPassword to return a dummy hashed password
	monkey.Patch(utils.HashPassword, func(password string) (string, error) {
		return "hashedPassword", nil
	})
	defer monkey.Unpatch(utils.HashPassword)

	// Patch the Save method on *models.User to simulate a save error
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.User{}), "Save", func(u *models.User) error {
		return errors.New("save error")
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.User{}), "Save")

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/signup", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"firstName": "John",
		"lastName": "Doe",
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	signup(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to save the user"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch the ValidateCredentials method on *models.User to simulate successful validation
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials", func(u *models.User) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials")

	// Patch utils.GenerateToken to return a dummy token
	monkey.Patch(utils.GenerateToken, func(email string, userID int64) (string, error) {
		return "dummyToken", nil
	})
	defer monkey.Unpatch(utils.GenerateToken)

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/login", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	login(c)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"User logged in successfully","token":"dummyToken"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestLogin_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the request with invalid JSON body
	c.Request = httptest.NewRequest("POST", "/login", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"password": "password123"  // Missing closing brace
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	login(c)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Could not parse the request body"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestLogin_ValidateCredentialsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch the ValidateCredentials method on *models.User to simulate validation error
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials", func(u *models.User) error {
		return errors.New("Invalid Credentials")
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials")

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/login", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	login(c)

	// Assert that the response status code is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to authenticate the user"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestLogin_GenerateTokenError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch the ValidateCredentials method on *models.User to simulate successful validation
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials", func(u *models.User) error {
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.User{}), "ValidateCredentials")

	// Patch utils.GenerateToken to simulate token generation error
	monkey.Patch(utils.GenerateToken, func(email string, userID int64) (string, error) {
		return "", errors.New("token generation error")
	})
	defer monkey.Unpatch(utils.GenerateToken)

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/login", io.NopCloser(bytes.NewBufferString(`{
		"email": "testuser@example.com",
		"password": "password123"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	login(c)

	// Assert that the response status code is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to authenticate the user"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
