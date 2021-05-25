package validator

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reecerussell/adaptive-password-hasher/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
)

func TestNewUserValidator(t *testing.T) {
	v := NewUserValidator()

	assert.NotNil(t, v.(*userValidator).h)
}

func TestUserValidator_ValidatePassword_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testPassword := "myPassword1"
	testUser := &dal.User{
		PasswordHash: "aGVsbG8K",
	}

	mockHasher := mock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(true)

	v := &userValidator{h: mockHasher}

	err := v.ValidatePassword(testUser, testPassword)
	assert.NoError(t, err)
}

func TestUserValidator_ValidatePassword_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Given Invalid Password", func(t *testing.T) {
		testPassword := "myPassword1"
		testUser := &dal.User{
			PasswordHash: "aGVsbG8K",
		}

		mockHasher := mock.NewMockHasher(ctrl)
		mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(false)

		v := &userValidator{h: mockHasher}

		err := v.ValidatePassword(testUser, testPassword)
		assert.Equal(t, ErrInvalidPassword, err)
	})

	t.Run("Given Empty Password", func(t *testing.T) {
		testPassword := ""
		testUser := &dal.User{
			PasswordHash: "aGVsbG8K",
		}

		mockHasher := mock.NewMockHasher(ctrl)
		v := &userValidator{h: mockHasher}

		err := v.ValidatePassword(testUser, testPassword)
		assert.Equal(t, ErrInvalidPassword, err)
	})

	t.Run("Given Non-Base64 Password Hash", func(t *testing.T) {
		testPassword := "myPassword1"
		testUser := &dal.User{
			PasswordHash: "werowerlew",
		}

		mockHasher := mock.NewMockHasher(ctrl)
		v := &userValidator{h: mockHasher}

		err := v.ValidatePassword(testUser, testPassword)
		assert.Equal(t, ErrInvalidPasswordHash, err)
	})
}
