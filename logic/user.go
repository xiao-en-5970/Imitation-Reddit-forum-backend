package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snow"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// judge user exist or not

	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// generate UID
	userID := snow.GenID()
	//create a user
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//save Mysql
	return mysql.InsertUser(user)
	//redis
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{

		Username: p.Username,
		Password: p.Password,
	}
	//get userID
	if err = mysql.Login(user); err != nil {
		return nil, err
	}
	//gen JWT token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
