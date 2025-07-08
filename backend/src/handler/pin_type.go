package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// AddPinType //
// AddPinTypeリクエストに必要な型
type AddPinTypeRequest struct {
	PinTypeName        string `form:"pin_type_name" binding:"required"`
	PinTypeDescription string `form:"pin_type_description" binding:"required"`
}

// AddPinTypeのメイン処理
func AddPinType(c echo.Context) error {
	req := new(AddPinTypeRequest)
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

	// Adminかどうか
	isAdmin, err := model.IsAdminOfGropMember(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking admin status: "+err.Error())
	}
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "You are not an admin of this group")
	}

	pinType := &model.PinType{
		PinTypeID:          uuid.New(),
		GroupID:            groupID,
		PinTypeName:        req.PinTypeName,
		PinTypeDescription: req.PinTypeDescription,
		PinTypeCreatedByID: user.UserID,
	}

	if err := model.CreatePinType(pinType); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error creating pin type: "+err.Error())
	}

	return c.JSON(http.StatusOK, pinType)
}

// GetPinTypesByGroupID //
func GetPinTypesByGroupID(c echo.Context) error {
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

	// Groupに所属しているか
	isMember, err := model.IsMemberOfGroup(groupID, user.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error checking group membership: "+err.Error())
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, "You are not a member of this group")
	}

	pinTypes, err := model.FindPinTypesByGroupID(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error fetching pin types: "+err.Error())
	}

	return c.JSON(http.StatusOK, pinTypes)
}
