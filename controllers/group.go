package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/wgplaner/wg_planer_server/models"
	"github.com/wgplaner/wg_planer_server/modules/base"
	"github.com/wgplaner/wg_planer_server/modules/setting"
	"github.com/wgplaner/wg_planer_server/restapi/operations/group"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

var groupLog = logging.MustGetLogger("Group")

func GetGroup(params group.GetGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`Get group "%s"`, params.GroupID)

	var g *models.Group
	var err error

	// Database
	if g, err = models.GetGroupByUID(params.GroupID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, g.UID)
		return NewNotFoundResponse("Group not found on server")

	} else if err != nil {
		groupLog.Critical(`Database Error!`, err)
		return NewInternalServerError("Internal Database Error")
	}

	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("User is a member of the specified group")
	}

	return group.NewGetGroupOK().WithPayload(g)
}

func GetGroupImage(params group.GetGroupImageParams, principal *models.User) middleware.Responder {
	var g *models.Group
	var err error

	if g, err = models.GetGroupByUID(params.GroupID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, g.UID)
		return NewNotFoundResponse("Group not found on server")

	} else if err != nil {
		groupLog.Critical(`Database Error getting group!`, err)
		return NewInternalServerError("Internal Database Error")
	}

	//if !base.StringInSlice(*principal.UID, g.Members) {
	//	return NewUnauthorizedResponse("User not a member of the group.")
	//}

	var imgFile *os.File
	var fileErr error

	// Get default image if normal one does no exist
	if imgFile, fileErr = models.GetGroupImage(g.UID); os.IsNotExist(fileErr) {
		imgFile, fileErr = models.GetGroupImageDefault()
	}

	if fileErr != nil {
		groupLog.Error("Error getting group's profile image ", fileErr.Error())
		return NewInternalServerError("Internal Server Error with profile image")
	}

	return group.NewGetGroupImageOK().WithPayload(imgFile)

}

func CreateGroupCode(params group.CreateGroupCodeParams, principal *models.User) middleware.Responder {
	groupLog.Debugf(`Generate group code for group "%s"!`, params.GroupID)

	var (
		c   *models.GroupCode
		err error

		groupUid = strfmt.UUID(params.GroupID)
	)

	if principal.GroupUID != groupUid {
		return NewUnauthorizedResponse("Can't create group code for other groups")
	}

	// Group MUST exist or we have inconsistencies
	if _, err = models.GetGroupByUID(groupUid); err != nil {
		groupLog.Debugf(`Error validating group "%s": "%s"`, params.GroupID, err.Error())
		return NewInternalServerError("Internal Server Error")
	}

	// TODO: Check authorization for user in the group

	if c, err = models.CreateGroupCode(groupUid); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	return group.NewCreateGroupCodeOK().WithPayload(c)
}

func CreateGroup(params group.CreateGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debug(`Start creating group`)

	// Create new group
	newGroupUid := strfmt.UUID(uuid.NewV4().String())

	theUser := &models.User{
		UID:      principal.UID,
		GroupUID: newGroupUid,
	}

	theGroup := models.Group{
		UID:         newGroupUid,
		DisplayName: params.Body.DisplayName,
		Currency:    params.Body.Currency,
	}

	// TODO: Check if user has already a group

	if err := models.UpdateUserCols(theUser, "group_uid"); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	// Insert new user into database
	if err := models.CreateGroup(&theGroup); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	groupLog.Infof(`Created group "%s"`, theGroup.UID)

	return group.NewCreateGroupOK().WithPayload(&theGroup)
}

func UpdateGroup(params group.UpdateGroupParams, principal *models.User) middleware.Responder {
	groupLog.Debug(`Start updating group`)

	var g *models.Group
	var err error

	if g, err = models.GetGroupByUID(params.Body.UID); models.IsErrGroupNotExist(err) {
		groupLog.Debugf(`Update group: "%s" does not exist: %s"`, g.UID, err)
		return NewNotFoundResponse(`Group does not exist.`)

	} else if err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	g.DisplayName = params.Body.DisplayName
	g.Currency = params.Body.Currency

	// TODO: Check Auth (is admin)

	// Update user into database
	if err := models.UpdateGroupCols(g, `display_name`, `currency`); err != nil {
		groupLog.Critical("Database error!", err)
		return NewInternalServerError("Internal Database Error")
	}

	groupLog.Infof(`Updated group "%s"`, g.UID)

	return group.NewCreateGroupOK().WithPayload(g)
}

func JoinGroup(params group.JoinGroupParams, principal *models.User) middleware.Responder {
	theGroup, err := principal.JoinGroupWithCode(params.GroupCode)

	if err == nil {
		return group.NewJoinGroupOK().WithPayload(theGroup)
	}

	// TODO: Handle different errors
	groupLog.Error(`Unknown Internal Server Error: `, err)
	return NewBadRequest("Unknown Server Error")
}

func JoinGroupHelp(params group.JoinGroupHelpParams) middleware.Responder {
	groupLog.Debug(`Get help site for joining group`)

	var (
		filepath = path.Join(setting.AppWorkPath, "views/group_code.html")
		templ    = template.Must(template.ParseFiles(filepath))
		buf      = bytes.NewBuffer([]byte{})
		content  = map[string]string{"GroupCode": params.GroupCode}
	)

	r := regexp.MustCompile(`^[A-Z0-9]{12}$`)
	if !r.MatchString(params.GroupCode) {
		return group.NewJoinGroupHelpDefault(http.StatusBadRequest).WithPayload("Error. Your Code is invalid!")
	}

	if err := templ.Execute(buf, content); err != nil {
		groupLog.Error(`Can't execute template'`, err)
		return group.NewJoinGroupHelpOK().WithPayload("Internal Server Error")
	}

	return group.NewJoinGroupHelpOK().WithPayload(buf.String())
}

func LeaveGroup(params group.LeaveGroupParams, principal *models.User) middleware.Responder {
	if err := principal.LeaveGroup(); err != nil {
		groupLog.Critical("Database error updating group!", err)
		return NewInternalServerError("Internal Database Error")
	}

	return group.NewLeaveGroupOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully left group"),
		Status:  swag.Int64(http.StatusOK),
	})
}

func UpdateGroupImage(params group.UpdateGroupImageParams, principal *models.User) middleware.Responder {
	groupLog.Debug("Start put group image")

	var (
		err error
		g   *models.Group
	)

	// Database
	if g, err = models.GetGroupByUID(params.GroupID); models.IsErrUserNotExist(err) {
		groupLog.Debugf(`Can't find database group with id "%s"!`, params.GroupID)
		return NewNotFoundResponse("Unknown group")

	} else if err != nil {
		groupLog.Critical("Database Error getting group!", err)
		return NewInternalServerError("Internal Server Error")
	}

	// TODO: use isAdmin
	if !g.HasMember(*principal.UID) {
		return NewUnauthorizedResponse("User not a member of the group.")
	}

	data, err := ioutil.ReadAll(params.ProfileImage)
	if err != nil {
		return NewInternalServerError("Internal Server Error")
	}

	if !base.IsFileJPG(data) {
		msg := fmt.Sprintf(
			`Invalid file type. Only "image/jpeg" allowed. Mime was "%s"`,
			base.GetMimeType(data),
		)
		userLog.Debugf(`Invalid mime type "%s"`, msg)
		return NewBadRequest(msg)
	}

	if err = g.UploadGroupImage(data); err != nil {
		userLog.Critical(`Error uploading group avatar.`)
		return NewInternalServerError("Internal Server Error")
	}

	return group.NewUpdateGroupImageOK().WithPayload(&models.SuccessResponse{
		Message: swag.String("Successfully uploaded image file"),
		Status:  swag.Int64(http.StatusOK),
	})
}