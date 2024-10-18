package service

import (
	"fmt"

	"github.com/agusespa/a3n/internal/logger"
	"github.com/agusespa/a3n/internal/models"
	"github.com/agusespa/a3n/internal/repository"
)

type DefaultRealmService struct {
	AuthRepo repository.AuthRepository
	Logger   logger.Logger
}

type RealmService interface {
	GetRealmById(id int64) (models.RealmEntity, error)
}

func NewDefaultRealmService(authRepo *repository.MySqlRepository, logger logger.Logger) *DefaultRealmService {
	return &DefaultRealmService{
		AuthRepo: authRepo,
		Logger:   logger}
}

func (rs *DefaultRealmService) GetRealmById(realmID int64) (models.RealmEntity, error) {
	realm, err := rs.AuthRepo.ReadRealmById(realmID)
	if err != nil {
		rs.Logger.LogError(err)
		return models.RealmEntity{}, err
	}

	if err := validateRealm(realm); err != nil {
		rs.Logger.LogError(err)
		return models.RealmEntity{}, err
	}

	return realm, nil
}

func validateRealm(r models.RealmEntity) error {
	if r.RealmID == 0 {
		return fmt.Errorf("RealmID cannot be zero")
	}
	if r.RealmName == "" {
		return fmt.Errorf("RealmName cannot be empty")
	}
	if r.RealmDomain == "" {
		return fmt.Errorf("RealmDomain cannot be empty")
	}
	if r.RefreshExp <= 0 {
		return fmt.Errorf("RefreshExp must be greater than zero")
	}
	if r.AccessExp <= 0 {
		return fmt.Errorf("AccessExp must be greater than zero")
	}
	if r.EmailVerify && (r.EmailProvider == "" || r.EmailSender == "" || r.EmailAddr == "") {
		return fmt.Errorf("EmailProvider, EmailSender, and EmailAddr must be set when email verification is enabled")
	}
	return nil
}
