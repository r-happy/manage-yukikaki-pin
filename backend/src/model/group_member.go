package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GroupMember struct {
	GroupMemberID uuid.UUID `json:"group_member_id" gorm:"primaryKey"`
	GroupID       uuid.UUID `json:"group_id"`
	UserID        uuid.UUID `json:"user_id" gorm:"foreignKey:UserID"`
	User          User      `json:"user" gorm:"references:UserID"`
	NotAllowed    bool      `json:"not_allowed"`
	Admin         bool      `json:"admin"`
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
	// すでに同じグループに同じユーザーが存在するか確認
	existingMember, _ := FindGroupMemberByGroupIDAndUserID(groupMember.GroupID, groupMember.UserID)
	if existingMember != nil {
		return errors.New("User is already a member of this group")
	}

	r := db.Create(groupMember)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

// AddGroupMemberByUserIDsWithAllowedは、グループIDとユーザーIDのカンマ区切り文字列を受け取り、
// 指定されたグループにユーザーを追加します。
func AddGroupMemberByUserIDsWithAllowed(groupID uuid.UUID, userIDs string, admin bool) error {
	// UserIDsをカンマ区切りで分割し、各ユーザーIDを検証し追加
	for userID := range strings.SplitSeq(userIDs, ",") {
		groupMember := GroupMember{
			GroupMemberID: uuid.New(),
			GroupID:       groupID,
			UserID:        uuid.MustParse(userID),
			NotAllowed:    false,
			Admin:         admin,
		}

		if err := CreateGroupMember(&groupMember); err != nil {
			return errors.New("Failed to add group member: " + err.Error())
		}
	}
	return nil
}

// GroupMemberを探す関数群 //

// GroupMemberIDを用いてGroupMemberを探す
func FindGroupMemberByGroupMemberID(group_member_id uuid.UUID) (*GroupMember, error) {
	var groupMember GroupMember
	r := db.Preload("User").Where("group_member_id = ?", group_member_id).First(&groupMember)
	if r.Error != nil {
		return nil, r.Error
	}
	return &groupMember, nil
}

// GroupIDを用いてGroupMemberを探す
func FindGroupMemberByGroupID(group_id uuid.UUID) ([]GroupMember, error) {
	var groupMembers []GroupMember
	r := db.Preload("User").Where("group_id = ?", group_id).Find(&groupMembers)
	if r.Error != nil {
		return nil, r.Error
	}
	return groupMembers, nil
}

// UserIDを用いてGroupMemberを探す
func FindGroupMemberByUserID(user_id uuid.UUID) ([]GroupMember, error) {
	var groupMembers []GroupMember
	r := db.Preload("User").Where("user_id = ? AND not_allowed = ?", user_id, false).Find(&groupMembers)
	if r.Error != nil {
		return nil, r.Error
	}
	return groupMembers, nil
}

// GroupIDとUserIDを用いてGroupMemberを探す
func FindGroupMemberByGroupIDAndUserID(group_id uuid.UUID, user_id uuid.UUID) (*GroupMember, error) {
	var groupMember GroupMember
	r := db.Preload("User").Where("group_id = ? AND user_id = ?", group_id, user_id).First(&groupMember)
	if r.Error != nil {
		return nil, r.Error
	}
	return &groupMember, nil
}

// 自分が所属しているグループのすべてのGroupMemberを探す
func FindAllGroupMemberByUserID(user_id uuid.UUID) ([]GroupMember, error) {
	var groupIDs []uuid.UUID
	if err := db.Model(&GroupMember{}).
		Where("user_id = ? AND not_allowed = ?", user_id, false).
		Pluck("group_id", &groupIDs).Error; err != nil {
		return nil, err
	}
	if len(groupIDs) == 0 {
		return []GroupMember{}, nil
	}

	var allGroupMembers []GroupMember
	if err := db.Preload("User").
		Where("group_id IN ?", groupIDs).
		Find(&allGroupMembers).Error; err != nil {
		return nil, err
	}
	return allGroupMembers, nil
}

// UserID + GroupIDでAdminかどうかを確認する関数
func IsAdminOfGropMember(groupID uuid.UUID, userID uuid.UUID) (bool, error) {
	var groupMember GroupMember
	r := db.Where("group_id = ? AND user_id = ? AND admin = ?", groupID, userID, true).First(&groupMember)
	if r.Error != nil {
		if r.Error.Error() == "record not found" {
			return false, nil
		}
		return false, r.Error
	}
	return true, nil
}

// Userがグループに所属しているかどうかを確認する関数
func IsMemberOfGroup(groupID uuid.UUID, userID uuid.UUID) (bool, error) {
	var groupMember GroupMember
	r := db.Where("group_id = ? AND user_id = ? AND not_allowed = ?", groupID, userID, false).First(&groupMember)
	if r.Error != nil {
		if r.Error.Error() == "record not found" {
			return false, nil
		}
		return false, r.Error
	}
	return true, nil
}
