package middleware

import (
	"admin_demo/common"
	"admin_demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("authorization")

		//验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString,"Bearer") {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg":"权限不足"})
			context.Abort()
			return
		}

		//截取头部"Bearer"以外的字符串
		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg":"权限不足"})
			context.Abort()
			return
		}

		//验证通过后获取claims中的userId
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户不存在
		if user.ID == 0 {
			context.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg":"权限不足"})
			context.Abort()
			return
		}

		//用户存在，将信息存入上下文
		context.Set("user", user)
		context.Next()
	}
}



