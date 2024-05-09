package user

import (
	"context"
	"errors"
	"time"

	"github.com/isdzulqor/donation-hub/internal/core/entity"
	"github.com/isdzulqor/donation-hub/internal/driver/rest"
)

type Service interface {
	RegisterUser(ctx context.Context, req rest.RegisterRequestBody) (err error)
	LoginUser(ctx context.Context, req rest.LoginRequestBody) (user entity.User, err error)
	GetListUser(ctx context.Context, limit int, page int, role string) (users []entity.User, err error)
}

type service struct {
	storage UserStorage
}

func NewService(storage UserStorage) Service {
	return &service{
		storage: storage,
	}
}

func (s *service) RegisterUser(ctx context.Context, req rest.RegisterRequestBody) (err error) {
	user := entity.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	// check if username, email and password are valid
	err = user.Validate()
	if err != nil {
		return
	}

	// check if user has already used
	exist, err := s.storage.IsExist(ctx, user.Email)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("user already exist")
		return
	}

	user.CreatedAt = time.Now()

	err = s.storage.RegisterNewUser(ctx, &user)

	return
}

func (s *service) LoginUser(ctx context.Context, req rest.LoginRequestBody) (user entity.User, err error) {
	exist, err := s.storage.IsExist(ctx, req.Email)
	if err != nil {
		return
	}
	if !exist {
		err = rest.ErrInvalidUsernameOrPassword
		return
	}

	user, err = s.storage.LoginUser(ctx, req)
	if err != nil {
		return
	}

	return user, nil
}

func (s *service) GetListUser(ctx context.Context, limit int, page int, role string) (users []entity.User, err error) {
	return s.storage.ListUser(ctx, limit, page, role)
}
