package service

import (
	"errors"
	"go-backend/api/common/common"
	"go-backend/api/server/dao"
	"go-backend/api/server/entity"
	"go-backend/api/server/tools/server"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary API of golang gin backend
// @Tags User
// @description user register : 用户注册 参数列表：[用户名、密码、电话号码、邮箱地址] 
// @version 1.0
// @accept mpfd
// @param Name formData string true "username"
// @param Password formData string true "password"
// @param Telephone formData string true "telephone"
// @param Email formData string true "email"
// @Success 200 {object} server.SuccessResponse200 "注册成功"
// @Failure 422 {object} server.FailureResponse422 "输入参数错误"
// @Failure 500 {object} server.FailureResponse500 "系统异常"
// @router /user/register [post]
func RegisterService(name string, password string, telephone string, email string) (uint, string, error) {
	if len(name) < 2 {
		err := errors.New(server.NameTooShort)
		return 0, "", err
	}
	if len(password) < 6 {
		err := errors.New(server.PasswordTooShort)
		return 0, "", err
	}
	if user := dao.GetUserInfoByName(name); user.ID != 0 {
		err := errors.New(server.UsernameAlreadyExist)
		return 0, "", err
	}
	hasePassword, errEncryp := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errEncryp != nil {
		err := errors.New(server.PasswordEncryptionFailed)
		return 0, "", err
	}
	password = string(hasePassword)
	id := dao.CreateUser(entity.User{
		Name: name,
		Password: password,
		Telephone: telephone,
		Email: email,
	})
	token, errReleaseToken := common.ReleaseToken(id)
	if errReleaseToken != nil {
		err := errors.New(server.TokenGenerateFailed)
		return 0, "", err
	}
	return id, token, nil
}

// @Summary API of golang gin backend
// @Tags User
// @description user login : 用户登录 参数列表：[用户名、密码] 
// @version 1.0
// @accept mpfd
// @param Name formData string true "username"
// @param Password formData string true "password"
// @Success 200 {object} server.SuccessResponse200 "登录成功"
// @Failure 422 {object} server.FailureResponse422 "输入参数错误"
// @Failure 500 {object} server.FailureResponse500 "系统异常"
// @router /user/login [post]
func LoginService(name string, password string) (uint, string, uint, error) {
	if len(name) < 2 {
		panic((server.NameTooShort))
	}
	if len(password) < 6 {
		err := errors.New(server.PasswordTooShort)
		return 0, "", 0, err
	}
	user := dao.GetUserInfoByName(name)
	if user.ID == 0 {
		panic(server.UserNotExist)
	}
	if errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); errPassword != nil {
		panic(server.WrongPassword)
	}
	token, errToken := common.ReleaseToken(user.ID)
	if errToken != nil {
		panic(server.TokenGenerateFailed)
	}
	return user.ID, token, user.DefaultCompany, nil
}

// @Summary API of golang gin backend
// @Tags User
// @description get user information : 获取当前用户的详细信息 参数列表：[] 访问携带token
// @version 1.0
// @accept application/json
// @param Id query string true "Id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "查询成功"
// @Failure 400 {object} server.FailureResponse400 "用户信息不存在"
// @Failure 401 {object} server.FailureResponse401 "权限不足"
// @router /user/info [get]
func InfoService(id uint) (*gin.H, error) {
	user := dao.GetUserInfoById(id)
	if user.ID == 0 {
		err := errors.New(server.UserNotExist)
		return nil, err
	}
	infoMap := gin.H{
		"name":        user.Name,
		"id":          user.ID,
		"create_time": user.CreatedAt,
		"update_time": user.UpdatedAt,
		"telephone":   user.Telephone,
		"email":	   user.Email,
	}
	return &infoMap, nil
}

// @Summary API of golang gin backend
// @Tags User
// @description update user information : 更新当前用户的详细信息 参数列表：[] 访问携带token
// @version 1.0
// @accept application/json
// @param Name query string true "username"
// @param Password query string true "password"
// @param Telephone query string true "telephone"
// @param Email query string true "email"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "更新成功"
// @Failure 401 {object} server.FailureResponse401 "权限不足"
// @router /user/update [put]
func UpdateUserInfoService(userId uint, name string, password string, telephone string, email string) {
	if len(name) < 2 {
		err := errors.New(server.NameTooShort)
		panic(err.Error())
	}
	if len(password) < 6 {
		err := errors.New(server.PasswordTooShort)
		panic(err.Error())
	}
	if user := dao.GetUserInfoByName(name); user.ID != 0 {
		err := errors.New(server.UsernameAlreadyExist)
		panic(err.Error())
	}
	hasePassword, errEncryp := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errEncryp != nil {
		err := errors.New(server.PasswordEncryptionFailed)
		panic(err.Error())
	}
	password = string(hasePassword)
	dao.UpdateUserInfo(userId, name, password, telephone, email)
}

// @Summary API of golang gin backend
// @Tags User
// @description update user default company : 更新当前用户的首页默认显示企业 参数列表：[公司ID] 访问携带token
// @version 1.0
// @accept application/json
// @param CompanyId query int true "company id"
// @param Authorization header string true "token"
// @Success 200 {object} server.SuccessResponse200 "更新成功"
// @Failure 401 {object} server.FailureResponse401 "权限不足"
// @router /user/update_default_company [put]
func UpdateUserDefaultCompanyService(userId uint, companyId uint) {
	dao.UpdateUserDefaultCompany(userId, companyId)
}