package controller

import (
	"crypto/sha256"
	"github.com/YJ9938/DouYin/model"
	"github.com/YJ9938/DouYin/service"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
)

const (
	maxUsernameLen = 32
	maxPasswordLen = 32
	randomSaltLen  = 32
)

// generateSalt generates a salt of a given length.
func generateSalt(length int) []byte {
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = byte(rand.Intn(256))
	}
	return salt
}

// hashPassword calculates SHA256(plain password, salt).
func hashPassword(password string, salt []byte) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write(salt)
	return h.Sum(nil)
}

func checkUserValid(username, password string) bool {
	return len(username) != 0 && len(username) <= maxUsernameLen &&
		len(password) != 0 && len(password) <= maxPasswordLen
}

//用户注册、登录返回的响应
type UserResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

//用户查询返回的响应
type UserInfoResponse struct {
	Response
	User *model.UserInfo `json:"user"`
}

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if !checkUserValid(username, password) {
		Error(c, 1, "用户名或密码不符合规范")
	}
	// Generate random salt.
	salt := generateSalt(32)
	hashedPassword := hashPassword(password, salt)
	newUser := model.User{
		Username: username,
		Password: hashedPassword,
		Salt:     salt,
	}

	// Insert a user record.
	userDao := model.NewUserDao()
	if err := userDao.AddUser(&newUser); err != nil {
		if err == model.ErrUserExists {
			Error(c, 400, "用户名已存在")
		} else {
			log.Printf("error when adding user %v: %s\n", newUser, err)
			Error(c, 500, err.Error())
		}
		return
	}

	token, err := signJWT(newUser.Id)
	if err != nil {
		InternalError(c)
		return
	}
	c.JSON(200, UserResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserId: newUser.Id,
		Token:  token,
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if !checkUserValid(username, password) {
		Error(c, 1, "用户名或密码不符合要求")
		return
	}
	// Query the user by the username.
	userDao := model.NewUserDao()
	user, err := userDao.QueryUserByName(username)
	if err != nil {
		if err == model.ErrUserNotExists {
			Error(c, 400, "用户名不存在")
		} else {
			log.Printf("error when query user by ID: %s\n", err)
		}
		return
	}

	// Check if username and password matches.
	hashedPassword := hashPassword(password, user.Salt)
	if !reflect.DeepEqual(hashedPassword, user.Password) {
		Error(c, 400, "用户名密码不匹配")
		return
	}

	token, err := signJWT(user.Id)
	if err != nil {
		InternalError(c)
		return
	}
	c.JSON(200, UserResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserId: int64(user.Id),
		Token:  token,
	})
}

func UserInfo(c *gin.Context) {
	//取出待查询用户ID
	rawId := c.Query("user_id")
	queryId, _ := strconv.ParseInt(rawId, 10, 64)
	//从token取出当前用户ID
	token := c.Query("token")
	rawCurrentId := parseToken(token).Id
	currentId, _ := strconv.ParseInt(rawCurrentId, 10, 64)
	userInfoService := service.UserInfoService{
		CurrentUser: currentId,
		QueryUser:   queryId,
	}
	user, err := userInfoService.QueryUserInfoById()
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: Response{StatusCode: 0},
		User:     user,
	})
}
