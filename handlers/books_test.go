package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"connection_to_pg/mocks"
	"connection_to_pg/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TEST CASES FOR CREATE OPERATION
func TestCreateBookSuccess(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	// Sample book data
	book := models.Book{Name: "Test Book", Description: "A test book", Author: "Test Author"}
	bookJSON, _ := json.Marshal(book)

	// Mock DB behavior
	mockDB.On("Create", mock.AnythingOfType("*models.Book")).Return(mockDB)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.Create(rr, req)

	// Validate response
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, `{"message": "Book created successfully"}`, rr.Body.String())

	// Assert that Create was called
	mockDB.AssertExpectations(t)
}

func TestCreateBookBadRequest(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	// Invalid JSON input
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"invalid"`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateBook_DatabaseError(t *testing.T) {
	mockDB := new(mocks.MockDB)
	h := &Handler{DB: mockDB}

	book := models.Book{Name: "Test Book", Description: "A test book", Author: "Test Author"}
	bookJSON, _ := json.Marshal(book)

	mockDB.On("Create", mock.Anything).Return(errors.New("Database Error"))

	r := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(bookJSON))
	w := httptest.NewRecorder()

	h.Create(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TEST CASES FOR READ OPERATION
func TestGetAll_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	// Mock data
	expectedBooks := []models.Book{
		{ID: 1, Name: "Book One", Author: "Author One"},
		{ID: 2, Name: "Book Two", Author: "Author Two"},
	}

	// Corrected mock
	mockDB.On("Find", mock.AnythingOfType("*[]models.Book"), mock.Anything).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]models.Book) // Correctly get the first argument
		*arg = expectedBooks                // Assign mock data to the argument
	}).Return(&gorm.DB{}) // Simulate successful DB query

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	// Call the GetAll function
	handler.GetAll(w, req)
	fmt.Println("my code======>", w.Code)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	// assert.JSONEq(t, `[{"ID":1,"Title":"Book One","Author":"Author One"},{"ID":2,"Title":"Book Two","Author":"Author Two"}]`, w.Body.String())

	// Verify mock expectations
	mockDB.AssertExpectations(t)
}

func TestGetAll_DatabaseError(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	// Mock the Find method to return an error
	mockDB.On("Find", mock.Anything, mock.Anything).
		Return(nil, errors.New("database error"))

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	// Call the GetAll function
	handler.GetAll(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve books")

	// Verify mock expectations
	mockDB.AssertExpectations(t)
}

func TestGet_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	book := models.Book{ID: 1, Name: "Test Book", Description: "A test book", Author: "Author Name"}

	mockDB.On("First", mock.AnythingOfType("*models.Book"), 1).
		Run(func(args mock.Arguments) {
			argBook := args.Get(0).(*models.Book)
			*argBook = book
		}).
		Return(&gorm.DB{}) // Return a valid *gorm.DB instance

	handler := Handler{DB: mockDB}

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("query", "1") // Use "1" to match the integer conversion

	r := httptest.NewRequest("GET", "/books/1", nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()
	handler.Get(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBook models.Book
	err := json.NewDecoder(resp.Body).Decode(&responseBook)
	assert.NoError(t, err)
	assert.Equal(t, book, responseBook)

	mockDB.AssertExpectations(t)
}

func TestGet_MissingQuery(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	// Call handler without setting chi param
	handler.Get(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	fmt.Println("error message =====>", w.Body.String())
	assert.Contains(t, w.Body.String(), `"error": "Invalid ID format"`)
}

func TestGet_InvalidIDFormat(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	invalidID := "abc"

	req, err := http.NewRequest("GET", "/books/"+invalidID, nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("query", invalidID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Get(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error": "Invalid ID format"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}

func TestGet_DBNotInitialized(t *testing.T) {
	handler := Handler{DB: nil} // DB is not initialized

	validID := "1"

	req, err := http.NewRequest("GET", "/books/"+validID, nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("query", validID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()

	// Expecting a panic, so we need to recover from it
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "h.DB is nil - database connection not initialized", r)
		}
	}()

	handler.Get(rr, req)
}

// TEST CASES FOR UPDATE OPERATION
func TestUpdate_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	bookID := 1
	existingBook := models.Book{
		ID:          bookID,
		Name:        "Old Name",
		Description: "Old Desc",
		Author:      "Old Author",
	}

	mockDB.On("First", mock.AnythingOfType("*models.Book"), bookID).
		Run(func(args mock.Arguments) {
			argBook := args.Get(0).(*models.Book)
			*argBook = existingBook
		}).Return(&gorm.DB{})

	mockDB.On("Save", mock.AnythingOfType("*models.Book")).Return(&gorm.DB{})

	updateData := models.Book{
		Name:        "New Name",
		Description: "New Desc",
		Author:      "New Author",
	}
	jsonBody, _ := json.Marshal(updateData)

	req, err := http.NewRequest("PUT", "/books/"+strconv.Itoa(bookID), bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(bookID)) // Ensure this matches the Update method
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"message": "Book updated successfully"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}

func TestUpdate_Failure(t *testing.T) {
	// Create a new Chi router
	router := chi.NewRouter()

	// Create a new mock database
	mockDB := new(mocks.MockDB)
	handler := &Handler{DB: mockDB}

	// Define the book ID as an int (as it will be converted in the handler)
	id := 99

	// Mock the `First` method (simulate record retrieval failure)
	mockDB.On("First", mock.AnythingOfType("*models.Book"), id).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	// Mock the `Save` method (should not be called in case of failure)
	mockDB.On("Save", mock.Anything).Return(&gorm.DB{Error: errors.New("unexpected Save call")}).Maybe()

	// Define the update route
	router.Put("/books/{id}", handler.Update)

	// Create a test request
	req, err := http.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(id), nil)
	require.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assertions
	require.Equal(t, http.StatusNotFound, w.Code) // Expecting 404 instead of 400
	require.Contains(t, w.Body.String(), "Book not found")

	// Ensure all expectations were met
	mockDB.AssertExpectations(t)
}
func TestUpdate_InvalidJSON(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	bookID := 1

	mockDB.On("First", mock.AnythingOfType("*models.Book"), bookID).Return(&gorm.DB{})

	invalidJSON := `{"name": "New Name", "description": "New Desc",` // Missing closing brace

	req, err := http.NewRequest("PUT", "/books/"+strconv.Itoa(bookID), bytes.NewBufferString(invalidJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(bookID)) // Ensure this matches the Update method
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Update(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error": "Invalid request payload"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}

//TEST CASES FOR DELETE OPERATION

func TestDelete_Success(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	bookID := 1
	existingBook := models.Book{
		ID:          bookID,
		Name:        "Sample Book",
		Description: "Sample Description",
		Author:      "Sample Author",
	}

	mockDB.On("First", mock.AnythingOfType("*models.Book"), bookID).
		Run(func(args mock.Arguments) {
			argBook := args.Get(0).(*models.Book)
			*argBook = existingBook
		}).Return(&gorm.DB{})

	mockDB.On("Delete", mock.AnythingOfType("*models.Book")).Return(&gorm.DB{})

	req, err := http.NewRequest("DELETE", "/books/"+strconv.Itoa(bookID), nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(bookID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"message": "Book deleted successfully"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}

func TestDelete_BookNotFound(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	bookID := 99

	mockDB.On("First", mock.AnythingOfType("*models.Book"), bookID).
		Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

	req, err := http.NewRequest("DELETE", "/books/"+strconv.Itoa(bookID), nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(bookID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.JSONEq(t, `{"error": "Book not found"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}

func TestDelete_InvalidBookID(t *testing.T) {
	mockDB := new(mocks.MockDB)
	handler := Handler{DB: mockDB}

	invalidBookID := "abc"

	req, err := http.NewRequest("DELETE", "/books/"+invalidBookID, nil)
	assert.NoError(t, err)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", invalidBookID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.Delete(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error": "Invalid book ID"}`, rr.Body.String())

	mockDB.AssertExpectations(t)
}
