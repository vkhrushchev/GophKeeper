package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vkhrushchev/GophKeeper/internal/server/entity"
	"github.com/vkhrushchev/GophKeeper/internal/server/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log/slog"
	"os"
	"testing"
)

type MongoUserRepositoryTestSuite struct {
	suite.Suite
	log        *slog.Logger
	container  *mongodb.MongoDBContainer
	repository *MongoUserRepository
}

func (s *MongoUserRepositoryTestSuite) SetupSuite() {
	s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	s.log.Info("setup test suite...")

	mongoDBContainer, err := mongodb.Run(
		context.TODO(),
		"mongo:8",
		mongodb.WithUsername("goph-keeper-test"),
		mongodb.WithPassword("goph-keeper-test"),
		testcontainers.WithEnv(map[string]string{
			"MONGO_INITDB_DATABASE": "GophKeeper",
		}),
	)
	if err != nil {
		s.log.Error("unexpected error when run mongodb container", "error", err)
		return
	}

	s.log.Info("mongodb container created...")

	s.container = mongoDBContainer
	connectionDSN, err := mongoDBContainer.ConnectionString(context.TODO())
	if err != nil {
		s.log.Error("unexpected error when get database connection string", "error", err)
	}

	s.log.Info("", "connectionDSN", connectionDSN)

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(connectionDSN))
	if err != nil {
		s.log.Error("unexpected error when create mongodb mongoClient", "error", err)
	}

	s.repository = &MongoUserRepository{
		log:         s.log,
		mongoClient: mongoClient,
	}
}

func (s *MongoUserRepositoryTestSuite) TearDownSuite() {
	s.log.Info("tear down test suite...")

	if err := s.repository.mongoClient.Disconnect(context.TODO()); err != nil {
		s.log.Error("unexpected error when disconnect mongodb client", "error", err)
	}

	if err := s.container.Terminate(context.TODO()); err != nil {
		s.log.Error("unexpected error when terminate mongodb container", "error", err)
	}
}

func (s *MongoUserRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	s.log.Debug("before test...", "suiteName", suiteName, "testName", testName)
}

func (s *MongoUserRepositoryTestSuite) AfterTest(suiteName, testName string) {
	s.log.Debug("after test...", "suiteName", suiteName, "testName", testName)

	err := s.repository.mongoClient.
		Database("GophKeeper").
		Collection("users").
		Drop(context.TODO())
	if err != nil {
		s.log.Error("unexpected error when drop mongodb collection", "collection", "users", "error", err)
	}
}

func (s *MongoUserRepositoryTestSuite) TestSave() {
	a := s.Assert()

	userEntity := entity.User{
		ID:           uuid.NewString(),
		Username:     "test_user",
		PasswordHash: "test_password_hash",
		PasswordSalt: "test_password_salt",
		Email:        "test_email@email.com",
	}

	savedUserEntity, err := s.repository.Save(context.Background(), userEntity)
	if !a.NoError(err) {
		s.FailNow("error not expected here: {}", err)
	}

	a.NotNil(savedUserEntity, "saved user entity is nil")
	a.Equal(userEntity.ID, savedUserEntity.ID, "saved entity id not equals generated")
}

func (s *MongoUserRepositoryTestSuite) TestSaveExisted() {
	a := s.Assert()

	userEntity := entity.User{
		ID:           uuid.NewString(),
		Username:     "test_user",
		PasswordHash: "test_password_hash",
		PasswordSalt: "test_password_salt",
		Email:        "test_email@email.com",
	}

	_, err := s.repository.Save(context.TODO(), userEntity)
	if !a.NoError(err) {
		s.FailNow("error not expected here: {}", err)
	}

	_, err = s.repository.Save(context.TODO(), userEntity)
	if !a.Error(err) {
		s.FailNow("expected error here", err)
	}

	a.Equal(usecase.ErrUserAlreadyExists, err, "expected ErrUserAlreadyExists here")
}

func (s *MongoUserRepositoryTestSuite) TestGet() {
	a := s.Assert()

	userEntity := entity.User{
		ID:           uuid.NewString(),
		Username:     "test_user",
		PasswordHash: "test_password_hash",
		PasswordSalt: "test_password_salt",
		Email:        "test_email@email.com",
	}

	savedUserEntity, err := s.repository.Save(context.Background(), userEntity)
	if !a.NoError(err) {
		s.Fail("error not expected here: {}", err)
	}

	userEntityByUsername, err := s.repository.Get(context.TODO(), savedUserEntity.Username)
	if !a.NoError(err) {
		s.Fail("error not expected here: {}", err)
	}

	a.NotNil(userEntityByUsername, "saved user entity is nil by username", "username", savedUserEntity.Username)
	a.Equal(savedUserEntity, userEntityByUsername, "founded user entity not equals saved")
}

func (s *MongoUserRepositoryTestSuite) TestGetNotFound() {
	a := s.Assert()

	_, err := s.repository.Get(context.TODO(), "not_existed_test")
	if !a.Error(err) {
		s.Fail("error expected here")
	}

	a.Equal(usecase.ErrUserNotFound, err, "expected ErrUserNotFound here")
}

func TestMongoUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MongoUserRepositoryTestSuite))
}
