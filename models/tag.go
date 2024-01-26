package models

import (
	"errors"

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

func (tag *Tag) FindOrCreateByName(name string) (*Tag, error) {
	obj, err := tag.FindByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if obj != nil {
		return obj, nil
	}

	return tag.Create(name)
}

func (tag *Tag) Create(name string) (*Tag, error) {
	newTag := Tag{
		Name:     name,
		Identify: utils.GenerateMD5Hash(name),
		Num:      0,
	}
	if err := global.GVA_DB.Create(&newTag).Error; err != nil {
		return nil, err
	}
	return &newTag, nil
}

func (tag *Tag) Update(obj *Tag, updates map[string]interface{}) error {
	if name, ok := updates["name"].(string); ok {
		// 如果存在，更新 identify
		updates["identify"] = utils.GenerateMD5Hash(name)
	}
	// 更新字段
	tx := global.GVA_DB.Begin()
	if err := tx.Model(obj).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (tag *Tag) DeleteByTagIds(tagIds []int) (int64, error) {
	result := global.GVA_DB.Where("id IN ?", tagIds).Unscoped().Delete(&Tag{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (tag *Tag) Count() int {
	var total int64
	global.GVA_DB.Model(&Tag{}).Count(&total)
	return int(total)
}

func (tag *Tag) FindByTagId(tagId int) (*Tag, error) {
	var obj *Tag
	err := global.GVA_DB.Where("id = ?", tagId).First(&obj).Error
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (tag *Tag) FindByTagIdArray(tagIds []int) ([]Tag, error) {
	var tags []Tag
	result := global.GVA_DB.Where("id IN ?", tagIds).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (tag *Tag) FindByName(name string) (*Tag, error) {
	var obj Tag
	err := global.GVA_DB.Where("name = ?", name).First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (tag *Tag) FindByIdentify(identify string) (*Tag, error) {
	var obj Tag
	err := global.GVA_DB.Where("identify = ?", identify).First(&obj).Error
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func (tag *Tag) FindByKeyword(keyword string, sort string) ([]*Tag, error) {
	var tags []*Tag

	db := global.GVA_DB

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 处理排序参数
	if sort != "" {
		orderType := "DESC"
		if sort[0] == '+' {
			orderType = "ASC"
			sort = sort[1:]
		} else if sort[0] == '-' {
			sort = sort[1:]
		}

		// 使用处理后的排序参数构建排序子句
		orderClause := sort + " " + orderType
		db = db.Order(orderClause)
	} else {
		// 默认按照id降序
		db = db.Order("id DESC")
	}

	err := db.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
