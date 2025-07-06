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
	if err := model.AddGroupMemberByUserIDsWithAllowed(group.GroupID, req.UserIDs, true); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, group)
}

// GetGroup処理 //
// GetGroupリクエストに必要な型
type getGroupRequest struct {
	GroupID uuid.UUID `param:"groupID" binding:"required"`
}

// GetGroupのメイン処理
func GetGroup(c echo.Context) error {
	groupIDstr := c.Param("groupID")
	groupID, err := uuid.Parse(groupIDstr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid group ID")
	}

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	user.UserPassword = ""

	// Groupを取得
	group, err := model.FindGroupByGroupIDAndUserID(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group not found")
	}

	return c.JSON(http.StatusOK, group)
}

// GetGroupsByUser処理 //
// GetGroupsByUserで返す時に必要な型
type groupType struct {
	Group        model.Group         `json:"group"`
	GroupMembers []model.GroupMember `json:"group_members"`
}
type returnTyoeGetGroupsByUser struct {
	Groups []groupType `json:"groups"`
}

// GetGroupByUserのメイン処理
func GetGroupsByUser(c echo.Context) error {
	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	user.UserPassword = ""

	groups := make([]groupType, 0)

	// GroupMemberを取得
	groupMembers, err := model.FindGroupMemberByUserID(user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group members not found")
	}

	// Groupを取得
	for _, groupMember := range groupMembers {
		group, err := model.FindGroupByGroupID(groupMember.GroupID)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Group not found")
		}
		groups = append(groups, groupType{
			Group:        *group,
			GroupMembers: []model.GroupMember{groupMember},
		})
	}

	return c.JSON(http.StatusOK, groups)
}
