package controller

import (
	"admin_demo/common"
	"admin_demo/dto"
	"admin_demo/model"
	"admin_demo/response"
	"admin_demo/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

/**
用户注册
 */
func Register(c *gin.Context) {
	db := common.GetDB()
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	//数据验证
	if len(telephone) < 11{
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6{
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码必须不小于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}
	//创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "密码加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Telephone: telephone,
		Password: string(hashedPassword),
	}
	db.Create(&newUser)
	fmt.Println(newUser)
	response.Success(c, nil, "注册成功")
}


/**
用户登录
 */
func Login(c *gin.Context) {
	db := common.GetDB()
	//获取参数
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	//数据验证
	if len(telephone) < 11{
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6{
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码必须不小于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	//验证密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "token发放失败")
		log.Printf("token generate err: %v", err)
		return
	}

	//返回结果
	response.Success(c, gin.H{"token":token}, "登录成功")
}

/**
从上下文中获取登录用户的信息
 */
func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code":200, "data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
}


func isTelephoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("telephone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
