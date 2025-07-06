package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// AddGroupMember //
// AddGroupMemberリクエストに必要な型
type AddGroupMemberRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserIDs string `json:"user_id" binding:"required"`
	Admin   bool   `json:"admin" binding:"required"`
}

// AddGroupMemberのメイン処理
func AddGroupMemberByAdmin(c echo.Context) error {
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

	// Adminかどうか
	isAdmin, err := model.IsAdminOfGropMember(uuid.MustParse(req.GroupID), user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking admin status: "+err.Error())
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "You are not an admin of this group")
	}

	// groupMemberを追加
	if err := model.AddGroupMemberByUserIDsWithAllowed(uuid.MustParse(req.GroupID), req.UserIDs, req.Admin); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, "Group members added successfully")
}
