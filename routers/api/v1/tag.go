package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/wudidayan/go-blog/models"
	"github.com/wudidayan/go-blog/pkg/app"
	"github.com/wudidayan/go-blog/pkg/e"
	"github.com/wudidayan/go-blog/pkg/setting"
	"github.com/wudidayan/go-blog/pkg/util"
)

// @Summary 获取多个文章标签
// @Produce json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Param page query int false "PageNum"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	name := c.Query("name")
	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	var err error
	appResp := app.Gin{c}
	code := e.SUCCESS

	data["lists"], err = models.GetTags(util.GetPageOffset(c), setting.AppSetting.PageSize, maps)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	data["total"], err = models.GetTagTotal(maps)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, data)
	return
}

// @Summary 添加文章标签
// @Produce json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query string true "CreatedBy"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err := models.ExistTagByName(name)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if exist {
		code = e.ERROR_EXIST_TAG
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	err = models.AddTag(name, state, createdBy)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
	return
}

// @Summary 修改文章标签
// @Produce json
// @Param id path int true "ID"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	valid := validation.Validation{}
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err := models.ExistTagByID(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if !exist {
		code = e.ERROR_NOT_EXIST_TAG
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	data := make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}
	if state != -1 {
		data["state"] = state
	}

	err = models.EditTag(id, data)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
	return
}

// @Summary 删除文章标签
// @Produce json
// @Param id path int true "ID"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err := models.ExistTagByID(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if !exist {
		code = e.ERROR_NOT_EXIST_TAG
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	err = models.DeleteTag(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
	return
}
