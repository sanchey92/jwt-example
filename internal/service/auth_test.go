package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	appError "github.com/sanchey92/jwt-example/internal/errors"
	"github.com/sanchey92/jwt-example/internal/models"
	"github.com/sanchey92/jwt-example/internal/service/mocks"
)

const (
	testEmail    = "test@example.com"
	testPassword = "password123"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		password     string
		mockUserRepo func(m *mocks.MockUserRepository)
		wantErr      bool
	}{
		{
			name:     "success registration",
			email:    testEmail,
			password: testPassword,
			mockUserRepo: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, user *models.User) error {
						assert.Equal(t, testEmail, user.Email)
						assert.NotEmpty(t, user.ID)
						assert.NotEmpty(t, user.Password)
						assert.Equal(t, models.RoleUser, user.Role)
						assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
						assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
						return nil
					})
			},
			wantErr: false,
		},
		{
			name:     "failure registration",
			email:    testEmail,
			password: testPassword,
			mockUserRepo: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(appError.ErrUserAlreadyExists)
			},
			wantErr: true,
		},
		{
			name:     "password hashing failure",
			email:    testEmail,
			password: string(make([]byte, 10000)),
			mockUserRepo: func(m *mocks.MockUserRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepository(ctrl)

			s := newTestAuthService(mockRepo)

			tt.mockUserRepo(mockRepo)

			user, err := s.Register(context.Background(), tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)

				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))

				assert.NoError(t, err)
			}
		})
	}
}

func newTestAuthService(repo UserRepository) *AuthService {
	return &AuthService{
		userRepo: repo,
		log:      zap.NewNop(),
	}
}
