package models

import (
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

// 标签模型
type Journal struct {
	gorm.Model
	Content     string `json:"content"`
	HtmlContent string `json:"html_content"`
}

func NewJournal() *Journal {
	return &Journal{}
}

func (journal *Journal) Create(content string) error {
	htmlContent, _ := utils.MarkdownToHTML(content)
	newJournal := Journal{
		Content:     content,
		HtmlContent: htmlContent,
	}
	if err := global.GVA_DB.Create(&newJournal).Error; err != nil {
		return err
	}
	return nil
}

func (journal *Journal) Update(updates map[string]interface{}) error {
	if content, ok := updates["content"].(string); ok {
		htmlContent, _ := utils.MarkdownToHTML(content)
		// 如果存在，更新 identify
		updates["html_content"] = htmlContent
	}
	// 更新字段
	tx := global.GVA_DB.Begin()
	if err := tx.Model(journal).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *Journal) DeleteByJournalIds(journalIds []int) (int64, error) {
	result := global.GVA_DB.Where("id IN ?", journalIds).Unscoped().Delete(&Journal{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (journal *Journal) FindAll() ([]*Journal, error) {
	var journals []*Journal
	err := global.GVA_DB.Order("id asc").Find(&journals).Error
	if err != nil {
		return nil, err
	}
	return journals, nil
}

func (t *Journal) FindByJournalId(journalId int) (*Journal, error) {
	var journal Journal
	err := global.GVA_DB.Where("id = ?", journalId).First(&journal).Error
	if err != nil {
		return nil, err
	}
	return &journal, nil
}
