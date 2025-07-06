package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	GroupMemberID uuid.UUID `json:"group_member_id" gorm:"primaryKey"`
	GroupID       uuid.UUID `json:"group_id"`
	UserID        uuid.UUID `json:"user_id" gorm:"foreignKey:UserID"`
	User          User      `json:"user" gorm:"references:UserID"`
	NotAllowed    bool      `json:"not_allowed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

// GroupMemberの作成
func CreateGroupMember(groupMember *GroupMember) error {
	if g, _ := FindGroupByGroupID(groupMember.GroupID); g == nil {
		return errors.New("Group not found")
	}
	if u, _ := FindUserByUserID(groupMember.UserID); u == nil {
		return errors.New("User not found")
	}

	r := db.Create(groupMember)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

// GroupMemberを探す関数群 //

// GroupMemberIDを用いてGroupMemberを探す
func FindGroupMemberByGroupMemberID(group_member_id uuid.UUID) (*GroupMember, error) {
	var groupMember GroupMember
	r := db.Where("group_member_id = ?", group_member_id).First(&groupMember)
	if r.Error != nil {
		return nil, r.Error
	}
	return &groupMember, nil
}

// GroupIDを用いてGroupMemberを探す
func FindGroupMemberByGroupID(group_id uuid.UUID) ([]GroupMember, error) {
	var groupMembers []GroupMember
	r := db.Where("group_id = ?", group_id).Find(&groupMembers)
	if r.Error != nil {
		return nil, r.Error
	}
	return groupMembers, nil
}

// UserIDを用いてGroupMemberを探す
func FindGroupMemberByUserID(user_id uuid.UUID) ([]GroupMember, error) {
	var groupMembers []GroupMember
	r := db.Where("user_id = ? AND not_allowed = ?", user_id, false).Find(&groupMembers)
	if r.Error != nil {
		return nil, r.Error
	}
	return groupMembers, nil
}
