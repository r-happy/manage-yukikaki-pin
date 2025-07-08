package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PinType struct {
	PinTypeID          uuid.UUID      `json:"pin_type_id" gorm:"primaryKey"`
	GroupID            uuid.UUID      `json:"group_id"`
	Group              Group          `json:"group" gorm:"foreignKey:GroupID;references:GroupID"`
	PinTypeName        string         `json:"pin_type_name" gorm:"not null"`
	PinTypeDescription string         `json:"pin_type_description" gorm:"not null"`
	PinTypeCreatedByID uuid.UUID      `json:"pin_type_created_by_id"`
	PinTypeCreatedBy   User           `json:"pin_type_created_by" gorm:"foreignKey:PinTypeCreatedByID;references:UserID"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// PinTypeの作成
func CreatePinType(pinType *PinType) error {
	if g, _ := FindGroupByGroupID(pinType.GroupID); g == nil {
		return errors.New("Group not found")
	}
	if u, _ := FindUserByUserID(pinType.PinTypeCreatedByID); u == nil {
		return errors.New("User not found")
	}

	r := db.Create(pinType)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

// PinTypeを探す関数群 //

// PinTypeIDで探す
func FindPinTypeByPinTypeID(pinTypeID uuid.UUID) (*PinType, error) {
	var pinType PinType
	r := db.Preload("Group").Preload("PinTypeCreatedBy").Where("pin_type_id = ?", pinTypeID).First(&pinType)

	if r.Error != nil {
		return nil, r.Error
	}

	return &pinType, nil
}

// GroupIDで探す
func FindPinTypesByGroupID(groupID uuid.UUID) ([]PinType, error) {
	var pinTypes []PinType
	r := db.Preload("Group").Preload("PinTypeCreatedBy").Where("group_id = ?", groupID).Find(&pinTypes)

	if r.Error != nil {
		return nil, r.Error
	}

	return pinTypes, nil
}

// グループに所属しているかどうかを確認
func IsPinTypeOfGroup(groupID uuid.UUID, pinTypeID uuid.UUID) (bool, error) {
	var count int64
	r := db.Model(&PinType{}).Where("group_id = ? AND pin_type_id = ?", groupID, pinTypeID).Count(&count)

	if r.Error != nil {
		return false, r.Error
	}

	return count > 0, nil
}
