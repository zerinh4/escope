package services

import (
	"errors"
	"fmt"
	"github.com/mertbahardogan/escope/internal/config"
	"github.com/mertbahardogan/escope/internal/connection"
	"github.com/mertbahardogan/escope/internal/constants"
	"os"
)

type ConfigService interface {
	SaveHost(alias string, config config.ConnectionConfig) error
	LoadHost(alias string) (config.ConnectionConfig, error)
	ListHosts() ([]string, error)
	DeleteHost(alias string) error
	ClearConfig() error
	ValidateConfig(config config.ConnectionConfig) error
	SetActiveHost(alias string) error
	GetActiveHost() (string, error)
	ClearActiveHost() error
	SetConnectionTimeout(timeout int) error
	GetConnectionTimeout() (int, error)
	GetAppConfig() (config.AppConfig, error)
}

type configService struct{}

func NewConfigService() ConfigService {
	return &configService{}
}

func (s *configService) SaveHost(alias string, cfg config.ConnectionConfig) error {
	if err := s.ValidateConfig(cfg); err != nil {
		return fmt.Errorf(constants.ErrConfigValidationFailed, err)
	}

	if err := config.SaveHost(alias, cfg); err != nil {
		return fmt.Errorf(constants.ErrFailedToSaveHostConfig, err)
	}

	return nil
}

func (s *configService) LoadHost(alias string) (config.ConnectionConfig, error) {
	savedConfig, err := config.LoadHost(alias)
	if err != nil {
		return config.ConnectionConfig{}, fmt.Errorf(constants.ErrFailedToLoadHostConfig, err)
	}

	return savedConfig, nil
}

func (s *configService) ListHosts() ([]string, error) {
	aliases, err := config.ListHosts()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrFailedToListHosts, err)
	}

	return aliases, nil
}

func (s *configService) DeleteHost(alias string) error {
	if err := config.DeleteHost(alias); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf(constants.ErrHostNotFound, alias)
		}
		return fmt.Errorf(constants.ErrFailedToDeleteHost, err)
	}

	return nil
}

func (s *configService) ClearConfig() error {
	cfgPath := os.ExpandEnv(constants.ConfigFileEnvPath)
	if err := os.Remove(cfgPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf(constants.ErrFailedToRemoveConfigFile, err)
	}

	connection.ClearConfig()
	return nil
}

func (s *configService) ValidateConfig(cfg config.ConnectionConfig) error {
	if cfg.Host == constants.EmptyString {
		return fmt.Errorf(constants.ErrHostIsRequired)
	}

	if cfg.Secure {
		if cfg.Username == constants.EmptyString {
			return fmt.Errorf(constants.ErrUsernameRequired)
		}
		if cfg.Password == constants.EmptyString {
			return fmt.Errorf(constants.ErrPasswordRequired)
		}
	}

	connConfig := connection.Config{
		Host:     cfg.Host,
		Username: cfg.Username,
		Password: cfg.Password,
		Secure:   cfg.Secure,
	}

	timeout := constants.DefaultTimeout
	if appTimeout, err := s.GetConnectionTimeout(); err == nil {
		timeout = appTimeout
	}

	if err := connection.TestConnection(connConfig, timeout); err != nil {
		return fmt.Errorf(constants.ErrConnectionFailed2, err)
	}

	return nil
}

func (s *configService) SetConnectionTimeout(timeout int) error {
	if err := config.SetConnectionTimeout(timeout); err != nil {
		return fmt.Errorf(constants.ErrFailedToSetTimeout, err)
	}
	return nil
}

func (s *configService) GetConnectionTimeout() (int, error) {
	timeout, err := config.GetConnectionTimeout()
	if err != nil {
		return 0, fmt.Errorf(constants.ErrFailedToGetTimeout, err)
	}
	return timeout, nil
}

func (s *configService) GetAppConfig() (config.AppConfig, error) {
	appCfg, err := config.GetAppConfig()
	if err != nil {
		return config.AppConfig{}, fmt.Errorf(constants.ErrFailedToGetAppConfig, err)
	}
	return appCfg, nil
}

func (s *configService) SetActiveHost(alias string) error {
	_, err := config.LoadHost(alias)
	if err != nil {
		if err == os.ErrNotExist {
			return fmt.Errorf(constants.ErrHostNotFound, alias)
		}
		return fmt.Errorf(constants.ErrFailedToLoadHost, err)
	}

	if err := config.SetActiveHost(alias); err != nil {
		return fmt.Errorf(constants.ErrFailedToSetActiveHost, err)
	}
	return nil
}

func (s *configService) GetActiveHost() (string, error) {
	activeHost, err := config.GetActiveHost()
	if err != nil {
		return constants.EmptyString, fmt.Errorf(constants.ErrFailedToGetActiveHost, err)
	}
	return activeHost, nil
}

func (s *configService) ClearActiveHost() error {
	if err := config.ClearActiveHost(); err != nil {
		return fmt.Errorf(constants.ErrFailedToClearActiveHost, err)
	}
	return nil
}
