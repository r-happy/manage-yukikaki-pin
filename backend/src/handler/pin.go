package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// AddPinByMember //
// AddPinByMemberリクエストに必要な型
type AddPinByMemberRequest struct {
	PinTypeID string  `form:"pin_type_id" binding:"required"`
	PinName   string  `form:"pin_name" binding:"required"`
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
}

// AddPinByMemberのメイン処理
func AddPinByMember(c echo.Context) error {
	req := new(AddPinByMemberRequest)
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

	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// PinTypeの存在確認
	pinTypeID := uuid.MustParse(req.PinTypeID)
	if _, err := model.FindPinTypeByPinTypeID(pinTypeID); err != nil {
		return c.JSON(http.StatusNotFound, "Pin Type not found")
	}

	// グループに所属してるか、PinTypeはそのグループに属しているか確認
	isMember, err := model.IsMemberOfGroup(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking group membership: "+err.Error())
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, "You are not a member of this group")
	}
	isPinType, err := model.IsPinTypeOfGroup(groupID, pinTypeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking pin type group: "+err.Error())
	}
	if !isPinType {
		return c.JSON(http.StatusForbidden, "Pin Type does not belong to this group")
	}

	pin := &model.Pin{
		PinID:              uuid.New(),
		GroupID:            groupID,
		PinTypeID:          pinTypeID,
		PinName:            req.PinName,
		PinCreatedByID:     user.UserID,
		PinTypeDescription: "", // 必要に応じて設定
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		NotAllowed:         false,
	}

	if err := model.CreatePin(pin); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error creating pin: "+err.Error())
	}

	return c.JSON(http.StatusOK, pin)
}

// GetPinsByGroupID //
func GetPinsByGroupID(c echo.Context) error {
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

	pins, err := model.FindPinsByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error retrieving pins: "+err.Error())
	}

	return c.JSON(http.StatusOK, pins)
}

// GetPublicPinsByGroupID //
// 認証不要でグループのピン情報を取得
func GetPublicPinsByGroupID(c echo.Context) error {
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// グループの存在確認
	group, err := model.FindGroupByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group not found")
	}
	if group == nil {
		return c.JSON(http.StatusNotFound, "Group not found")
	}

	pins, err := model.FindPinsByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error retrieving pins: "+err.Error())
	}

	return c.JSON(http.StatusOK, pins)
}

// UpdatePinByMember //
// UpdatePinByMemberリクエストに必要な型
type UpdatePinByMemberRequest struct {
	PinTypeID string  `form:"pin_type_id" binding:"required"`
	PinName   string  `form:"pin_name" binding:"required"`
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
}

// UpdatePinByMemberのメイン処理
func UpdatePinByMember(c echo.Context) error {
	req := new(UpdatePinByMemberRequest)
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

	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// pinIDを取得
	pinIDstr := c.Param("pinID")
	if pinIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Pin ID is required")
	}
	pinID := uuid.MustParse(pinIDstr)

	// Pinの存在確認とグループの確認
	pin, err := model.FindPinByPinID(pinID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Pin not found")
	}
	if pin.GroupID != groupID {
		return c.JSON(http.StatusForbidden, "Pin does not belong to this group")
	}

	// PinTypeの存在確認
	pinTypeID := uuid.MustParse(req.PinTypeID)
	if _, err := model.FindPinTypeByPinTypeID(pinTypeID); err != nil {
		return c.JSON(http.StatusNotFound, "Pin Type not found")
	}

	// グループに所属してるか、PinTypeはそのグループに属しているか確認
	isMember, err := model.IsMemberOfGroup(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking group membership: "+err.Error())
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, "You are not a member of this group")
	}
	isPinType, err := model.IsPinTypeOfGroup(groupID, pinTypeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking pin type group: "+err.Error())
	}
	if !isPinType {
		return c.JSON(http.StatusForbidden, "Pin Type does not belong to this group")
	}

	// ピンの更新
	updatedPin, err := model.UpdatePin(pinID, req.PinName, pinTypeID, req.Latitude, req.Longitude)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error updating pin: "+err.Error())
	}

	return c.JSON(http.StatusOK, updatedPin)
}

// DeletePinByMember //
func DeletePinByMember(c echo.Context) error {
	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// groupIDを取得
	groupIDstr := c.Param("groupID")
	if groupIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Group ID is required")
	}
	groupID := uuid.MustParse(groupIDstr)

	// pinIDを取得
	pinIDstr := c.Param("pinID")
	if pinIDstr == "" {
		return c.JSON(http.StatusBadRequest, "Pin ID is required")
	}
	pinID := uuid.MustParse(pinIDstr)

	// Pinの存在確認とグループの確認
	pin, err := model.FindPinByPinID(pinID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Pin not found")
	}
	if pin.GroupID != groupID {
		return c.JSON(http.StatusForbidden, "Pin does not belong to this group")
	}

	// グループに所属しているか確認
	isMember, err := model.IsMemberOfGroup(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking group membership: "+err.Error())
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, "You are not a member of this group")
	}

	// ピンの削除
	if err := model.DeletePin(pinID); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error deleting pin: "+err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Pin deleted successfully"})
}
