package models

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
)

// 标签模型
type Setting struct {
	gorm.Model
	Name  string `json:"name"`  // 键
	Value string `json:"value"` // 值
}

func NewSetting() *Setting {
	return &Setting{}
}

func (s *Setting) FindByName(name string) (*Setting, error) {
	var setting Setting
	err := global.GVA_DB.Where("name = ?", name).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (s *Setting) Create(name string, value string) error {
	setting := Setting{Name: name, Value: value}
	return global.GVA_DB.Create(&setting).Error
}

func (s *Setting) Update(updates map[string]interface{}) error {
	// 更新字段
	tx := global.GVA_DB.Begin()
	if err := tx.Model(s).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 定义一个方法，接收value和name，当name存在时，更新value，否则就创建
func (s *Setting) UpdateOrCreate(name string, value string) error {
	if setting, _ := s.FindByName(name); setting == nil {
		return s.Create(name, value)
	} else {
		return setting.Update(map[string]interface{}{"value": value})
	}
}

func (s *Setting) GetValue(name string) string {
	if setting, _ := s.FindByName(name); setting == nil {
		return ""
	} else {
		return setting.Value
	}
}
