package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/vkhrushchev/GophKeeper/internal/server/domain"
	"github.com/vkhrushchev/GophKeeper/internal/server/entity"
	"log/slog"
)

//go:generate mockgen -source=./usecase.go -destination=./mock/mock.go

var (
	ErrUnexpected        = errors.New("unexpected error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongCredentials  = errors.New("wrong credentials")
)

type PasswordHashCalculator interface {
	CalculateHash(password string) (string, string, error)
}

type PasswordHashValidator interface {
	ValidateHash(expectedHash string, password string, salt string) (bool, error)
}

type UserSaver interface {
	Save(ctx context.Context, userEntity entity.User) (entity.User, error)
}

type UserProvider interface {
	Get(ctx context.Context, username string) (entity.User, error)
}

type TokenProvider interface {
	ProvideToken(ctx context.Context, user domain.User) (string, error)
}

type TokenInvalidator interface {
	InvalidateToken(ctx context.Context, user domain.User, token string) error
}

type RegisterUserUseCase struct {
	log                    *slog.Logger
	passwordHashCalculator PasswordHashCalculator
	userSaver              UserSaver
}

func (u *RegisterUserUseCase) RegisterUser(ctx context.Context, user domain.User) (domain.User, error) {
	u.log.Debug("usecase: RegisterUser", "user", user)

	passwordHash, passwordSalt, err := u.passwordHashCalculator.CalculateHash(user.Password)
	if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return domain.User{}, ErrUnexpected
	}

	userEntity := entity.User{
		ID:           uuid.NewString(),
		Username:     user.Username,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
		Email:        user.Email,
	}

	savedUserEntity, err := u.userSaver.Save(ctx, userEntity)
	if err != nil && errors.Is(err, ErrUserAlreadyExists) {
		u.log.Debug("usecase: user already exists", "username", user.Username, "email", user.Email)
		return domain.User{}, err
	} else if err != nil && errors.Is(err, ErrUnexpected) {
		u.log.Error("usecase: save user failed", "error", err)
		return domain.User{}, err
	} else if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return domain.User{}, ErrUnexpected
	}

	return domain.User{
		ID:       savedUserEntity.ID,
		Username: savedUserEntity.Username,
		Email:    savedUserEntity.Email,
	}, nil
}

type LoginUseCase struct {
	log                   *slog.Logger
	userProvider          UserProvider
	passwordHashValidator PasswordHashValidator
	tokenProvider         TokenProvider
}

func (u *LoginUseCase) Login(ctx context.Context, user domain.User) (domain.User, string, error) {
	u.log.Debug("usecase: Login", "user", user)

	userEntity, err := u.userProvider.Get(ctx, user.Username)
	if (err != nil) && errors.Is(err, ErrUserNotFound) {
		u.log.Debug("usecase: user not found", "username", user.Username)
		return domain.User{}, "", ErrUserNotFound
	} else if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return domain.User{}, "", ErrUnexpected
	}

	passwordHashValid, err := u.passwordHashValidator.ValidateHash(
		userEntity.PasswordHash,
		user.Password,
		userEntity.PasswordSalt,
	)
	if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return domain.User{}, "", ErrUnexpected
	}

	if !passwordHashValid {
		u.log.Debug("usecase: password hash validation failed", "user", user)
		return domain.User{}, "", ErrWrongCredentials
	}

	user = domain.User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Email:    userEntity.Email,
	}

	token, err := u.tokenProvider.ProvideToken(ctx, user)
	if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return domain.User{}, "", ErrUnexpected
	}

	return user, token, nil
}

type LogoutUseCase struct {
	log              *slog.Logger
	tokenInvalidator TokenInvalidator
}

func (u *LogoutUseCase) Logout(ctx context.Context, user domain.User, token string) error {
	u.log.Debug("usecase: Logout", "user", user)

	err := u.tokenInvalidator.InvalidateToken(ctx, user, token)
	if err != nil {
		u.log.Error("usecase: unexpected error", "error", err)
		return ErrUnexpected
	}

	return nil
}
