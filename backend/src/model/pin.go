package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pin struct {
	PinID              uuid.UUID      `json:"pin_id" gorm:"primaryKey"`
	GroupID            uuid.UUID      `json:"group_id"`
	Group              Group          `json:"group" gorm:"foreignKey:GroupID;references:GroupID"`
	PinTypeID          uuid.UUID      `json:"pin_type_id"`
	PinType            PinType        `json:"pin_type" gorm:"foreignKey:PinTypeID;references:PinTypeID"`
	PinName            string         `json:"pin_name" gorm:"not null"`
	PinTypeDescription string         `json:"pin_type_description" gorm:"not null"` // この名前が適切か確認
	PinCreatedByID     uuid.UUID      `json:"pin_created_by_id"`
	PinCreatedBy       User           `json:"pin_created_by" gorm:"foreignKey:PinCreatedByID;references:UserID"`
	Latitude           float64        `json:"latitude" gorm:"not null"`
	Longitude          float64        `json:"longitude" gorm:"not null"`
	NotAllowed         bool           `json:"not_allowed"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Pinの作成
func CreatePin(pin *Pin) error {
	if pin.Latitude < -90 || pin.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90 degrees")
	}
	if pin.Longitude < -180 || pin.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180 degrees")
	}

	if g, _ := FindGroupByGroupID(pin.GroupID); g == nil {
		return errors.New("group not found")
	}
	if pt, _ := FindPinTypeByPinTypeID(pin.PinTypeID); pt == nil {
		return errors.New("pin type not found")
	}
	if u, _ := FindUserByUserID(pin.PinCreatedByID); u == nil {
		return errors.New("user not found")
	}

	r := db.Create(pin)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

// Pinを探す関数群 //

// PinIDで探す
func FindPinByPinID(pinID uuid.UUID) (*Pin, error) {
	var pin Pin
	r := db.Preload("Group").Preload("PinType").Preload("PinCreatedBy").Where("pin_id = ?", pinID).First(&pin)

	if r.Error != nil {
		return nil, r.Error
	}

	return &pin, nil
}

// GroupIDで探す
// GroupIDで探す（深いPreload）
func FindPinsByGroupID(groupID uuid.UUID) ([]Pin, error) {
	var pins []Pin
	r := db.Preload("Group").
		Preload("Group.GroupCreatedBy").
		Preload("PinType").
		Preload("PinType.Group").
		Preload("PinType.PinTypeCreatedBy").
		Preload("PinCreatedBy").
		Where("group_id = ?", groupID).
		Find(&pins)

	if r.Error != nil {
		return nil, r.Error
	}

	return pins, nil
}
