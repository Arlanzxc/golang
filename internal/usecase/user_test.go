package usecase

import (
	"fmt"
	"testing"

	"GOLANG/internal/entity"
	"GOLANG/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserWithEmailCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	uc := NewUserUseCase(mockRepo)
	user := &entity.User{Username: "Test", Email: "test@mail.com"}

	mockRepo.EXPECT().GetByEmail("test@mail.com").Return(&entity.User{}, nil)
	err := uc.RegisterUserWithEmailCheck(user, "test@mail.com")
	assert.EqualError(t, err, "user with this email already exists")

	mockRepo.EXPECT().GetByEmail("new@mail.com").Return(nil, fmt.Errorf("user not found"))
	mockRepo.EXPECT().RegisterUser(user).Return(user, nil)
	err = uc.RegisterUserWithEmailCheck(user, "new@mail.com")
	assert.NoError(t, err)

	mockRepo.EXPECT().GetByEmail("err@mail.com").Return(nil, fmt.Errorf("user not found"))
	mockRepo.EXPECT().RegisterUser(user).Return(nil, fmt.Errorf("db error"))
	err = uc.RegisterUserWithEmailCheck(user, "err@mail.com")
	assert.EqualError(t, err, "db error")
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	uc := NewUserUseCase(mockRepo)
	user := &entity.User{Username: "OldName"}

	err := uc.UpdateUserName("123", "")
	assert.EqualError(t, err, "name cannot be empty")

	mockRepo.EXPECT().GetUserByID("123").Return(nil, fmt.Errorf("not found"))
	err = uc.UpdateUserName("123", "NewName")
	assert.EqualError(t, err, "not found")

	mockRepo.EXPECT().GetUserByID("123").Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *entity.User) error {
		assert.Equal(t, "NewName", u.Username) 
		return nil
	})
	err = uc.UpdateUserName("123", "NewName")
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	uc := NewUserUseCase(mockRepo)

	err := uc.DeleteUser("1")
	assert.EqualError(t, err, "it is not allowed to delete admin user")

	mockRepo.EXPECT().DeleteUser("123").Return(nil) 
	err = uc.DeleteUser("123")
	assert.NoError(t, err)

	mockRepo.EXPECT().DeleteUser("999").Return(fmt.Errorf("db error"))
	err = uc.DeleteUser("999")
	assert.EqualError(t, err, "db error")
}