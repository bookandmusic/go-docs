package models

import (
	"errors"

	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

// 分类模型
type Category struct {
	gorm.Model
	Name     string    `json:"name"`     // 分类名称
	Identify string    `json:"identify"` // 唯一标识
	Num      int       `json:"num"`      // 分类文章数
	Articles []Article `json:"articles"` // 分类包含的文章
}

func NewCategory() *Category {
	return &Category{}
}

func (c *Category) FindOrCreateByName(name string) (*Category, error) {
	obj, err := c.FindByName(name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if obj != nil {
		return obj, nil
	}
	return c.Create(name)
}

func (t *Category) Count() int {
	var total int64
	global.GVA_DB.Model(&Category{}).Count(&total)
	return int(total)
}

func (c *Category) FindByCategoryId(categoryId int) (*Category, error) {
	var category Category
	err := global.GVA_DB.Where("id = ?", categoryId).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *Category) DeleteByCategoryIds(categoryIds []int) (int64, error) {
	result := global.GVA_DB.Where("id IN ?", categoryIds).Unscoped().Delete(&Category{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (c *Category) FindByName(name string) (*Category, error) {
	var category Category
	err := global.GVA_DB.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (c *Category) FindByIdentify(identify string) (*Category, error) {
	var category Category
	err := global.GVA_DB.Where("identify = ?", identify).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (category *Category) Create(name string) (*Category, error) {
	newcategory := Category{
		Name:     name,
		Identify: utils.GenerateMD5Hash(name),
		Num:      0,
	}
	if err := global.GVA_DB.Create(&newcategory).Error; err != nil {
		return nil, err
	}
	return &newcategory, nil
}

func (category *Category) Update(obj *Category, updates map[string]interface{}) error {
	// 更新字段
	if name, ok := updates["name"].(string); ok {
		// 如果存在，更新 identify
		updates["identify"] = utils.GenerateMD5Hash(name)
	}
	tx := global.GVA_DB.Begin()
	if err := tx.Model(obj).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (category *Category) FindByKeyword(keyword string, sort string) ([]*Category, error) {
	var categorys []*Category

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

	err := db.Find(&categorys).Error
	if err != nil {
		return nil, err
	}
	return categorys, nil
}
