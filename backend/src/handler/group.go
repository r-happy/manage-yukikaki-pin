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
type pinWithType struct {
	Pin     model.Pin     `json:"pin"`
	PinType model.PinType `json:"pin_type"`
}

type groupType struct {
	Group        model.Group         `json:"group"`
	GroupMembers []model.GroupMember `json:"group_members"`
	Pins         []pinWithType       `json:"pins"` // ピンとそのタイプの情報
}

// GetGroupByUserのメイン処理
func GetGroupsByUser(c echo.Context) error {
	// user認証
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}
	user.UserPassword = ""

	// 自分が所属しているグループID一覧を取得
	groupMembers, err := model.FindAllGroupMemberByUserID(user.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Group members not found")
	}

	// グループIDの重複排除
	groupIDmap := make(map[uuid.UUID]bool)
	for _, gm := range groupMembers {
		groupIDmap[gm.GroupID] = true
	}

	groups := make([]groupType, 0, len(groupIDmap))
	for groupID := range groupIDmap {
		// グループ情報取得
		group, err := model.FindGroupByGroupID(groupID)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Group not found")
		}

		// 各グループの全メンバーを取得
		members, err := model.FindGroupMemberByGroupID(groupID)
		if err != nil {
			return c.JSON(http.StatusNotFound, "Group members not found")
		}

		// グループに属するすべてのPinを取得
		pins, err := model.FindPinsByGroupID(groupID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Error retrieving pins: "+err.Error())
		}

		// 各PinにそのPinTypeをマッピング
		pinsWithTypes := make([]pinWithType, 0, len(pins))
		for _, pin := range pins {
			// PinTypeを取得
			pinType, err := model.FindPinTypeByPinTypeID(pin.PinTypeID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, "Error retrieving pin type: "+err.Error())
			}

			pinsWithTypes = append(pinsWithTypes, pinWithType{
				Pin:     pin,
				PinType: *pinType,
			})
		}

		groups = append(groups, groupType{
			Group:        *group,
			GroupMembers: members,
			Pins:         pinsWithTypes,
		})
	}

	return c.JSON(http.StatusOK, groups)
}
