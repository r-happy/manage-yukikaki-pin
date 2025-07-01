package model

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	GroupID          uuid.UUID `json:"group_id" gorm:"primaryKey"`
	GroupName        string    `json:"group_name"`
	GroupDescription string    `json:"group_description"`
	GroupCreatedBy   User      `json:"group_created_by" gorm:"references:UserID"`
	CreatedAt        time.Time `json:"createad_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

// Groupの作成
func CreateGroup(group *Group) error {
	r := db.Create(group)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

// Groupを探す関数群 //

// GroupIDを用いてGroupを探す
func FindGroupByGroupID(group_id uuid.UUID) (*Group, error) {
	var group Group
	r := db.Where("group_id = ?", group_id).First(&group)
	if r.Error != nil {
		return nil, r.Error
	}
	return &group, nil
}

// グループを作った人（GroupCreatedBy）でGroupを探す
func FindGroupByGroupCreatedBy(user *User) ([]Group, error) {
	var groups []Group
	r := db.Where("group_created_by = ?", user).Find(&groups)
	if r.Error != nil {
		return nil, r.Error
	}
	return groups, nil
}
