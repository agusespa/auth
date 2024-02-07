package service

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/agusespa/autz/internal/httperrors"
	"github.com/agusespa/autz/internal/models"
	"github.com/agusespa/autz/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	AuthRepo *repository.AuthRepository
}

func NewProductService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{AuthRepo: authRepo}
}

func (as *AuthService) RegisterNewUser(body models.AuthRequest) (int64, error) {
	uuidStr := uuid.New().String()

	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		err := httperrors.NewError(err, http.StatusBadRequest)
		return 0, err
	}

	id, err := as.AuthRepo.CreateUser(uuidStr, body.Email, hashedPassword)
	return id, err
}

func (as *AuthService) LoginUser(username, password string) (models.UserAuthData, error) {
	var userAuthData models.UserAuthData

	userData, err := as.AuthRepo.QueryUserByEmail(username)
	if err != nil {
		return userAuthData, err
	}

	if err := verifyPassword(userData.PasswordHash, password); err != nil {
		err := httperrors.NewError(err, http.StatusUnauthorized)
		return userAuthData, err
	}

	accessExpiresBy := time.Now().Add(5 * time.Minute).Unix()
	accessToken, err := generateJWT(userData.Email, userData.UserUUID, accessExpiresBy)
	if err != nil {
		return userAuthData, err
	}

	refreshExpiresBy := time.Now().AddDate(0, 12, 0).Unix()
	refreshToken, err := generateJWT(userData.Email, userData.UserUUID, refreshExpiresBy)
	if err != nil {
		return userAuthData, err
	}

	userAuthData = models.NewUserAuthData(userData.UserUUID, userData.Email, accessToken, refreshToken)
	return userAuthData, err
}

func (as *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := as.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessExpiresBy := time.Now().Add(5 * time.Minute).Unix()
	accessToken, err := generateJWT(claims.Username, claims.UserUUID, accessExpiresBy)
	if err != nil {
		return "", err
	}

	return accessToken, err
}

func (as *AuthService) ValidateToken(refreshToken string) (*models.CustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		refreshToken,
		&models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return privateKey, nil
		},
	)
	if err != nil {
		err := httperrors.NewError(err, http.StatusUnauthorized)
		return &models.CustomClaims{}, err
	}

	claims, ok := parsedToken.Claims.(*models.CustomClaims)
	if !ok {
		err := httperrors.NewError(errors.New("failed to parse token claims"), http.StatusUnauthorized)
		return &models.CustomClaims{}, err
	}

	return claims, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func verifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func generateJWT(username, useruuid string, expiration int64) (string, error) {
	claims := models.CustomClaims{
		Username: username,
		UserUUID: useruuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
			// TODO: use environment variable
			Issuer: "ecom",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		err := httperrors.NewError(err, http.StatusInternalServerError)
		return "", err
	}

	return tokenString, nil
}

var privateKey = []byte("AAAAB3NzaC1yc2EAAAADAQABAAABgQDVWRMI4ex2hUf0Lz/lPhbxIq3m28agw4XTOzYE2BwbHOlrAs23+rRyAW0jlaS1dCRz09fUGAqlxV13sQinS/VACXzvKzdCOxxGno2hGuIbxH6baXVmDRbFlK9qdeMtzXnppZ4cIVq33Y1IJYwZ1erj6QYqPhHcl4FmYuOL76/A6RptF3njBFqfU241lZuuDnbe2cFeihj0TFUOQVoH0Y/JK+Gwy0pebNy8hjnyGQZNBVeZw9R5UMxphtb2pbL1lKCoM7MDPLKGN+hhjRZyeLYEy/8AR1xiwE+R7LDaG/Zik5xQJ/YXYXMQBN2Ip4dTZdn40iuk+IWmaNT92Q5CpPvZO0aU5LWxPSLlZot4IloQZXr11ZKUXxzZvAh7OQXbolN/qTdBtKeOjw7iKvKiKGXTw6Uoq8fEUglPhX6ZcdGmELpHMx8VliXUjNPXbm9mSPk6Izx+HkcK2Zg5JLoqGNXf3wcOfbeJvEAAPafPlKFqoL/Okxgn/+fXuCh//z5Hrf0=")
