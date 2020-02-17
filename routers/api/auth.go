package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/wudidayan/go-blog/models"
	"github.com/wudidayan/go-blog/pkg/app"
	"github.com/wudidayan/go-blog/pkg/e"
	"github.com/wudidayan/go-blog/pkg/logging"
	"github.com/wudidayan/go-blog/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary 获取权限
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	appResp := app.Gin{c}
	code := e.SUCCESS
	data := make(map[string]string)

	if !ok {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, data)
		return
	}

	isExist := models.CheckAuth(username, password)
	if !isExist {
		code = e.ERROR_AUTH
		appResp.Response(http.StatusOK, code, data)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		code = e.ERROR_AUTH_TOKEN
		appResp.Response(http.StatusOK, code, data)
		return
	}

	data["token"] = token
	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, data)
	return
}
