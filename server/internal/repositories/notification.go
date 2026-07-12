package repositories

import (
	"crosshatch/internal/database/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NotificationRepository struct {
	db *gorm.DB
}

func (r *NotificationRepository) UpsertSubscription(deviceID, userID, endpoint, p256dh, auth string) error {
	sub := models.PushSubscription{
		DeviceID: deviceID,
		UserID:   userID,
		Endpoint: endpoint,
		P256dh:   p256dh,
		Auth:     auth,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "endpoint", "p256dh", "auth"}),
	}).Create(&sub).Error
}

func (r *NotificationRepository) DeleteSubscriptionByEndpoint(endpoint string) error {
	return r.db.Delete(&models.PushSubscription{}, "endpoint = ?", endpoint).Error
}

func (r *NotificationRepository) DeleteSubscriptionByDevice(deviceID string) error {
	return r.db.Delete(&models.PushSubscription{}, "device_id = ?", deviceID).Error
}

func (r *NotificationRepository) GetSubscriptionByDevice(deviceID string) (*models.PushSubscription, error) {
	sub := models.PushSubscription{}
	err := r.db.First(&sub, "device_id = ?", deviceID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *NotificationRepository) DeleteSubscriptionsForUser(userID string) error {
	return r.db.Delete(&models.PushSubscription{}, "user_id = ?", userID).Error
}

func (r *NotificationRepository) GetSetting(deviceID, serial string) (models.NotificationSetting, error) {
	setting := models.NotificationSetting{}
	err := r.db.First(&setting, "device_id = ? AND printer_serial = ?", deviceID, serial).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.NotificationSetting{DeviceID: deviceID, PrinterSerial: serial}, nil
		}
		return models.NotificationSetting{}, err
	}
	return setting, nil
}

func (r *NotificationRepository) UpsertSetting(setting models.NotificationSetting) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "printer_serial"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "enabled", "notify_complete", "notify_error"}),
	}).Create(&setting).Error
}

func (r *NotificationRepository) DeleteSettingsForUser(userID string) error {
	return r.db.Delete(&models.NotificationSetting{}, "user_id = ?", userID).Error
}

// SubscriptionsForEvent returns the push subscriptions of every device whose
// setting for the serial is enabled and opts into the given event
// ("complete" or "error").
func (r *NotificationRepository) SubscriptionsForEvent(serial, event string) ([]models.PushSubscription, error) {
	query := r.db.
		Joins("JOIN notification_setting ns ON ns.device_id = push_subscription.device_id").
		Where("ns.printer_serial = ? AND ns.enabled = ?", serial, true)

	switch event {
	case "complete":
		query = query.Where("ns.notify_complete = ?", true)
	case "error":
		query = query.Where("ns.notify_error = ?", true)
	default:
		return []models.PushSubscription{}, nil
	}

	subs := []models.PushSubscription{}
	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}
