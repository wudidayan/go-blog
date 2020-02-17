package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wudidayan/go-blog/pkg/app"
	"github.com/wudidayan/go-blog/pkg/e"
	"github.com/wudidayan/go-blog/pkg/logging"
	"github.com/wudidayan/go-blog/pkg/upload"
)

// @Summary 上传图片
// @Accept mpfd
// @Produce json
// @Param upload formData file true "Image File"
// @Success 200 {object} app.Response
// @Router /upload [post]
func UploadImage(c *gin.Context) {
	appResp := app.Gin{c}
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("upload")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		appResp.Response(http.StatusOK, code, data)
		return
	}

	if image == nil {
		code = e.INVALID_PARAMS
		appResp.Response(http.StatusOK, code, data)
		return
	}

	if !upload.CheckImageExt(image.Filename) || !upload.CheckImageSize(file) {
		code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		appResp.Response(http.StatusOK, code, data)
		return
	}

	filePath := upload.GetImageFullPath()
	fileDate := time.Now().Format("20060102")
	filePath = filePath + fileDate + "/"
	err = upload.CheckImagePath(filePath)
	if err != nil {
		logging.Warn(err)
		code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		appResp.Response(http.StatusOK, code, data)
		return
	}

	fileName := upload.GetNewImageName(image.Filename)
	err = c.SaveUploadedFile(image, filePath+fileName)
	if err != nil {
		logging.Warn(err)
		code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
		appResp.Response(http.StatusOK, code, data)
		return
	}

	data["image_url"] = upload.GetImagePrefixUrl() + "/" + upload.GetImagePath() + fileDate + "/" + fileName
	data["image_name"] = fileName
	code = e.SUCCESS
	appResp.Response(http.StatusOK, code, data)
	return
}
