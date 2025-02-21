package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"connection_to_pg/mocks"
	"connection_to_pg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// type MockDB struct {
// 	mock.Mock
// }

// func (m *MockDB) Create(value interface{}) *gorm.DB {
// 	args := m.Called(value)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(*gorm.DB)
// 	}
// 	return &gorm.DB{}
// }

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
