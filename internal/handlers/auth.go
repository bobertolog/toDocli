package handlers

import (
	"net/http"
	"time"

	"todocli/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
}

var (
	userRepo  service.UserRepository
	jwtSecret = []byte("mysecret") // Используется и здесь, и в middleware
)

// SetUserRepo устанавливает реализацию UserRepository (например, PostgresRepo)
func SetUserRepo(repo service.UserRepository) {
	userRepo = repo
}

// GetJWTSecret возвращает секрет (используется в middleware)
func GetJWTSecret() []byte {
	return jwtSecret
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует пользователя с логином и паролем
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   body body registerReq true "Данные для регистрации"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := userRepo.Register(req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

// Login godoc
// @Summary Авторизация
// @Description Возвращает JWT токен при успешной авторизации
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   body body loginReq true "Данные для входа"
// @Success 200 {object} authResp
// @Failure 400,401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := userRepo.Login(req.Username, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, authResp{Token: signed})
}
