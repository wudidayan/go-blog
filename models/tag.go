package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageOffset int, pageSize int, maps interface{}) (tags []Tag, err error) {
	err = db.Where(maps).Offset(pageOffset).Limit(pageSize).Find(&tags).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func GetTagTotal(maps interface{}) (count int, err error) {
	err = db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	var err error

	err = db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	var err error

	err = db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil {
		return err
	}

	return nil
}

func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
