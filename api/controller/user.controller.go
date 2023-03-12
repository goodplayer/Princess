package controller

import (
	"errors"
	"net/http"

	"github.com/goodplayer/Princess/domain/user"
	"github.com/goodplayer/Princess/framework/app"
	"github.com/goodplayer/Princess/framework/ginplugin"
	"github.com/goodplayer/Princess/framework/ginsupport"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var userControllerLogger = logrus.New()

func initUser(r *gin.Engine, ac *app.ApplicationContainer) error {
	uc := new(UserController)
	ac.Set(uc)

	userGroup := r.Group("/user")
	{
		ginsupport.GroupGET(userGroup, "/register", uc.showRegister)
		ginsupport.GroupPOST(userGroup, "/register", uc.doUserRegister)

		ginsupport.GroupGET(userGroup, "/login", uc.showLogin)
		ginsupport.GroupPOST(userGroup, "/login", uc.doLogin)
	}

	return nil
}

type UserController struct {
	userDomainService *user.UserDomainService
}

func (uc *UserController) DependOn(userDomainService *user.UserDomainService) {
	uc.userDomainService = userDomainService
}

func (uc *UserController) showRegister(ctx *gin.Context) ginsupport.Render {
	return ginsupport.NewRenderTemplateSuccess("default/user/register.html", gin.H{})
}

func (uc *UserController) doUserRegister(ctx *gin.Context) ginsupport.Render {
	params := new(struct {
		UserName string `form:"email" binding:"required,email,max=120"`
		Password string `form:"password" binding:"required,min=8,max=20"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		return ginsupport.NewErrorTemplate(http.StatusBadRequest, "default/user/register.html", err)
	}

	if exists, err := uc.userDomainService.UserExistsByEmail(params.UserName); err != nil {
		userControllerLogger.Errorf("check user exists by email error: %s", err)
		return ginsupport.NewErrorTemplate(http.StatusInternalServerError, "default/user/register.html", errors.New("internal error"))
	} else if exists {
		return ginsupport.NewErrorTemplate(http.StatusOK, "default/user/register.html", errors.New("user exists"))
	}

	if u, err := uc.userDomainService.DoUserRegister(&struct {
		Email    string
		Password string
	}{
		Email:    params.UserName,
		Password: params.Password,
	}); err != nil {
		userControllerLogger.Errorf("do register user error: %s", err)
		return ginsupport.NewErrorTemplate(http.StatusInternalServerError, "default/user/register.html", errors.New("register user failed"))
	} else {
		var _ = u
	}

	return ginsupport.NewRenderTemplateSuccess("default/user/register_success.html", gin.H{})
}

func (uc *UserController) showLogin(ctx *gin.Context) ginsupport.Render {
	return ginsupport.NewRenderTemplateSuccess("default/user/login.html", gin.H{})
}

func (uc *UserController) doLogin(ctx *gin.Context) ginsupport.Render {
	params := new(struct {
		UserName string `form:"username" binding:"required,max=120"`
		Password string `form:"password" binding:"required,min=8,max=20"`
	})

	if err := ctx.ShouldBind(params); err != nil {
		return ginsupport.NewErrorTemplate(http.StatusBadRequest, "default/user/login.html", err)
	}

	if u, err := uc.userDomainService.UserAuthByPassword(params.UserName, params.Password); err != nil {
		return ginsupport.NewErrorTemplate(http.StatusBadRequest, "default/user/login.html", errors.New("login failed"))
	} else {
		// init session with user
		session := ginplugin.Session(ctx)
		session.Set("user_id", []byte(u.UserId))
		if err := session.SaveAndFreeze(); err != nil {
			userControllerLogger.Errorf("save login user session failed: %s", err)
			return ginsupport.NewErrorTemplate(http.StatusBadRequest, "default/user/login.html", errors.New("login failed"))
		}
	}

	return ginsupport.NewRenderTemplateSuccess("default/user/login_success.html", gin.H{
		"redirect_url": "/",
	})
}

func (uc *UserController) userLogout() {

}

func (uc *UserController) forgetPassword() {

}
