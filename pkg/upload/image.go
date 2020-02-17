package upload

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/wudidayan/go-blog/pkg/file"
	"github.com/wudidayan/go-blog/pkg/logging"
	"github.com/wudidayan/go-blog/pkg/setting"
	"github.com/wudidayan/go-blog/pkg/util"
)

func GetNewImageName(fileName string) string {
	ext := path.Ext(fileName)
	fileNamePrifix := strings.TrimSuffix(fileName, ext)
	randStr := util.GetRandomString(16)
	fileNameNew := util.EncodeMD5(fileNamePrifix + randStr)
	return fileNameNew + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func GetImagePrefixUrl() string {
	return setting.AppSetting.PrefixUrl
}

func CheckImageExt(fileName string) bool {
	ext := path.Ext(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func GetImageSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)
	return len(content), err
}

func CheckImageSize(f multipart.File) bool {
	size, err := GetImageSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImagePath(filePath string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + filePath)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(dir + "/" + filePath)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied: %s", filePath)
	}

	return nil
}
