package service

import (
	"context"
	"errors"
	"time"

	"github.com/Zhandos28/ticket-booking/config"
	"github.com/Zhandos28/ticket-booking/model"
	"github.com/Zhandos28/ticket-booking/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	InvalidCredentialsError = errors.New("invalid credentials")
	UserAlreadyExistsError  = errors.New("user already exists")
)

type AuthService struct {
	repository model.AuthRepository
	envConfig  *config.EnvConfig
}

func (s *AuthService) Login(ctx context.Context, loginData *model.AuthCredentials) (string, *model.User, error) {
	user, err := s.repository.GetUser(ctx, "email = ?", loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, InvalidCredentialsError
		}
		return "", nil, nil
	}

	if ok := model.MatchesHash(loginData.Password, user.Password); !ok {
		return "", nil, InvalidCredentialsError
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 168).Unix(),
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, s.envConfig.JWTSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Register(ctx context.Context, registerData *model.AuthCredentials) (string, *model.User, error) {
	if _, err := s.repository.GetUser(ctx, "email = ?", registerData.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, UserAlreadyExistsError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	registerData.Password = string(hashedPassword)

	user, err := s.repository.RegisterUser(ctx, registerData)
	if err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 168).Unix(),
	}

	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, s.envConfig.JWTSecret)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func NewAuthService(repository model.AuthRepository, envConfig *config.EnvConfig) *AuthService {
	return &AuthService{repository: repository, envConfig: envConfig}
}
