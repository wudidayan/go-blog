package models

import (
	"github.com/jinzhu/gorm"

	"time"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetArticleTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Article{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetArticles(pageOffset int, pageSize int, maps interface{}) (articles []Article, err error) {
	err = db.Preload("Tag").Where(maps).Offset(pageOffset).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func EditArticle(id int, data interface{}) error {
	err := db.Model(&Article{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}

func AddArticle(data map[string]interface{}) error {
	err := db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	err := db.Model(&Article{}).Where("id = ?", id).Delete(Article{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
