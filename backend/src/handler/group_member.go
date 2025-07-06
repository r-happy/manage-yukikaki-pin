package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// AddGroupMember //
// AddGroupMemberリクエストに必要な型
type AddGroupMemberRequest struct {
	GroupID    string `json:"group_id" binding:"required"`
	UserIDs    string `json:"user_id" binding:"required"`
	NotAllowed bool   `json:"not_allowed"`
}

func AddGroupMemberByUserIDsWithAllowed(groupID uuid.UUID, userIDs string) error {
	// UserIDsをカンマ区切りで分割し、各ユーザーIDを検証し追加
	for userID := range strings.SplitSeq(userIDs, ",") {
		groupMember := &model.GroupMember{
			GroupMemberID: uuid.New(),
			GroupID:       groupID,
			UserID:        uuid.MustParse(userID),
			NotAllowed:    false,
		}

		if err := model.CreateGroupMember(groupMember); err != nil {
			return errors.New("Failed to add group member: " + err.Error())
		}
	}
	return nil
}

// AddGroupMemberのメイン処理
func AddGroupMember(c echo.Context) error {
	req := new(AddGroupMemberRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// すべてのフィールドが埋まっている
	if err := ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// groupMemberを追加
	if err := AddGroupMemberByUserIDsWithAllowed(uuid.MustParse(req.GroupID), req.UserIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, "Group members added successfully")
}
