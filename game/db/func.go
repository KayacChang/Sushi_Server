package db

import (
	"database/sql"

	"github.com/YWJSonic/ServerUtility/code"
	"github.com/YWJSonic/ServerUtility/dbservice"
	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/messagehandle"
)

// GetSetting get db setting
func GetSetting(db *sql.DB) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "SettingGet_Read")
	return result, err
}

// GetSettingKey get db setting
func GetSettingKey(db *sql.DB, key string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "SettingKeyGet_Read", key)
	return result, err
}

// NewSetting ...
func NewSetting(db *sql.DB, args ...interface{}) {
	dbservice.CallWrite(
		db,
		dbservice.MakeProcedureQueryStr("SettingNew_Write", len(args)),
		args...,
	)
}

// UpdateSetting ...
func UpdateSetting(db *sql.DB, args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbservice.CallWrite(db, dbservice.MakeProcedureQueryStr("SettingSet_Update", len(args)), args...)
	return err
}

// ReflushSetting ...
func ReflushSetting(db *sql.DB, args ...interface{}) messagehandle.ErrorMsg {
	args = append(args, foundation.ServerNowTime())
	_, err := dbservice.CallWrite(db, dbservice.MakeProcedureQueryStr("SettingSet_Update_v2", len(args)), args...)
	return err
}

// GetAttachTypeRange ...
func GetAttachTypeRange(db *sql.DB, playerid, kind, miniAttType, maxAttType int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "AttachTypeRangeGet_Read", playerid, kind, miniAttType, maxAttType)
	return result, err
}

// GetAttachType ...
func GetAttachType(db *sql.DB, playerid int64, kind int64, attType int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "AttachTypeGet_Read", playerid, kind, attType)
	return result, err
}

// GetAttachKind get db attach kind
func GetAttachKind(db *sql.DB, playerid int64, kind int64) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "AttachKindGet_Read", playerid, kind)
	return result, err
}

// NewAttach ...
func NewAttach(db *sql.DB, args ...interface{}) {
	dbservice.CallWrite(
		db,
		dbservice.MakeProcedureQueryStr("AttachNew_Write", len(args)),
		args...,
	)
}

// UpdateAttach ...
func UpdateAttach(db *sql.DB, args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbservice.CallWrite(db, dbservice.MakeProcedureQueryStr("AttachSet_Write", len(args)), args...)
	return err
}

// GetAccountInfo Check Account existence and get
func GetAccountInfo(db *sql.DB, account string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "AccountGet_Read", account)
	return result, err
}

// NewAccount new goruting set new Account
func NewAccount(db *sql.DB, args ...interface{}) { //messagehandle.ErrorMsg {
	dbservice.CallWrite(
		db,
		dbservice.MakeProcedureQueryStr("AccountNew_Write", len(args)),
		args...,
	)
}

// UpdateAccount update
func UpdateAccount(db *sql.DB, args ...interface{}) messagehandle.ErrorMsg {
	_, err := dbservice.CallWrite(db, dbservice.MakeProcedureQueryStr("AccountSet_Update", len(args)), args...)
	return err
}

// NewGameAccount gameaccount, money, gametoken
func NewGameAccount(db *sql.DB, args ...interface{}) (int64, messagehandle.ErrorMsg) {
	QuertStr := "INSERT INTO gameaccount VALUE (NULL,"
	if len(args) > 0 {
		for range args {
			QuertStr += "?,"
		}
		QuertStr = QuertStr[:len(QuertStr)-1]
	}
	QuertStr += ");"

	request, err := dbservice.Exec(db, QuertStr, args...)
	if err.ErrorCode != code.OK {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = "NewGameAccountError"
		messagehandle.ErrorLogPrintln("NewGameAccount-1", err, QuertStr, args)
		return -1, err
	}
	playerID, errMsg := request.LastInsertId()
	if errMsg != nil {
		messagehandle.ErrorLogPrintln("NewGameAccount-2", errMsg)
	}
	// err := messagehandle.New()
	return playerID, err
}

// GetPlayerInfoByGameAccount ...
func GetPlayerInfoByGameAccount(db *sql.DB, gameAccount string) ([]map[string]interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "GameAccountGet_Read", gameAccount)
	return result, err
}

// GetPlayerInfoByPlayerID ...
func GetPlayerInfoByPlayerID(db *sql.DB, playerID int64) (interface{}, messagehandle.ErrorMsg) {
	result, err := dbservice.CallReadOutMap(db, "GameAccountGet_Read", playerID)
	return result, err
}

// UpdatePlayerInfo ...
func UpdatePlayerInfo(db *sql.DB, args ...interface{}) {
	dbservice.CallWrite(
		db,
		dbservice.MakeProcedureQueryStr("GameAccountSet_Update", len(args)),
		args...,
	)
}
