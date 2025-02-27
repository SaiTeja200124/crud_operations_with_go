package mocks

import (
	"fmt"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB implements the Database interface
type MockDB struct {
	mock.Mock
	Error error
}

// func (m *MockDB) Create(value interface{}) *gorm.DB {
// 	args := m.Called(value)
// 	if args.Get(0) != nil {
// 		return args.Get(0).(*gorm.DB)
// 	}
// 	return &gorm.DB{}
// }

// func (m *MockDB) Create(value interface{}) *gorm.DB {
// 	args := m.Called(value)
// 	if err, ok := args.Get(0).(error); ok {
// 		return &gorm.DB{Error: err}
// 	}
// 	return &gorm.DB{} // Return an empty GORM DB instance if no error
// }

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	if err, ok := args.Get(0).(error); ok {
		return &gorm.DB{Error: err}
	}
	return &gorm.DB{}
}
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)

	if db, ok := args.Get(0).(*gorm.DB); ok {
		return db
	}

	return &gorm.DB{Error: args.Error(1)}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	fmt.Printf("[MockDB] Where() called with query: %v, args: %v\n", query, args)
	m.Called(query, args)
	return &gorm.DB{} // ✅ Always return a valid *gorm.DB instance
}

// Mocking First method
// Mock First method
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where[0]) // Use where[0] to pass only the ID as an int

	if db, ok := args.Get(0).(*gorm.DB); ok {
		fmt.Println("MockDB First called, returning valid *gorm.DB instance")
		return db
	}

	return &gorm.DB{}
}

// Mocking Save method
func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	if db, ok := args.Get(0).(*gorm.DB); ok {
		return db
	}
	return &gorm.DB{} // Ensure a valid GORM instance is always returned
}

func (m *MockDB) Delete(value interface{}) *gorm.DB {
	args := m.Called(value)
	if db, ok := args.Get(0).(*gorm.DB); ok {
		return db
	}
	return &gorm.DB{}
}
