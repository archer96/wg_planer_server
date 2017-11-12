package controllers

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/op/go-logging"
	"github.com/wgplaner/wg_planer_server/gen/models"
)

var authLog = logging.MustGetLogger("Auth")

func UserIDAuth(token string) (interface{}, error) {
	theUser := models.User{UID: &token}
	authLog.Debugf(`Check userId authorization for user id "%s"`, *theUser.UID)

	if isRegistered, err := isUserRegistered(&theUser); err != nil {
		authLog.Error(`DB error with isUserRegistered`, err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		authLog.Debugf(`Unauthorized database user "%s"`, *theUser.UID)
		return nil, errors.Unauthenticated("invalid credentials")
	}

	return theUser, nil
}

func FirebaseIDAuth(token string) (interface{}, error) {
	theUser := models.User{UID: &token}
	authLog.Debugf(`Check firebaseId authorization for user id "%s"`, *theUser.UID)

	if isRegistered, err := isUserOnFirebase(&theUser); err != nil {
		authLog.Error(`DB error with isUserOnFirebase`, err.Error())
		return nil, errors.New(http.StatusInternalServerError, "Internal Server Error")

	} else if !isRegistered {
		authLog.Debugf(`Unauthorized firebase user "%s"`, *theUser.UID)
		return nil, errors.Unauthenticated("invalid credentials")
	}

	return theUser, nil
}
