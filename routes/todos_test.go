package routes

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"project_todo/models"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTodosSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockDate := time.Now().Truncate(time.Second)
	mockUserId := int64(10)

	mockTodos := []models.Todo{
		models.Todo{
			ID:    1,
			Title: "Test todo1",
			List: []models.TodoItem{
				models.TodoItem{
					Item:    "sample data",
					Checked: false,
				},
			},
			IsActive:  true,
			CreatedAt: mockDate,
			UpdatedAt: mockDate,
			UserID:    10,
		},
		models.Todo{
			ID:    2,
			Title: "Test todo2",
			List: []models.TodoItem{
				models.TodoItem{
					Item:    "sample data2",
					Checked: false,
				},
			},
			IsActive:  true,
			CreatedAt: mockDate,
			UpdatedAt: mockDate,
			UserID:    10,
		},
	}
	monkey.Patch(models.GetAllTodos, func(userId int64) ([]models.Todo, error) {
		t.Log("userId", userId)
		if userId == 10 {
			return mockTodos, nil
		}
		return []models.Todo{}, nil
	})
	defer monkey.Unpatch(models.GetAllTodos)

	c.Set("userId", mockUserId)
	getAllTodos(c)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `[
		{
		  "id": 1,
		  "title": "Test todo1",
		  "list": [
			{
			  "item": "sample data",
			  "checked": false
			}
		  ],
		  "isActive": true,
		  "createdAt":  "` + mockDate.Format(time.RFC3339Nano) + `",
		  "updatedAt":  "` + mockDate.Format(time.RFC3339Nano) + `",
		  "userId": 10
		},
		{
		  "id": 2,
		  "title": "Test todo2",
		  "list": [
			{
			  "item": "sample data2",
			  "checked": false
			}
		  ],
		  "isActive": true,
		  "createdAt":  "` + mockDate.Format(time.RFC3339Nano) + `",
		  "updatedAt":  "` + mockDate.Format(time.RFC3339Nano) + `",
		  "userId": 10
		}
	  ]`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetAllTodos_FetchError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch models.GetAllTodos to simulate a fetch error
	monkey.Patch(models.GetAllTodos, func(userId int64) ([]models.Todo, error) {
		return nil, errors.New("fetch error")
	})
	defer monkey.Unpatch(models.GetAllTodos)

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request
	c.Request = httptest.NewRequest("GET", "/todos", nil)

	// Call the handler function
	getAllTodos(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to fetch todos"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCreateTodoSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock current time to ensure consistency
	mockTime := time.Date(2024, time.August, 26, 18, 2, 39, 0, time.Local)
	monkey.Patch(time.Now, func() time.Time {
		return mockTime
	})

	mockTodo := models.Todo{
		Title:     "Test todo",
		List:      []models.TodoItem{{Item: "sample data", Checked: false}},
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    10,
	}

	// Patch the Save method
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.Todo{}), "Save", func(t *models.Todo) error {
		// Simulate a successful save by setting the ID
		t.ID = 1
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.Todo{}), "Save")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate JSON request body
	c.Request = httptest.NewRequest("POST", "/todos", io.NopCloser(bytes.NewBufferString(`{
			"title": "Test todo",
			"list": [{"item": "sample data", "checked": false}]
		}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	createTodo(c)

	// Assert that the response status code is 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)

	// Assert the response body
	expectedResponse := `{
			"message": "Todo created",
			"todo": {
				"id": 1,
				"title": "Test todo",
				"list": [{"item": "sample data", "checked": false}],
				"isActive": true,
				"createdAt": "` + mockTodo.CreatedAt.Format(time.RFC3339Nano) + `",
				"updatedAt": "` + mockTodo.UpdatedAt.Format(time.RFC3339Nano) + `",
				"userId": 10
			}
		}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCreateTodo_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the request with invalid JSON body
	c.Request = httptest.NewRequest("POST", "/todos", io.NopCloser(bytes.NewBufferString(`{
		"title": "Test Todo"  // Missing closing brace
	`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	createTodo(c)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to parse request data"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestCreateTodo_SaveError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock the current time
	mockTime := time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC)
	monkey.Patch(time.Now, func() time.Time {
		return mockTime
	})
	defer monkey.Unpatch(time.Now)

	// Patch the Save method on *models.Todo to simulate a save error
	monkey.PatchInstanceMethod(reflect.TypeOf(&models.Todo{}), "Save", func(t *models.Todo) error {
		return errors.New("save error")
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&models.Todo{}), "Save")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request with JSON body
	c.Request = httptest.NewRequest("POST", "/todos", io.NopCloser(bytes.NewBufferString(`{
		"title": "Test Todo",
		"list": [{"item": "sample data", "checked": false}]
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	createTodo(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to create todo"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestDeleteTodoById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:       1,
		Title:    "Test todo",
		List:     []models.TodoItem{},
		IsActive: true,
		UserID:   10,
	}

	// Patch models.GetTodoById
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		t.Log("Actual todo id ", todoId)
		if todoId == mockTodo.ID {
			return mockTodo, nil
		}
		return nil, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Patch the Delete method on *models.Todo
	monkey.PatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Delete", func(td models.Todo) error {
		// Simulate a successful deletion without accessing the database
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Delete")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("DELETE", "/todos/1", nil)

	// Call the handler function
	deleteTodoById(c)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Todo deleted successfully"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestDeleteTodoById_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:     1,
		UserID: 20, // Different user ID to simulate unauthorized access
	}

	// Patch models.GetTodoById
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context to a different user
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("DELETE", "/todos/1", nil)

	// Call the handler function
	deleteTodoById(c)

	// Assert that the response status code is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unauthorized delete"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestDeleteTodoById_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:     1,
		UserID: 10,
	}

	// Patch models.GetTodoById
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Patch the Delete method on *models.Todo to simulate an error
	monkey.PatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Delete", func(t models.Todo) error {
		return errors.New("delete error") // Simulate a delete failure
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Delete")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("DELETE", "/todos/1", nil)

	// Call the handler function
	deleteTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to delete the todo"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetTodoById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockTime := time.Date(2024, time.August, 26, 18, 2, 39, 0, time.Local)
	// Mock Todo data
	mockTodo := &models.Todo{
		ID:    1,
		Title: "Test Todo",
		List: []models.TodoItem{
			models.TodoItem{
				Item:    "sample data",
				Checked: true,
			},
		},
		IsActive:  true,
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
		UserID:    10,
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("GET", "/todos/1", nil)

	// Call the handler function
	getTodoById(c)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	expectedResponse := `{
		"id": 1,
		"title": "Test Todo",
		"list": [
			{
			  "item": "sample data",
			  "checked": true
			}
		  ],
		"isActive": true,
		"createdAt":  "` + mockTime.Format(time.RFC3339Nano) + `",
		"updatedAt":  "` + mockTime.Format(time.RFC3339Nano) + `",
		"userId": 10
		
	}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetTodoById_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:     1,
		Title:  "Test Todo",
		UserID: 20, // Different user ID to simulate unauthorized access
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context to a different user
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("GET", "/todos/1", nil)

	// Call the handler function
	getTodoById(c)

	// Assert that the response status code is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unauthorized access"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetTodoById_ParseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the request with an invalid ID
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	c.Request = httptest.NewRequest("GET", "/todos/invalid", nil)

	// Call the handler function
	getTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to parse todo id"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetTodoById_FetchError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch models.GetTodoById to simulate a fetch error
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return nil, errors.New("fetch error")
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/todos/1", nil)

	// Call the handler function
	getTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to fetch todo"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockTime := time.Date(2024, time.August, 26, 18, 2, 39, 0, time.Local)
	// Mock Todo data
	mockTodo := &models.Todo{
		ID:    1,
		Title: "Test Todo",
		List: []models.TodoItem{
			models.TodoItem{
				Item:    "sample data",
				Checked: true,
			},
		},
		IsActive:  true,
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
		UserID:    10,
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Patch the Update method on *models.Todo
	monkey.PatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Update", func(t models.Todo) error {
		// Simulate a successful update
		return nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Update")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request with JSON body
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("PUT", "/todos/1", io.NopCloser(bytes.NewBufferString(`{
		"title": "Updated Title"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Todo updated successfully"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_ParseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the request with an invalid ID
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}
	c.Request = httptest.NewRequest("PUT", "/todos/invalid", nil)

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to parse todo id"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:     1,
		Title:  "Original Title",
		UserID: 20, // Different user ID to simulate unauthorized access
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context to a different user
	c.Set("userId", int64(10))

	// Simulate the request with JSON body
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("PUT", "/todos/1", io.NopCloser(bytes.NewBufferString(`{
		"title": "Updated Title"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unauthorized update"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_FetchError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Patch models.GetTodoById to simulate a fetch error
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return nil, errors.New("fetch error")
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request with JSON body
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/todos/1", io.NopCloser(bytes.NewBufferString(`{
		"title": "Updated Title"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to fetch todo to update"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_BindJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock Todo data
	mockTodo := &models.Todo{
		ID:     1,
		Title:  "Original Title",
		UserID: 10,
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request with invalid JSON body
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("PUT", "/todos/1", io.NopCloser(bytes.NewBufferString(`{
		"title": "Updated Title"  // Missing closing quotes
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to parse request data"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestUpdateTodoById_UpdateError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockTime := time.Date(2024, time.August, 26, 18, 2, 39, 0, time.Local)
	// Mock Todo data
	mockTodo := &models.Todo{
		ID:    1,
		Title: "Test Todo",
		List: []models.TodoItem{
			models.TodoItem{
				Item:    "sample data",
				Checked: true,
			},
		},
		IsActive:  true,
		CreatedAt: mockTime,
		UpdatedAt: mockTime,
		UserID:    10,
	}

	// Patch models.GetTodoById to return the mock todo
	monkey.Patch(models.GetTodoById, func(todoId int64) (*models.Todo, error) {
		return mockTodo, nil
	})
	defer monkey.Unpatch(models.GetTodoById)

	// Patch the Update method on *models.Todo to simulate an update error
	monkey.PatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Update", func(t models.Todo) error {
		return errors.New("update error")
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(models.Todo{}), "Update")

	// Set the userId in the context
	c.Set("userId", int64(10))

	// Simulate the request with JSON body
	c.Params = gin.Params{{Key: "id", Value: strconv.FormatInt(mockTodo.ID, 10)}}
	c.Request = httptest.NewRequest("PUT", "/todos/1", io.NopCloser(bytes.NewBufferString(`{
		"title": "Updated Title"
	}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call the handler function
	updateTodoById(c)

	// Assert that the response status code is 500 Internal Server Error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert the response body
	expectedResponse := `{"message":"Unable to update todo"}`
	assert.JSONEq(t, expectedResponse, w.Body.String())
}
