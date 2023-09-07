package models

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

// 标签模型
type Tag struct {
	gorm.Model
	Name     string    `json:"name"`                                              // 标签名称
	Identify string    `json:"identify"`                                          // 唯一标识
	Num      int       `json:"num"`                                               // 标签文章数
	Articles []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"` // 标签关联的文章
}

func NewTag() *Tag {
	return &Tag{}
}

func (tag *Tag) Create(name string) error {
	newTag := Tag{
		Name:     name,
		Identify: utils.GenerateMD5Hash(name),
		Num:      0,
	}
	if err := global.GVA_DB.Create(&newTag).Error; err != nil {
		return err
	}
	return nil
}

func (tag *Tag) Update(updates map[string]interface{}) error {
	if name, ok := updates["name"].(string); ok {
		// 如果存在，更新 identify
		updates["identify"] = utils.GenerateMD5Hash(name)
	}
	// 更新字段
	tx := global.GVA_DB.Begin()
	if err := tx.Model(tag).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *Tag) DeleteByTagIds(tagIds []int) (int64, error) {
	result := global.GVA_DB.Where("id IN ?", tagIds).Unscoped().Delete(&Tag{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (tag *Tag) FindAll() ([]*Tag, error) {
	var tags []*Tag
	err := global.GVA_DB.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *Tag) Count() int {
	var total int64
	global.GVA_DB.Model(&Tag{}).Count(&total)
	return int(total)
}

func (t *Tag) FindByTagId(tagId int) (*Tag, error) {
	var tag *Tag
	err := global.GVA_DB.Where("id = ?", tagId).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (c *Tag) FindByTagIdArray(tagIds []int) ([]Tag, error) {
	var tags []Tag
	result := global.GVA_DB.Where("id IN ?", tagIds).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (c *Tag) FindByName(name string) (*Tag, error) {
	var tag Tag
	err := global.GVA_DB.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (c *Tag) FindByIdentify(identify string) (*Tag, error) {
	var tag Tag
	err := global.GVA_DB.Where("identify = ?", identify).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (tag *Tag) FindByKeyword(keyword string) ([]*Tag, error) {
	var tags []*Tag
	err := global.GVA_DB.Where("name LIKE ?", "%"+keyword+"%").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
