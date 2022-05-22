package api

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"reflect"

	"github.com/aeof/douyin/model"
	"github.com/gin-gonic/gin"
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

// checkUserPass checks if username and password is valid.
func checkUserPass(username, password string) bool {
	return len(username) != 0 && len(username) <= maxUsernameLen &&
		len(password) != 0 && len(password) <= maxPasswordLen
}

type UserResponse struct {
	Status
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Register is a router to register a user.
func Register(c *gin.Context) {
	// Check username and password.
	username := c.Query("username")
	password := c.Query("password")
	if !checkUserPass(username, password) {
		Error(c, 400, "username or password invalid")
		return
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
	if err := model.AddUser(&newUser); err != nil {
		if err == model.ErrUserExists {
			Error(c, 400, "用户名已存在")
		} else {
			log.Printf("error when adding user %v: %s\n", newUser, err)
			Error(c, 500, err.Error())
		}
		return
	}

	token, err := signJWT(newUser.ID)
	if err != nil {
		InternalError(c)
		return
	}
	c.JSON(200, UserResponse{
		Status: Status{
			StatusCode:    0,
			StatusMessage: "注册成功",
		},
		UserID: int64(newUser.ID),
		Token:  token,
	})
}

// Login is a router to login a user.
func Login(c *gin.Context) {
	// Check username and password.
	username := c.Query("username")
	password := c.Query("password")
	if !checkUserPass(username, password) {
		Error(c, 1, "用户名或密码不符合要求")
		return
	}

	// Query the user by the username.
	user, err := model.GetUserByUsername(username)
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

	token, err := signJWT(user.ID)
	if err != nil {
		InternalError(c)
		return
	}
	c.JSON(200, UserResponse{
		Status: Status{
			StatusCode:    0,
			StatusMessage: "登录成功",
		},
		UserID: int64(user.ID),
		Token:  token,
	})
}

type UserInfoResponse struct {
	Status
	User struct {
		ID            int64
		Name          string
		FollowCount   int64 `json:"follow_count"`
		FollowerCount int64 `json:"follower_count"`
		IsFollow      bool  `json:"is_follow"`
	}
}

// QueryUserInfo queries the info of a user.
// Note it is an authorized API.
func QueryUserInfo(c *gin.Context) {
	id := getUserID(c)
	var targetID int64
	fmt.Sscanf(c.Query("user_id"), "%d", &targetID)

	// Username
	username, err := model.GetUsernameByID(targetID)
	if err != nil {
		if err == model.ErrUserNotExists {
			Error(c, 400, "用户不存在")
		} else {
			log.Printf("error when querying username of id %d: %s", targetID, err)
			InternalError(c)
		}
		return
	}
	// Follower count
	followerCount, err := model.GetFollowerCount(targetID)
	if err != nil {
		log.Printf("error when querying user follower count: %s\n", err)
		InternalError(c)
		return
	}
	// Following count
	followingCount, err := model.GetFollowingCount(targetID)
	if err != nil {
		log.Printf("error when querying user following count: %s\n", err)
		InternalError(c)
		return
	}
	// Following status
	isFollowing, err := model.IsFollowing(id, targetID)
	if err != nil {
		log.Printf("error when querying following status: %s\n", err)
	}

	var resp UserInfoResponse
	resp.Status = Status{StatusCode: 0, StatusMessage: "获取成功"}
	resp.User.FollowCount = followingCount
	resp.User.FollowerCount = followerCount
	resp.User.IsFollow = isFollowing
	resp.User.Name = username
	resp.User.ID = targetID
	c.JSON(200, resp)
}
