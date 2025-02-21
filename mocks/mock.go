package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB implements the Database interface
type MockDB struct {
	mock.Mock
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
	err, _ := args.Get(0).(error) // Get error if available

	return &gorm.DB{Error: err} // Return gorm.DB with error
}
