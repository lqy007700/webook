package web

import (
	"errors"
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
	svc     *service.UserService
	codeSvc *service.CodeService
}

func NewUserHandler(svc *service.UserService, codeSvc *service.CodeService) *UserHandler {
	return &UserHandler{
		svc:     svc,
		codeSvc: codeSvc,
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
	ug.POST("/login_sms/code/send", u.SendLoginSmsCode)
	ug.POST("/login_sms/code/login", u.LoginSms)
}

func (u *UserHandler) LoginSms(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "参数错误")
		return
	}

	const biz = "login"
	err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "验证码错误")
		return
	}

	user, err := u.svc.LoginSms(ctx.Request.Context(), req.Phone)
	if err != nil {
		ctx.String(http.StatusBadRequest, "登录失败")
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

func (u *UserHandler) SendLoginSmsCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	const biz = "login"
	err := u.codeSvc.Send(ctx, biz, req.Phone)
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "短信发送成功")
	case errors.Is(err, service.ErrCodeSendTooMany):
		ctx.String(http.StatusOK, "短信发送太频繁")
	default:
		ctx.String(http.StatusOK, "短信发送失败")
	}
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

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		ctx.String(http.StatusBadRequest, "参数错误")
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
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	u.svc.Profile(ctx, 1)

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
