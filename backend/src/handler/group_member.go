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
	UserIDs string `form:"user_ids"`
	Admin   bool   `form:"admin"`
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

	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// Adminかどうか
	isAdmin, err := model.IsAdminOfGropMember(uuid.MustParse(groupIDstr), user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking admin status: "+err.Error())
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "You are not an admin of this group")
	}

	// groupMemberを追加
	if err := model.AddGroupMemberByUserIDsWithAllowed(uuid.MustParse(groupIDstr), req.UserIDs, req.Admin); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error adding group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, "Group members added successfully")
}

// GetGroupMembers //
// グループメンバー一覧を取得
func GetGroupMembers(c echo.Context) error {
	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// グループに所属しているか確認
	isMember, err := model.IsMemberOfGroup(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking group membership: "+err.Error())
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, "You are not a member of this group")
	}

	// グループメンバー一覧を取得
	members, err := model.FindGroupMemberByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error retrieving group members: "+err.Error())
	}

	return c.JSON(http.StatusOK, members)
}

// UpdateGroupMember //
// UpdateGroupMemberリクエストに必要な型
type UpdateGroupMemberRequest struct {
	Admin      bool `form:"admin"`
	NotAllowed bool `form:"not_allowed"`
}

// UpdateGroupMemberのメイン処理
func UpdateGroupMember(c echo.Context) error {
	req := new(UpdateGroupMemberRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}

	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// groupMemberIDを取得
	groupMemberIDstr := c.Param("memberID")
	if groupMemberIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group Member ID is required")
	}
	groupMemberID := uuid.MustParse(groupMemberIDstr)

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// Adminかどうか
	isAdmin, err := model.IsAdminOfGropMember(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking admin status: "+err.Error())
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "You are not an admin of this group")
	}

	// グループメンバーの存在確認
	groupMember, err := model.FindGroupMemberByGroupMemberID(groupMemberID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group member not found")
	}
	if groupMember.GroupID != groupID {
		return c.JSON(http.StatusForbidden, "Group member does not belong to this group")
	}

	// グループメンバーを更新
	updatedMember, err := model.UpdateGroupMember(groupMemberID, req.Admin, req.NotAllowed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error updating group member: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedMember)
}

// DeleteGroupMember //
func DeleteGroupMember(c echo.Context) error {
	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// groupMemberIDを取得
	groupMemberIDstr := c.Param("memberID")
	if groupMemberIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group Member ID is required")
	}
	groupMemberID := uuid.MustParse(groupMemberIDstr)

	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// Adminかどうか
	isAdmin, err := model.IsAdminOfGropMember(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking admin status: "+err.Error())
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "You are not an admin of this group")
	}

	// グループメンバーの存在確認
	groupMember, err := model.FindGroupMemberByGroupMemberID(groupMemberID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group member not found")
	}
	if groupMember.GroupID != groupID {
		return c.JSON(http.StatusForbidden, "Group member does not belong to this group")
	}

	// グループメンバーを削除
	if err := model.DeleteGroupMember(groupMemberID); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error deleting group member: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Group member deleted successfully"})
}
