package service

import (
	"context"
	"errors"
	"log"

	"github.com/SohamChatterG/uptime/auth"
	"github.com/SohamChatterG/uptime/model"
	"github.com/SohamChatterG/uptime/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
	jwtSvc   *auth.JWTService
}

func NewUserService(repo *repository.UserRepository, jwtSvc *auth.JWTService) *UserService {
	return &UserService{userRepo: repo, jwtSvc: jwtSvc}
}

func (s *UserService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwtSvc.GenerateToken(user)
}

func (s *UserService) FindOrCreateUser(ctx context.Context, email, name string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)

	// If user does not exist, create a new one.
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User with email %s not found. Creating new user.", email)
			newUser := &model.User{
				Name:     name,
				Email:    email,
				Password: "", // No password for OAuth users
			}
			if err := s.userRepo.Create(ctx, newUser); err != nil {
				return "", err
			}
			user = newUser
		} else {
			// A different database error occurred
			return "", err
		}
	}

	// User exists or was just created, generate a JWT for them.
	return s.jwtSvc.GenerateToken(user)
}
