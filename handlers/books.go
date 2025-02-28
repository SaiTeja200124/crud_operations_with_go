package handlers

import (
	// "connection_to_pg/db"
	"connection_to_pg/models"
	"errors"
	"fmt"
	"strconv"

	"encoding/json"
	"log"
	"net/http"

	// "strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Define an interface for database operations
type Database interface {
	Create(value interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
}

// Handler struct now depends on the interface, not on *gorm.DB directly
type Handler struct {
	DB Database
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ✅ Fix: Ensure result.Error is checked properly
	result := h.DB.Create(&book)
	fmt.Printf("Result: %+v\n", result)            // Print the entire result for inspection
	fmt.Printf("Result.Error: %v\n", result.Error) // Print the error specifically

	if result.Error != nil {
		fmt.Println("Error block entered") // Make sure this block is executed
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Success block entered")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book created successfully"})
}

func (h *Handler) GetAll(w http.ResponseWriter, _ *http.Request) {
	var books []models.Book

	// Query the database
	if err := h.DB.Find(&books).Error; err != nil {
		http.Error(w, `{"error": "Failed to retrieve books"}`, http.StatusInternalServerError)
		log.Printf("error querying books table: %v", err)
		return
	}

	// Marshal books to JSON
	j, _ := json.Marshal(books)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

var jsonMarshal = json.Marshal

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "query")
	fmt.Println("Received ID:", idStr)

	id, err := strconv.Atoi(idStr) // Convert searchQuery to an integer
	if err != nil {
		http.Error(w, `{"error": "Invalid ID format"}`, http.StatusBadRequest)
		return
	}

	if h.DB == nil {
		panic("h.DB is nil - database connection not initialized")
	}
	fmt.Println("Database connection initialized")

	var book models.Book
	if err := h.DB.First(&book, id).Error; err != nil {
		http.Error(w, `{"message": "Book not found"}`, http.StatusNotFound)
		return
	}
	fmt.Println("Book found:", book)

	j, err := jsonMarshal(book)
	if err != nil {
		http.Error(w, "Failed to marshal book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract and validate book ID from URL
	bookIDParam := chi.URLParam(r, "id") // Ensure this matches the test case
	bookID, err := strconv.Atoi(bookIDParam)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	// Find the book in the database
	var book models.Book
	if err := h.DB.First(&book, bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "Book not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Decode request body
	var updateData models.Book
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Update book fields
	book.Name = updateData.Name
	book.Description = updateData.Description
	book.Author = updateData.Author

	// Save updated book
	if err := h.DB.Save(&book).Error; err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book updated successfully"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract and validate book ID from URL
	bookIDParam := chi.URLParam(r, "id") // Ensure this matches the URL parameter
	bookID, err := strconv.Atoi(bookIDParam)
	if err != nil {
		http.Error(w, `{"error": "Invalid book ID"}`, http.StatusBadRequest)
		log.Printf("error parsing book ID from string to integer: %v", err)
		return
	}

	// Check if the book exists before attempting to delete
	var book models.Book
	if err := h.DB.First(&book, bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "Book not found"}`, http.StatusNotFound)
			log.Printf("book with ID %d not found", bookID)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Printf("error querying book from database: %v", err)
		return
	}

	// Delete the book
	if err := h.DB.Delete(&book).Error; err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		log.Printf("error deleting book with ID %d: %v", bookID, err)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Book deleted successfully"})
}
