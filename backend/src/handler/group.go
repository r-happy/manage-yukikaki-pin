package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// AddGroup処理 //
// AddGroupリクエストに必要な型
type addGroupRequest struct {
	GroupName        string `form:"group_name"`
	GroupDescription string `form:"group_description"`
	UserIDs          string `form:"user_ids"` // カンマ区切りのユーザーID
}

// AddGroupのメイン処理
func AddGroup(c echo.Context) error {
	req := new(addGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}
	// すべてのフィールドが埋まっているかどうか
	if err := ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	user.UserPassword = ""

	group := &model.Group{
		GroupID:          uuid.New(),
		GroupName:        req.GroupName,
		GroupDescription: req.GroupDescription,
		GroupCreatedByID: user.UserID,
	}

	if err := model.CreateGroup(group); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error")
	}

	// GroupMemberの作成
	if err := AddGroupMemberByUserIDsWithAllowed(group.GroupID, req.UserIDs); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, group)
}

// GetGroup処理 //
// GetGroupリクエストに必要な型
type getGroupRequest struct {
	GroupID uuid.UUID `form:"group_id"`
}

// GetGroupのメイン処理
func GetGroup(c echo.Context) error {
	req := new(getGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}

	// すべてのフィールドが埋まっているかどうか
	if err := ValidateStruct(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	user.UserPassword = ""

	// Groupを取得
	group, err := model.FindGroupByGroupID(req.GroupID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group not found")
	}

	return c.JSON(http.StatusOK, group)
}
