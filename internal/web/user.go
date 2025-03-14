package web

import (
	"fmt"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) RegisterRouters(s *gin.Engine) {
	ug := s.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/loginjwt", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
	ug.GET("/profilejwt", u.ProfileJWT)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	err := u.svc.SignUp(ctx.Request.Context(), domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		ctx.String(http.StatusOK, "注册失败")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "登录失败")
		return
	}

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("eHwX09d&*3KLs0^lm#PqA5RzVcT7NyU4QbFiGj2M8W!n@tYh"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "登录失败")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功")
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx.Request.Context(), req.Email, req.Password)
	fmt.Println(err)
	if err != nil {
		ctx.String(http.StatusOK, "登录失败")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
}

func (*UserHandler) Edit(ctx *gin.Context) {
}

func (*UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "登录成功")
}
func (*UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claim")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"uid": claims.Uid,
	})
}
