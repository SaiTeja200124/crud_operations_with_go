package handlers

import (
	"connection_to_pg/db"
	"connection_to_pg/models"
	"fmt"

	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	// "github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

// Define an interface for database operations
type Database interface {
	Create(value interface{}) *gorm.DB
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

func GetAll(w http.ResponseWriter, _ *http.Request) {
	var books []models.Book
	if err := db.DB.Find(&books).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error querying books table %v", err)
		return
	}

	// Marshalling books to JSON
	j, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling books into json %v\n", err)
		return
	}

	w.Write(j)
}

func Get(w http.ResponseWriter, r *http.Request) {
	searchQuery := chi.URLParam(r, "query") // Extract search term from URL

	if searchQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing search query"}`))
		return
	}

	var books []models.Book
	if err := db.DB.Where("name ILIKE ? OR description ILIKE ? OR author ILIKE ?",
		"%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%").Find(&books).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error querying books table: %v", err)
		return
	}

	if len(books) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "No books found"}`))
		return
	}

	j, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling books into json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func Update(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing %d from string into integer %v", bookID, err)
		return
	}

	var book models.Book
	if err := db.DB.First(&book, bookID).Error; err != nil {
		// if gorm.IsRecordNotFoundError(err) {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	log.Printf("book with id %d not found", bookID)
		// } else {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error querying book from books table with id %d %v", bookID, err)
		// }
		return
	}

	// Decode request body into a map (to check for missing fields)
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body %v", err)
		return
	}

	// Conditionally update fields
	if name, exists := body["name"].(string); exists {
		book.Name = name
	}
	if desc, exists := body["description"].(string); exists {
		book.Description = desc
	}
	if author, exists := body["author"].(string); exists {
		book.Author = author
	}

	// Save updated book
	if err := db.DB.Save(&book).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error updating book %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("book with id %d updated successfully", bookID)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID")) // Convert bookID from string to int
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing book ID from string to integer: %v", err)
		w.Write([]byte(`{"error": "Invalid book ID"}`))
		return
	}

	// Check if the book exists before attempting to delete
	var book models.Book
	if err := db.DB.First(&book, bookID).Error; err != nil {
		// if gorm.IsRecordNotFoundError(err) {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	log.Printf("book with ID %d not found", bookID)
		// 	w.Write([]byte(`{"error": "Book not found"}`))
		// } else {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error querying book from database: %v", err)
		// }
		return
	}

	// Delete the book
	if err := db.DB.Delete(&book).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error deleting book with ID %d: %v", bookID, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book deleted successfully"}`))
}
