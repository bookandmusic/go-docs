package models

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

// 文集模型
type Collection struct {
	gorm.Model
	Name     string    `json:"name"`     // 文集名称
	Identify string    `json:"identify"` // 唯一标识
	Author   string    `json:"author"`   // 作者
	Num      int       `json:"num"`      // 文章数量
	FirstDoc string    `json:"first_doc"`
	Articles []Article `json:"articles"` // 文集包含的文章
}

func NewCollection() *Collection {
	return &Collection{}
}

func (t *Collection) ArticleCountByCollectionIds(collectionIds []int) int {
	var total int64
	global.GVA_DB.Model(&Article{}).Where("collection_id IN ?", collectionIds).Count(&total)
	return int(total)
}

func (c *Collection) Count() int {
	var total int64
	global.GVA_DB.Model(&Collection{}).Count(&total)
	return int(total)
}

func (c *Collection) Create(name, author string) (*Collection, error) {
	newCollection := Collection{
		Name:     name,
		Identify: utils.GenerateMD5Hash(name),
		Author:   author,
	}
	if err := global.GVA_DB.Create(&newCollection).Error; err != nil {
		return nil, err
	}
	return &newCollection, nil
}

func (c *Collection) Update(obj *Collection, updates map[string]interface{}) error {
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

func (c *Collection) FindByIdentify(identify string) (*Collection, error) {
	var collection Collection
	err := global.GVA_DB.Where("identify = ?", identify).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (c *Collection) FindByName(name string) (*Collection, error) {
	var collection *Collection
	err := global.GVA_DB.Where("name = ?", name).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (collection *Collection) FindByKeyword(keyword string, sort string, contains bool, valid bool) ([]*Collection, error) {
	var collections []*Collection
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

	if valid {
		db = db.Where("first_doc != ''")
	}

	err := db.Find(&collections).Error
	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (c *Collection) FindByCollectionId(collectionId int) (*Collection, error) {
	var collection Collection
	err := global.GVA_DB.Where("id = ?", collectionId).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (c *Collection) DeleteByCollectionIds(collectionIds []int) (int64, error) {
	result := global.GVA_DB.Where("id IN ?", collectionIds).Unscoped().Delete(&Collection{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
