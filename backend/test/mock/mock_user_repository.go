package mock

import (
	"github.com/ahmadzainulhasan/optik-pelangi/backend/internal/entity"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserRepository struct {
	Mock mock.Mock
}


func (m *MockUserRepository) FindByEmail(tx *gorm.DB, email string) (entity.User, error) {
	args := m.Mock.Called(tx, email)
	if args.Get(0) == nil {
		return entity.User{}, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(tx *gorm.DB, id any) (entity.User, error) {
	args := m.Mock.Called(tx, id)
	if args.Get(0) == nil {
		return entity.User{}, args.Error(1)
	}
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) Save(tx *gorm.DB, user entity.User) error {
	args := m.Mock.Called(tx, user)
	return args.Error(0)
}
