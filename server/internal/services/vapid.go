package services

import (
	"os"

	"crosshatch/internal/database/models"
	"crosshatch/internal/dtos"

	webpush "github.com/SherClockHolmes/webpush-go"
	"gorm.io/gorm"
)

const (
	configVapidPublic  = "vapid_public_key"
	configVapidPrivate = "vapid_private_key"
)

type VapidService struct {
	publicKey  string
	privateKey string
	subject    string
}

func (s *VapidService) PublicKey() string  { return s.publicKey }
func (s *VapidService) PrivateKey() string { return s.privateKey }
func (s *VapidService) Subject() string    { return s.subject }

func NewVapidService(db *gorm.DB) (*VapidService, error) {
	subject := os.Getenv("VAPID_SUBJECT")
	if subject == "" {
		subject = "crosshatch@bwees.io"
	}

	public := os.Getenv("VAPID_PUBLIC_KEY")
	private := os.Getenv("VAPID_PRIVATE_KEY")

	if public == "" || private == "" {
		stored, err := loadVapidKeys(db)
		if err != nil {
			return nil, err
		}
		if stored != nil {
			public, private = stored[0], stored[1]
		}
	}

	if public == "" || private == "" {
		var err error
		private, public, err = webpush.GenerateVAPIDKeys()
		if err != nil {
			return nil, err
		}
		if err := persistVapidKeys(db, public, private); err != nil {
			return nil, err
		}
	}

	return &VapidService{publicKey: public, privateKey: private, subject: subject}, nil
}

func loadVapidKeys(db *gorm.DB) (*[2]string, error) {
	public, ok, err := getAppConfig(db, configVapidPublic)
	if err != nil || !ok {
		return nil, err
	}
	private, ok, err := getAppConfig(db, configVapidPrivate)
	if err != nil || !ok {
		return nil, err
	}
	return &[2]string{public, private}, nil
}

func persistVapidKeys(db *gorm.DB, public, private string) error {
	if err := db.Create(&models.AppConfig{Key: configVapidPublic, Value: public}).Error; err != nil {
		return err
	}
	return db.Create(&models.AppConfig{Key: configVapidPrivate, Value: private}).Error
}

func getAppConfig(db *gorm.DB, key string) (string, bool, error) {
	config := models.AppConfig{}
	err := db.First(&config, "key = ?", key).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", false, nil
		}
		return "", false, err
	}
	return config.Value, true, nil
}

func (s *VapidService) GetConfig() (dtos.VapidDto, error) {
	return dtos.VapidDto{PublicKey: s.publicKey}, nil
}
