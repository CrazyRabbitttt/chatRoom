package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goHttp/pkg/utils"
	"net/http"
	"time"
)

const ExpireTime = 7 * 24 // 过期时间为7天

var jwtKey = []byte("shao-gui-xin.sgx")

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// 颁发 token
func ReleaseToken(userId int64) (string, error) {
	expirationTime := time.Now().Add(ExpireTime * time.Hour)
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "lovexb",
			Subject:   "S_G_X",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// parse the token(string) To Claims
func ParseToken(tokenString string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}

// 鉴权的中间件， 鉴权并且设置 user_id(能够用 gin.Context 在中间件之间传递参数)
func JWTMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		// 真的就是不存在 token
		if tokenStr == "" {
			c.JSON(http.StatusOK, utils.CommonResponse{StatusCode: utils.Invalid, StatusMsg: "用户不存在"})
			c.Abort()
			return
		}
		// 验证token
		tokenStruct, ok := ParseToken(tokenStr)
		if !ok { // 如果说 token 验证不通过
			c.JSON(http.StatusOK, utils.CommonResponse{StatusCode: utils.Invalid, StatusMsg: "token不正确"})
			c.Abort()
			return
		}
		// token 已经是超时的了
		if time.Now().Unix() > tokenStruct.ExpiresAt {
			c.JSON(http.StatusOK, utils.CommonResponse{
				StatusCode: utils.Invalid,
				StatusMsg:  "token已经过期",
			})
			c.Abort()
			return
		}
		c.Set("user_id", tokenStruct.UserId) // 设置好 user_id, 后面能够 get 到
		c.Next()
	}
}
