package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go.uber.org/zap"
)

const secret = "theo5970"

// CheckUserExist check user exist
func CheckUserExist(username string) (err error) {
	sqlstr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		zap.L().Error(err.Error())
		return err
	}
	if count > 0 {

		return ErrorUserExist
	}
	return
}

// InsertUser insert user info
func InsertUser(user *models.User) error {
	//execute SQL insert
	user.Password = encryptPassword(user.Password)
	sqlStr := `insert into user (user_id,username,password) values (?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
func Login(user *models.User) (err error) {
	sqlstr := `select user_id,username,password from user where username = ?`
	oPassword := user.Password

	err = db.Get(user, sqlstr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}
