package v1

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/wudidayan/go-blog/models"
	"github.com/wudidayan/go-blog/pkg/app"
	"github.com/wudidayan/go-blog/pkg/e"
	"github.com/wudidayan/go-blog/pkg/gredis"
	"github.com/wudidayan/go-blog/pkg/logging"
	"github.com/wudidayan/go-blog/pkg/setting"
	"github.com/wudidayan/go-blog/pkg/util"
	"github.com/wudidayan/go-blog/service/cache"
)

// @Summary 获取单个文章
// @Description 通过id获取指定文章
// @Produce json
// @Param id path int true "ID"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	appResp := app.Gin{c}
	code := e.SUCCESS
	var data interface{}

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	var cacheArticle *models.Article
	cache_article := cache.Article{ID: id}
	key := cache_article.GetArticleKey()
	if gredis.Exists(key) {
		keydata, err := gredis.Get(key)
		if err != nil {
			logging.Warn("redis.Get() error: ", err)
		} else {
			json.Unmarshal(keydata, &cacheArticle)
			code = e.SUCCESS
			appResp.Response(http.StatusOK, code, cacheArticle)
			return
		}
	} else {
		logging.Debug("redis key[%s] not exist", key)
	}

	exist, err := models.ExistArticleByID(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if !exist {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	data = models.GetArticle(id)
	err = gredis.Set(key, data, 3600)
	if err != nil {
		logging.Warn("set redis key[%s] error[%v]", key, err)
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, data)
	return
}

// @Summary 获取多个文章
// @Produce json
// @Param tag_id query int false "TagID"
// @Param state query int false "State"
// @Param page query int false "PageNum"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	appResp := app.Gin{c}
	code := e.SUCCESS
	var err error

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	data["lists"], err = models.GetArticles(util.GetPageOffset(c), setting.AppSetting.PageSize, maps)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	data["total"], err = models.GetArticleTotal(maps)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, data)
}

// @Summary 新增文章
// @Produce json
// @Param tag_id query int true "TagID"
// @Param title query string true "Title"
// @Param desc query string true "Desc"
// @Param content query string true "Content"
// @Param created_by query string true "CreatedBy"
// @Param state query int true "State"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err := models.ExistTagByID(tagId)
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
	data["tag_id"] = tagId
	data["title"] = title
	data["desc"] = desc
	data["content"] = content
	data["created_by"] = createdBy
	data["state"] = state

	err = models.AddArticle(data)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
	return
}

// @Summary 修改文章
// @Produce json
// @Param id path int true "ID"
// @Param tag_id query string true "TagID"
// @Param title query string false "Title"
// @Param desc query string false "Desc"
// @Param content query string false "Content"
// @Param modified_by query string true "ModifiedBy"
// @Param state query int false "State"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")

	var state int = -1
	valid := validation.Validation{}
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	var exist bool
	var err error

	exist, err = models.ExistArticleByID(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if !exist {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err = models.ExistTagByID(tagId)
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

	if tagId > 0 {
		data["tag_id"] = tagId
	}

	if title != "" {
		data["title"] = title
	}

	if desc != "" {
		data["desc"] = desc
	}

	if content != "" {
		data["content"] = content
	}

	err = models.EditArticle(id, data)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	cache_article := cache.Article{ID: id}
	key := cache_article.GetArticleKey()
	gredis.Delete(key)

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
}

// @Summary 删除文章
// @Produce json
// @Param id path int true "ID"
// @Param token query string true "Token"
// @Success 200 {object} app.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	appResp := app.Gin{c}
	code := e.SUCCESS

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}

		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	exist, err := models.ExistArticleByID(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	if !exist {
		code = e.ERROR_NOT_EXIST_ARTICLE
		appResp.Response(http.StatusOK, code, nil)
		return
	}

	err = models.DeleteArticle(id)
	if err != nil {
		code = e.ERROR
		appResp.Response(http.StatusInternalServerError, code, nil)
		return
	}

	cache_article := cache.Article{ID: id}
	key := cache_article.GetArticleKey()
	gredis.Delete(key)

	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, nil)
	return
}
