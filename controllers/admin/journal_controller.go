package admin

import (
	"net/http"
	"strconv"

	pongo2 "github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/models"
	"github.com/bookandmusic/docs/utils"
)

type JournalController struct{}

func NewJournalController() *JournalController {
	return &JournalController{}
}

func (journal *JournalController) Index(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/journal/admin.html", pongo2.Context{
			"csrfToken": csrf.Token(c.Request),
		})
	}
}

func (u *JournalController) JournalList(c *gin.Context) {
	var journals []*models.Journal
	journals, _ = models.NewJournal().FindAll()

	data := make([]map[string]any, 0)

	for _, journal := range journals {
		tmp := make(map[string]any)
		tmp["id"] = journal.ID
		tmp["created_at"] = journal.CreatedAt.Format("2006年01月02日")
		data = append(data, tmp)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "", "count": len(journals), "data": data})
	return
}

func (u *JournalController) DeleteJournal(c *gin.Context) {
	var journals []*models.Journal
	if err := c.ShouldBindJSON(&journals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "无效的参数"})
		return
	}
	var journalIds []int

	for _, cate := range journals {
		journalIds = append(journalIds, int(cate.ID))
	}
	num, err := models.NewJournal().DeleteByJournalIds(journalIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "删除日志失败"})
		return
	}
	if num == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "日志不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除日志成功"})
	return
}

func (j *JournalController) EditJournal(c *gin.Context) {
	journalIDStr := c.Query("id")

	var (
		journalID int
		err       error
		journal   *models.Journal
	)

	if journalIDStr == "" {
		journal = models.NewJournal()
	} else {
		journalID, err = strconv.Atoi(journalIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的日志ID"})
			return
		}
		journal, _ = models.NewJournal().FindByJournalId(journalID)
	}

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "admin/journal/edit.html", pongo2.Context{
			"csrfToken":      csrf.Token(c.Request),
			"journal":        journal,
			"journalContent": utils.EscapeJS(journal.Content),
		})
	} else if c.Request.Method == "POST" {
		content := c.PostForm("content")

		if content == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "内容是必须的"})
			return
		}

		if journal == nil || journal.ID == 0 {
			if err := journal.Create(content); err != nil {
				global.GVA_LOG.Warn("Failed to add journal", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑日志失败"})
				return
			}
		} else {
			updates := map[string]interface{}{
				"content": content,
			}
			if err := journal.Update(updates); err != nil {
				global.GVA_LOG.Warn("Failed to edit journal", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": false, "msg": "编辑日志失败"})
				return
			}
		}
		return
	}
}
