package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vkhrushchev/GophKeeper/internal/server/domain"
	"github.com/vkhrushchev/GophKeeper/internal/server/entity"
	mock_usecase "github.com/vkhrushchev/GophKeeper/internal/server/usecase/mock"
	"log/slog"
	"os"
	"testing"
)

type RegisterUserUseCaseTestSuite struct {
	suite.Suite
	log                        *slog.Logger
	passwordHashCalculatorMock *mock_usecase.MockPasswordHashCalculator
	userSaverMock              *mock_usecase.MockUserSaver
	registerUserUseCase        *RegisterUserUseCase
}

func (s *RegisterUserUseCaseTestSuite) SetupSuite() {
	s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func (s *RegisterUserUseCaseTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(s.T())

	s.passwordHashCalculatorMock = mock_usecase.NewMockPasswordHashCalculator(mockCtrl)
	s.userSaverMock = mock_usecase.NewMockUserSaver(mockCtrl)

	s.registerUserUseCase = &RegisterUserUseCase{
		log:                    s.log,
		passwordHashCalculator: s.passwordHashCalculatorMock,
		userSaver:              s.userSaverMock,
	}
}

func (s *RegisterUserUseCaseTestSuite) TestRegisterUser() {
	a := s.Assert()

	s.passwordHashCalculatorMock.EXPECT().
		CalculateHash(gomock.Any()).
		Return("test_password_hash", "test_password_salt", nil)

	savedUserEntity := entity.User{
		ID:       uuid.New().String(),
		Username: "test_username",
		Email:    "test_email@example.com",
	}
	s.userSaverMock.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(savedUserEntity, nil)

	userDomain, err := s.registerUserUseCase.RegisterUser(context.Background(), domain.User{
		Username: "test_user_name",
		Password: "test_password",
		Email:    "test_email@example.com",
	})

	if !a.NoError(err) {
		s.Fail("not expected error here: {}", err)
	}

	expectedUserDomain := domain.User{
		ID:       savedUserEntity.ID,
		Username: savedUserEntity.Username,
		Email:    savedUserEntity.Email,
	}

	a.Equal(expectedUserDomain, userDomain, "user domain should be equal expected")
}

func (s *RegisterUserUseCaseTestSuite) TestRegisterUser_CalculateHashError() {
	a := s.Assert()

	s.passwordHashCalculatorMock.EXPECT().
		CalculateHash(gomock.Any()).
		Return("", "", errors.New("calculate hash error"))

	_, err := s.registerUserUseCase.RegisterUser(context.Background(), domain.User{
		Username: "test_user_name",
		Password: "test_password",
		Email:    "test_email@example.com",
	})

	if !a.Error(err) {
		s.Fail("expected error here: {}", err)
	}

	a.Equal(ErrUnexpected, err, "err mast be equal ErrUnexpected")
}

func (s *RegisterUserUseCaseTestSuite) TestRegisterUser_ErrUserAlreadyExists() {
	a := s.Assert()

	s.passwordHashCalculatorMock.EXPECT().
		CalculateHash(gomock.Any()).
		Return("test_password_hash", "test_password_salt", nil)

	s.userSaverMock.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(entity.User{}, ErrUserAlreadyExists)

	_, err := s.registerUserUseCase.RegisterUser(context.Background(), domain.User{
		Username: "test_user_name",
		Password: "test_password",
		Email:    "test_email@example.com",
	})

	if !a.Error(err) {
		s.Fail("expected error here: {}", err)
	}

	a.Equal(ErrUserAlreadyExists, err, "err mast be equal ErrUserAlreadyExists")
}

func (s *RegisterUserUseCaseTestSuite) TestRegisterUser_ErrUnexpected() {
	a := s.Assert()

	s.passwordHashCalculatorMock.EXPECT().
		CalculateHash(gomock.Any()).
		Return("test_password_hash", "test_password_salt", nil)

	s.userSaverMock.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(entity.User{}, ErrUnexpected)

	_, err := s.registerUserUseCase.RegisterUser(context.Background(), domain.User{
		Username: "test_user_name",
		Password: "test_password",
		Email:    "test_email@example.com",
	})

	if !a.Error(err) {
		s.Fail("expected error here: {}", err)
	}

	a.Equal(ErrUnexpected, err, "err mast be equal ErrUnexpected")
}

func (s *RegisterUserUseCaseTestSuite) TestRegisterUser_NotExpectedErr() {
	a := s.Assert()

	s.passwordHashCalculatorMock.EXPECT().
		CalculateHash(gomock.Any()).
		Return("test_password_hash", "test_password_salt", nil)

	s.userSaverMock.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(entity.User{}, errors.New("not expected error"))

	_, err := s.registerUserUseCase.RegisterUser(context.Background(), domain.User{
		Username: "test_user_name",
		Password: "test_password",
		Email:    "test_email@example.com",
	})

	if !a.Error(err) {
		s.Fail("expected error here: {}", err)
	}

	a.Equal(ErrUnexpected, err, "err mast be equal ErrUnexpected")
}

func TestRegisterUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterUserUseCaseTestSuite))
}

type LoginUseCaseTestSuite struct {
	suite.Suite
	log                       *slog.Logger
	userProviderMock          *mock_usecase.MockUserProvider
	passwordHashValidatorMock *mock_usecase.MockPasswordHashValidator
	tokenProviderMock         *mock_usecase.MockTokenProvider
	loginUseCase              *LoginUseCase
}

func (s *LoginUseCaseTestSuite) SetupSuite() {
	s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func (s *LoginUseCaseTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(s.T())

	s.userProviderMock = mock_usecase.NewMockUserProvider(mockCtrl)
	s.passwordHashValidatorMock = mock_usecase.NewMockPasswordHashValidator(mockCtrl)
	s.tokenProviderMock = mock_usecase.NewMockTokenProvider(mockCtrl)

	s.loginUseCase = &LoginUseCase{
		log:                   s.log,
		userProvider:          s.userProviderMock,
		passwordHashValidator: s.passwordHashValidatorMock,
		tokenProvider:         s.tokenProviderMock,
	}
}

func (s *LoginUseCaseTestSuite) TestLogin() {
	a := s.Assert()

	userEntity := entity.User{
		ID:       uuid.NewString(),
		Username: "test_user_name",
		Email:    "test_email@example.com",
	}

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(userEntity, nil)

	s.passwordHashValidatorMock.EXPECT().
		ValidateHash(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil)

	s.tokenProviderMock.EXPECT().
		ProvideToken(gomock.Any(), gomock.Any()).
		Return("test_token", nil)

	userDomain, token, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.NoError(err) {
		s.Fail("error not expected here: {}", err)
	}

	expectedUserDomain := domain.User{
		ID:       userEntity.ID,
		Username: userEntity.Username,
		Email:    userEntity.Email,
	}

	a.Equal(expectedUserDomain, userDomain, "user domain should be equal expected")
	a.NotEmpty(token, "token should not be empty")
}

func (s *LoginUseCaseTestSuite) TestLogin_ErrUserNotFound() {
	a := s.Assert()

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(entity.User{}, ErrUserNotFound)

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUserNotFound, err, "err must be equal ErrUserNotFound")
}

func (s *LoginUseCaseTestSuite) TestLogin_UserProvider_ErrUserNotFound() {
	a := s.Assert()

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(entity.User{}, ErrUserNotFound)

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUserNotFound, err, "err must be equal ErrUserNotFound")
}

func (s *LoginUseCaseTestSuite) TestLogin_UserProvider_ErrUnexpected() {
	a := s.Assert()

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(entity.User{}, ErrUnexpected)

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUnexpected, err, "err must be equal ErrUnexpected")
}

func (s *LoginUseCaseTestSuite) TestLogin_UserProvider_NotExpectedError() {
	a := s.Assert()

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(entity.User{}, errors.New("not expected error"))

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUnexpected, err, "err must be equal ErrUnexpected")
}

func (s *LoginUseCaseTestSuite) TestLogin_PasswordHashValidator_NotValid() {
	a := s.Assert()

	userEntity := entity.User{
		ID:       uuid.NewString(),
		Username: "test_user_name",
		Email:    "test_email@example.com",
	}

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(userEntity, nil)

	s.passwordHashValidatorMock.EXPECT().
		ValidateHash(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, nil)

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrWrongCredentials, err, "err must be equal ErrWrongCredentials")
}

func (s *LoginUseCaseTestSuite) TestLogin_PasswordHashValidator_Error() {
	a := s.Assert()

	userEntity := entity.User{
		ID:       uuid.NewString(),
		Username: "test_user_name",
		Email:    "test_email@example.com",
	}

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(userEntity, nil)

	s.passwordHashValidatorMock.EXPECT().
		ValidateHash(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, errors.New("test error"))

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUnexpected, err, "err must be equal ErrUnexpected")
}

func (s *LoginUseCaseTestSuite) TestLogin_TokenProvider_Error() {
	a := s.Assert()

	userEntity := entity.User{
		ID:       uuid.NewString(),
		Username: "test_user_name",
		Email:    "test_email@example.com",
	}

	s.userProviderMock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(userEntity, nil)

	s.passwordHashValidatorMock.EXPECT().
		ValidateHash(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil)

	s.tokenProviderMock.EXPECT().
		ProvideToken(gomock.Any(), gomock.Any()).
		Return("", errors.New("test error"))

	_, _, err := s.loginUseCase.Login(context.Background(), domain.User{})
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(ErrUnexpected, err, "err must be equal ErrUnexpected")
}

func TestLoginUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUseCaseTestSuite))
}

type LogoutUseCaseTestSuite struct {
	suite.Suite
	log                  *slog.Logger
	tokenInvalidatorMock *mock_usecase.MockTokenInvalidator
	logoutUseCase        *LogoutUseCase
}

func (s *LogoutUseCaseTestSuite) SetupSuite() {
	s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func (s *LogoutUseCaseTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(s.T())

	s.tokenInvalidatorMock = mock_usecase.NewMockTokenInvalidator(mockCtrl)

	s.logoutUseCase = &LogoutUseCase{
		log:              s.log,
		tokenInvalidator: s.tokenInvalidatorMock,
	}
}

func (s *LogoutUseCaseTestSuite) TestLogout() {
	a := s.Assert()

	s.tokenInvalidatorMock.EXPECT().
		InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	err := s.logoutUseCase.Logout(context.Background(), domain.User{}, "test_token")

	a.Nil(err, "err must be nil")
}

func (s *LogoutUseCaseTestSuite) TestLogout_TokenInvalidator_Error() {
	a := s.Assert()

	s.tokenInvalidatorMock.EXPECT().
		InvalidateToken(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(errors.New("test error"))

	err := s.logoutUseCase.Logout(context.Background(), domain.User{}, "test_token")

	a.Equal(ErrUnexpected, err, "err must be equal ErrUnexpected")
}

func TestLogoutUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutUseCaseTestSuite))
}
