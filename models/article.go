package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	bleve "github.com/blevesearch/bleve/v2"
	yaml "gopkg.in/yaml.v2"
	"gorm.io/gorm"

	"github.com/bookandmusic/docs/global"
	"github.com/bookandmusic/docs/utils"
)

// 定义自定义类型
type ArticleTyle int

type ArticleTocItem struct {
	ID       string            `json:"id"`
	Title    string            `json:"title"`
	Level    int               `json:"level"`
	Children []*ArticleTocItem `json:"children"`
}

// FrontMatter 结构体用于解析Front Matter中的数据
type FrontMatter struct {
	Title      string    `yaml:"title"`
	Tags       []string  `yaml:"tags"`
	Categories []string  `yaml:"categories"`
	Date       time.Time `yaml:"date"`
}

type ArticleItem struct {
	ID         int       `json:"id"`
	Identify   string    `json:"identify"` // 唯一标识
	CategoryID int       `json:"category_id"`
	Title      string    `json:"title"`
	Order      int       `json:"order"`
	ParentID   int       `json:"parent_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type CollectionTocItem struct {
	ID       int                  `json:"id"`
	Identify string               `json:"identify"` // 唯一标识
	Title    string               `json:"title"`
	Order    int                  `json:"order"`
	ParentID int                  `json:"parent_id"`
	Children []*CollectionTocItem `json:"children"`
}

const (
	Blog ArticleTyle = 1
	Doc  ArticleTyle = 2
)

// 实现 String() 方法，根据 Category 值返回相应的字符串
func (c ArticleTyle) String() string {
	switch c {
	case Blog:
		return "blog"
	case Doc:
		return "doc"
	default:
		return "blog"
	}
}

// 实现 Parse() 方法，接受字符串表示的分类并将其解析为 Category 类型
func (c ArticleTyle) Parse(str string) ArticleTyle {
	switch str {
	case "blog":
		return Blog
	case "doc":
		return Doc
	default:
		return Blog
	}
}

type Article struct {
	gorm.Model
	Identify     string      `json:"identify"`                                      // 唯一标识
	Type         ArticleTyle `json:"type"`                                          // 文章类型，博客还是文档
	Title        string      `json:"title"`                                         // 文章标题
	Content      string      `json:"content"`                                       // 文章内容
	HtmlContent  string      `json:"html_content"`                                  // 文章 HTML 内容
	CollectionID uint        `json:"collection_id,omitempty"`                       // 文章所属文集的外键，可以为空
	Collection   Collection  `json:"collection,omitempty"`                          // 文章所属文集
	ParentID     uint        `json:"patent_id"`                                     // 文章父类id
	Order        uint        `json:"order"`                                         // 文章在文集中的顺序
	CategoryID   uint        `json:"category_id,omitempty"`                         // 文章的分类外键，可以为空
	Category     Category    `json:"category,omitempty"`                            // 文章的分类
	Tags         []Tag       `gorm:"many2many:article_tags;" json:"tags,omitempty"` // 文章的标签，可以为空
}

func NewArticle() *Article {
	return &Article{}
}

func (a *Article) Create(article *Article) error {
	tx := global.GVA_DB.Begin()
	if article.CollectionID != 0 {
		collection, err := NewCollection().FindByCollectionId(int(article.CollectionID))
		if err != nil {
			tx.Rollback()
			return err
		}
		collection.Num += 1
		if err := tx.Model(collection).Where("id = ?", collection.ID).Updates(map[string]interface{}{"num": collection.Num, "first_doc": collection.FirstDoc}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if article.Category.ID != 0 {
		article.Category.Num += 1
		if err := tx.Model(article.Category).Where("id = ?", article.Category.ID).Updates(map[string]interface{}{"num": article.Category.Num}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(article.Tags) > 0 {
		for _, tag := range article.Tags {
			tag.Num += 1
			if err := tx.Model(tag).Where("id = ?", tag.ID).Updates(map[string]interface{}{"num": tag.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if article.Identify == "" {
		article.Identify = utils.GenerateMD5Hash(time.Now().Format("20060102150405") + article.Title)
	}
	if err := tx.Create(article).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := a.BleveIndexAddArticle(article); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("bleve index add article error: %v", err))
	}
	return tx.Commit().Error
}

func (a *Article) Update(obj *Article, updates map[string]interface{}) error {
	// 更新字段
	tx := global.GVA_DB.Begin()

	newCollection, _ := updates["Collection"].(*Collection)
	if obj.Collection.ID != 0 && newCollection != nil && newCollection.ID != 0 {
		if newCollection.ID != obj.Collection.ID {
			obj.Collection.Num -= 1
			if obj.Collection.FirstDoc == obj.Identify {
				obj.Collection.FirstDoc = ""
			}
			if err := tx.Model(obj.Collection).Where("id = ?", obj.Collection.ID).Updates(map[string]interface{}{"num": obj.Collection.Num, "first_doc": obj.Collection.FirstDoc}).Error; err != nil {
				tx.Rollback()
				return err
			}
			newCollection.Num += 1
			if err := tx.Model(newCollection).Where("id = ?", newCollection.ID).Updates(map[string]interface{}{"num": newCollection.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
			obj.Collection = *newCollection
		}
	} else if obj.Collection.ID != 0 && newCollection == nil {
		obj.Collection.Num -= 1
		if err := tx.Model(obj.Collection).Where("id = ?", obj.Collection.ID).Updates(map[string]interface{}{"num": obj.Collection.Num}).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Model(obj).Association("Collection").Clear()
	} else if obj.Collection.ID == 0 && newCollection != nil && newCollection.ID != 0 {
		newCollection.Num += 1
		if err := tx.Model(newCollection).Where("id = ?", newCollection.ID).Updates(map[string]interface{}{"num": newCollection.Num}).Error; err != nil {
			tx.Rollback()
			return err
		}
		obj.Collection = *newCollection
	}

	newCategory, _ := updates["Category"].(*Category)
	if obj.Category.ID != 0 && newCategory != nil {
		if newCategory.ID != obj.Category.ID {
			obj.Category.Num -= 1
			if err := tx.Model(obj.Category).Where("id = ?", obj.Category.ID).Updates(map[string]interface{}{"num": obj.Category.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
			newCategory.Num += 1
			if err := tx.Model(newCategory).Where("id = ?", newCategory.ID).Updates(map[string]interface{}{"num": newCategory.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
			obj.Category = *newCategory
		}
	} else if obj.Collection.ID != 0 && newCategory == nil {
		obj.Category.Num -= 1
		if err := tx.Model(obj.Category).Where("id = ?", obj.Category.ID).Updates(map[string]interface{}{"num": obj.Category.Num}).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Model(obj).Association("Category").Clear()
	} else if obj.Collection.ID == 0 && newCategory != nil {
		newCategory.Num += 1
		if err := tx.Model(newCategory).Where("id = ?", newCategory.ID).Updates(map[string]interface{}{"num": newCategory.Num}).Error; err != nil {
			tx.Rollback()
			return err
		}
		obj.Category = *newCategory
	}

	tagNum := make(map[uint]int)
	if len(obj.Tags) > 0 {
		for _, tag := range obj.Tags {
			tag.Num -= 1
			if err := tx.Model(tag).Where("id = ?", tag.ID).Updates(map[string]interface{}{"num": tag.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
			tagNum[tag.ID] = tag.Num
		}
	}

	tx.Model(obj).Association("Tags").Clear()

	// 更新 Tags 字段（如果更新映射中包含 "Tags" 键）
	var (
		tags []Tag
		ok   bool
	)
	if tags, ok = updates["Tags"].([]Tag); ok {
		for _, tag := range tags {
			if val, ok := tagNum[tag.ID]; ok {
				tag.Num = val + 1
			} else {
				tag.Num += 1
			}
			if err := tx.Model(tag).Where("id = ?", tag.ID).Updates(map[string]interface{}{"num": tag.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 更新其它指定字段
	if err := tx.Model(obj).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	obj.BleveIndexRemoveArticle(obj.Identify)
	obj.BleveIndexAddArticle(&Article{
		Title:    updates["title"].(string),
		Content:  updates["content"].(string),
		Identify: obj.Identify,
	})
	return tx.Commit().Error
}

func (a *Article) DeleteByArticleIds(articleIds []int) error {
	tx := global.GVA_DB.Begin()
	for _, articleId := range articleIds {
		var article Article
		if err := tx.Preload("Category").Preload("Tags").Preload("Collection").Where("id = ?", articleId).First(&article).Error; err != nil {
			tx.Rollback()
			return err
		}
		if article.Collection.ID != 0 {
			if article.Collection.FirstDoc == article.Identify {
				article.Collection.FirstDoc = ""
			}
			article.Collection.Num -= 1
			if err := tx.Model(article.Collection).Where("id = ?", article.Collection.ID).Updates(map[string]interface{}{"num": article.Collection.Num, "first_doc": article.Collection.FirstDoc}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if article.Category.ID != 0 {
			article.Category.Num -= 1
			if err := tx.Model(article.Category).Where("id = ?", article.Category.ID).Updates(map[string]interface{}{"num": article.Category.Num}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		if len(article.Tags) > 0 {
			for _, tag := range article.Tags {
				tag.Num -= 1
				if err := tx.Model(tag).Where("id = ?", tag.ID).Updates(map[string]interface{}{"num": tag.Num}).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		tx.Model(article).Where("id = ?", articleId).Unscoped().Delete(&Article{})
		if err := a.BleveIndexRemoveArticle(article.Identify); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("remove bleve index error by article %s : %v", article.Title, err))
		}

	}
	return tx.Commit().Error
}

type ArticleScoringCriteria struct{}

func (a *Article) Search(keyword string, page, size int) []*Article {
	docs, err := a.BleveIndexSearchArticle(keyword, page, size)
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("search article by bleve index error: %v", err))
	}
	return docs
}

func (a *Article) InitArticleIndex() error {
	var (
		articles []*Article
		err      error
	)
	if articles, err = a.FindAll(); err != nil {
		return fmt.Errorf(fmt.Sprintf("find all articles error: %v", err))
	}
	for _, article := range articles {
		if err := a.BleveIndexAddArticle(article); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("bleve index add article error: %v", err))
		}
	}
	return nil
}

func (a *Article) BleveIndexAddArticle(article *Article) error {
	return global.GVA_BLEVE_INDEX.Index(article.Identify, map[string]any{
		"Content": article.Content,
		"Title":   article.Title,
	})
}

func (a *Article) BleveIndexRemoveArticle(articleIdentify string) error {
	return global.GVA_BLEVE_INDEX.Delete(articleIdentify)
}

func (a *Article) BleveIndexSearchArticle(keyword string, page, size int) ([]*Article, error) {
	from := (page - 1) * size

	query := bleve.NewQueryStringQuery(fmt.Sprintf("Content:%s Title:%s", keyword, keyword))
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = from
	searchRequest.Size = size

	highlight := bleve.NewHighlight()
	highlight.AddField("Title")
	highlight.AddField("Content")
	searchRequest.Highlight = highlight

	results := []*Article{}

	searchResult, err := global.GVA_BLEVE_INDEX.Search(searchRequest)
	if err != nil {
		return results, err
	}

	for _, hit := range searchResult.Hits {
		item := &Article{}
		item.Content = hit.Fragments["Content"][0]
		item.Title = hit.Fragments["Title"][0]
		item.Identify = hit.ID
		results = append(results, item)
	}
	return results, nil
}

func (a *Article) FindAll() ([]*Article, error) {
	var articles []*Article
	err := global.GVA_DB.Find(&articles).Error
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *Article) Count() int {
	var total int64
	global.GVA_DB.Model(&Article{}).Count(&total)
	return int(total)
}

func (a *Article) FindByArticleId(articleID int) (*Article, error) {
	// 查询文章记录
	var article Article
	if err := global.GVA_DB.Preload("Category").Preload("Collection").Preload("Tags").Where("id = ?", articleID).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (a *Article) FindByIdentify(identify string) (*Article, error) {
	// 查询文章记录
	var article Article
	if err := global.GVA_DB.Preload("Category").Preload("Tags").Where("identify = ?", identify).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (a *Article) FindBlogsByCategoryID(page, size int, categoryID int) ([]*Article, int, error) {
	var (
		articles []*Article
		total    int64
	)
	offset := (page - 1) * size

	query := global.GVA_DB.Where("type = ?", Blog).Where("category_id = ?", categoryID)

	// 获取符合条件的总数量（不包含分页条件）
	if err := query.Model(&Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id desc").Offset(offset).Limit(size).Find(&articles).Preload("Category").Preload("Tags").Error; err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil
}

// GetArticlesByTagID 根据标签名查询所有的文章
func (a *Article) FindBlogsByTagID(page, size int, tagID int) ([]*Article, int, error) {
	var (
		articles []*Article
		total    int64
	)
	offset := (page - 1) * size

	query := global.GVA_DB.Where("type = ?", Blog).
		Joins("LEFT JOIN article_tags ON articles.id = article_tags.article_id").
		Where("article_tags.tag_id = ?", tagID)

	// 获取符合条件的总数量（不包含分页条件）
	if err := query.Model(&Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id desc").Offset(offset).Limit(size).Find(&articles).Preload("Category").Preload("Tags").Error; err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}

func (a *Article) FindBlogs(page, size int) ([]*Article, int, error) {
	var (
		articles []*Article
		total    int64
	)
	offset := (page - 1) * size

	query := global.GVA_DB.Where("type = ?", Blog).Preload("Category").Preload("Tags")

	// 获取符合条件的总数量（不包含分页条件）
	if err := query.Model(&Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id desc").Offset(offset).Limit(size).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, int(total), nil
}

func (a *Article) FindBlogItems() ([]*ArticleItem, error) {
	var articles []*ArticleItem

	query := global.GVA_DB.Table("articles").Select("`id`, `category_id`, `identify`, `title`, `order`, `parent_id`, `created_at`").Where("type = ?", Blog)

	if err := query.Order("id desc").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *Article) FindArticlesByParams(params map[string]any) ([]*Article, int, error) {
	var (
		articles     []*Article
		total        int64
		page, size   int
		keyword      string
		sort         string
		articleType  ArticleTyle
		categoryID   int
		collectionID int
		ok           bool
	)
	if page, ok = params["page"].(int); !ok {
		page = 1
	}
	if size, ok = params["size"].(int); !ok {
		size = 10
	}
	if keyword, ok = params["keyword"].(string); !ok {
		keyword = ""
	}
	if sort, ok = params["sort"].(string); !ok {
		sort = ""
	}
	if articleType, ok = params["articleType"].(ArticleTyle); !ok {
		articleType = -1
	}

	if categoryID, ok = params["categoryID"].(int); !ok {
		categoryID = -1
	}

	if collectionID, ok = params["collectionID"].(int); !ok {
		collectionID = -1
	}

	offset := (page - 1) * size
	query := global.GVA_DB

	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if articleType != -1 {
		query = query.Where("type = ?", articleType)
	}

	if categoryID != -1 {
		query = query.Where("category_id = ?", categoryID)
	}

	if collectionID != -1 {
		query = query.Where("collection_id = ?", collectionID)
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
		query = query.Order(orderClause)
	} else {
		// 默认按照id降序
		query = query.Order("id DESC")
	}

	// 获取符合条件的总数量（不包含分页条件）
	if err := query.Model(&Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id desc").Offset(offset).Limit(size).Preload("Category").Preload("Tags").Preload("Collection").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, int(total), nil
}

func (a *Article) FindPreviousBlogAndNextBlog(currentArticleID int) (*ArticleItem, *ArticleItem) {
	var previousItem ArticleItem
	var nextItem ArticleItem

	// 查找上一篇文章
	if err := global.GVA_DB.Table("articles").Select("`id`, `category_id`, `identify`, `title`, `order`, `parent_id`, `created_at`").Where("type = ?", Blog).Where("id < ?", currentArticleID).
		Order("id desc").
		First(&previousItem); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("query article err: %v", err))
	}

	// 查找下一篇文章
	if err := global.GVA_DB.Table("articles").Select("`id`, `category_id`, `identify`, `title`, `order`, `parent_id`, `created_at`").Where("type = ?", Blog).Where("id > ?", currentArticleID).
		Order("id asc").
		First(&nextItem); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("query article err: %v", err))
	}

	return &previousItem, &nextItem
}

func (a *Article) FindByCollectionId(collectionID int) ([]*ArticleItem, error) {
	// 查询文章记录
	var articles []*ArticleItem
	if err := global.GVA_DB.Table("articles").Select("`id`, `category_id`, `identify`, `title`, `order`, `parent_id`, `created_at`").Where("collection_id = ?", collectionID).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *Article) TocList(articles []*ArticleItem) []*CollectionTocItem {
	// 创建一个映射，将每个文章的ID与对应的 CollectionTocItem 关联起来
	tocItemMap := make(map[int]*CollectionTocItem)

	// 创建一个用于保存综合目录的根节点
	var root CollectionTocItem
	root.Children = []*CollectionTocItem{}

	// 遍历每一篇文章并将其转换为 CollectionTocItem
	for _, article := range articles {
		// 创建文章对应的 CollectionTocItem
		articleToc := CollectionTocItem{
			ID:       article.ID,
			Title:    article.Title,
			Identify: article.Identify,
			Order:    article.Order,
			ParentID: article.ParentID, // 使用 ParentID 作为父节点ID
			Children: []*CollectionTocItem{},
		}

		// 将文章的 CollectionTocItem 添加到映射中
		tocItemMap[article.ID] = &articleToc
	}
	for _, article := range tocItemMap {
		// 如果没有父节点（ParentID 为零），直接添加到根节点的子节点中
		if article.ParentID == 0 {
			root.Children = append(root.Children, article)
		} else {
			// 如果有父节点，将文章的 CollectionTocItem 添加到父节点的子节点中
			parentToc, ok := tocItemMap[article.ParentID]
			if ok {
				parentToc.Children = append(parentToc.Children, article)
			}
		}
	}

	// 遍历整个树，对每个子节点按照 Order 字段排序
	a.sortChildrenByOrder(root.Children)

	// 返回根节点的子节点，这些子节点包含了每篇文章的 TocItem，已按照顺序排列
	return root.Children
}

func (a *Article) flattenCollectionTocItems(items []*CollectionTocItem) []*CollectionTocItem {
	flatItems := []*CollectionTocItem{}

	// 定义一个辅助函数用于递归展平元素
	var flatten func(items []*CollectionTocItem)
	flatten = func(items []*CollectionTocItem) {
		for _, item := range items {
			flatItems = append(flatItems, item)
			flatten(item.Children)
		}
	}

	flatten(items)
	return flatItems
}

func (a *Article) FindPreviousDocAndNextDoc(items []*CollectionTocItem, currentArticleID int) (*CollectionTocItem, *CollectionTocItem) {
	var previousItem CollectionTocItem
	var nextItem CollectionTocItem

	flatItems := a.flattenCollectionTocItems(items)

	for i, item := range flatItems {
		if item.ID == currentArticleID {
			if i > 0 {
				previousItem = *flatItems[i-1]
			}
			if i < len(flatItems)-1 {
				nextItem = *flatItems[i+1]
			}
			break
		}
	}

	return &previousItem, &nextItem
}

// 辅助函数：对子节点及其子孙节点按照 Order 字段排序
func (a *Article) sortChildrenByOrder(children []*CollectionTocItem) {
	sort.SliceStable(children, func(i, j int) bool {
		return (children)[i].Order < (children)[j].Order
	})

	for _, child := range children {
		a.sortChildrenByOrder(child.Children)
	}
}

func (a *Article) UpdateArticleSort(tocList []*CollectionTocItem) {
	for _, tocItem := range tocList {
		// 根据ID查找数据库中的Article
		var (
			article *Article
			err     error
		)
		if err = global.GVA_DB.Where("id = ?", tocItem.ID).First(&article).Error; err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("TocList find article %d error: %v", tocItem.ID, err))
			continue
		} else {
			// 检查是否需要更新ParentID和Order
			if int(article.ParentID) != tocItem.ParentID || int(article.Order) != tocItem.Order {
				article.ParentID = uint(tocItem.ParentID)
				article.Order = uint(tocItem.Order)
				if err = global.GVA_DB.Save(&article).Error; err != nil {
					// 处理保存错误
					global.GVA_LOG.Warn(fmt.Sprintf("TocList update article %d error: %v", tocItem.ID, err))
				}
			}
		}

		// 递归更新子节点
		a.UpdateArticleSort(tocItem.Children)
	}
}

func (a *Article) ImportMds(zipPath string, articleType ArticleTyle, collectionID, categoryID int) error {
	// 判断是否存在
	if !utils.FileExists(zipPath) {
		return errors.New("文件不存在 => " + zipPath)
	}
	w := md5.New()
	io.WriteString(w, zipPath) // 将str写入到w中
	io.WriteString(w, time.Now().String())
	md5str := fmt.Sprintf("%x", w.Sum(nil)) // w.Sum(nil)将w的hash转成[]byte格式
	SiteName := NewSetting().GetValue("site_name")
	if SiteName == "" {
		SiteName = "GDocs"
	}

	tempPath := filepath.Join(os.TempDir(), SiteName, md5str)
	docMap := make(map[string]int, 0)

	if err := os.MkdirAll(tempPath, 0o766); err != nil {
		global.GVA_LOG.Error("创建导入目录出错 => ", err)
	}
	// 如果加压缩失败
	if err := utils.Unzip(zipPath, tempPath); err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("CAll ziptil.Unzip error, zipPath: %s, tempPath: %s, err: %v",
			zipPath, tempPath, err))
		return err
	}
	// 当导入结束后，删除临时文件
	defer os.RemoveAll(tempPath)
	defer os.Remove(zipPath)

	for {
		// 如果当前目录下只有一个目录，则重置根目录
		if entries, err := os.ReadDir(tempPath); err == nil && len(entries) == 1 {
			dir := entries[0]
			if dir.IsDir() && dir.Name() != "." && dir.Name() != ".." {
				tempPath = filepath.Join(tempPath, dir.Name())
			}
			break
		} else {
			break
		}
	}
	currentTimeStr := time.Now().Format("20060102150405")

	tagsMap := map[string]Tag{}
	categoryMap := map[string]Category{}

	err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
		path = strings.Replace(path, "\\", "/", -1)
		if path == tempPath {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			// 如果是Markdown文件
			if strings.EqualFold(ext, ".md") || strings.EqualFold(ext, ".markdown") {
				var (
					htmlContent string
					markdownStr string
					frontMatter *FrontMatter
				)

				global.GVA_LOG.Info("正在处理 =>", path, info.Name())
				doc := NewArticle()
				doc.CollectionID = uint(collectionID)
				doc.CategoryID = uint(categoryID)
				docIdentify := strings.Replace(strings.TrimPrefix(path, tempPath+"/"), "/", "-", -1)
				doc.Identify = utils.GenerateMD5Hash(currentTimeStr + docIdentify)
				// 匹配图片，如果图片语法是在代码块中，这里同样会处理
				re := regexp.MustCompile(`!\[(.*?)\]\((.*?)\)`)
				markdown, err := utils.ReadFileAndIgnoreUTF8BOM(path)
				if err != nil {
					return err
				}

				// 处理图片
				markdownStr = re.ReplaceAllStringFunc(string(markdown), func(image string) string {
					images := re.FindAllSubmatch([]byte(image), -1)
					if len(images) <= 0 || len(images[0]) < 3 {
						return image
					}
					originalImageUrl := string(images[0][2])
					imageUrl := strings.Replace(string(originalImageUrl), "\\", "/", -1)

					// 如果是本地路径，则需要将图片复制到项目目录
					if !strings.HasPrefix(imageUrl, "http://") &&
						!strings.HasPrefix(imageUrl, "https://") &&
						!strings.HasPrefix(imageUrl, "ftp://") {
						// 如果路径中存在参数
						if l := strings.Index(imageUrl, "?"); l > 0 {
							imageUrl = imageUrl[:l]
						}

						if strings.HasPrefix(imageUrl, "/") {
							imageUrl = filepath.Join(tempPath, imageUrl)
						} else if strings.HasPrefix(imageUrl, "./") {
							imageUrl = filepath.Join(filepath.Dir(path), strings.TrimPrefix(imageUrl, "./"))
						} else if strings.HasPrefix(imageUrl, "../") {
							imageUrl = filepath.Join(filepath.Dir(path), imageUrl)
						} else {
							imageUrl = filepath.Join(filepath.Dir(path), imageUrl)
						}
						imageUrl = strings.Replace(imageUrl, "\\", "/", -1)
						dstFile := filepath.Join(global.GVA_CONFIG.Server.UploadPath, time.Now().Format("200601"), strings.TrimPrefix(imageUrl, tempPath))

						if utils.FileExists(imageUrl) {
							utils.CopyFile(imageUrl, dstFile)

							imageUrl = strings.TrimPrefix(strings.Replace(dstFile, "\\", "/", -1), strings.Replace(global.GVA_CONFIG.Server.WorkingDirectory, "\\", "/", -1))

							if !strings.HasPrefix(imageUrl, "/") && !strings.HasPrefix(imageUrl, "\\") {
								imageUrl = "/" + imageUrl
							}
						}

					} else {
						imageExt := utils.Md5Crypt(imageUrl) + filepath.Ext(imageUrl)

						dstFile := filepath.Join(global.GVA_CONFIG.Server.UploadPath, time.Now().Format("200601"), imageExt)

						if err := utils.DownloadAndSaveFile(imageUrl, dstFile); err == nil {
							imageUrl = strings.TrimPrefix(strings.Replace(dstFile, "\\", "/", -1), strings.Replace(global.GVA_CONFIG.Server.WorkingDirectory, "\\", "/", -1))
							if !strings.HasPrefix(imageUrl, "/") && !strings.HasPrefix(imageUrl, "\\") {
								imageUrl = "/" + imageUrl
							}
						}
					}

					imageUrl = strings.Replace(strings.TrimSuffix(image, originalImageUrl+")")+imageUrl+")", "\\", "/", -1)
					return imageUrl
				})

				if htmlContent, markdownStr, frontMatter, err = a.MarkdownToHTML(markdownStr); err != nil {
					global.GVA_LOG.Error(fmt.Sprintf("转换markdown失败: %v", err))
					return err
				}
				doc.HtmlContent = htmlContent
				doc.Content = markdownStr

				if frontMatter != nil {
					doc.Title = frontMatter.Title
					doc.CreatedAt = frontMatter.Date
					for _, category := range frontMatter.Categories {
						var (
							categoryObj Category
							ok          bool
						)
						if categoryObj, ok = categoryMap[category]; !ok {
							if obj, _ := NewCategory().FindOrCreateByName(category); obj != nil {
								categoryObj = *obj
							}
						}
						if categoryObj.ID != 0 {
							categoryMap[category] = categoryObj
							doc.Category = categoryObj
						}

					}
					for _, tag := range frontMatter.Tags {
						var (
							tagObj Tag
							ok     bool
						)
						if tagObj, ok = tagsMap[tag]; !ok {
							if obj, _ := NewTag().FindOrCreateByName(tag); obj != nil {
								tagObj = *obj
							}
						}
						if tagObj.ID != 0 {
							doc.Tags = append(doc.Tags, tagObj)
						}

					}
				} else {
					// 解析文档名称，默认使用第一个h标签为标题
					docName := strings.TrimSuffix(info.Name(), ext)
					if docName == "" {
						for _, line := range strings.Split(doc.Content, "\n") {
							if strings.HasPrefix(line, "#") {
								docName = strings.TrimLeft(line, "#")
								break
							}
						}
					}
					doc.Title = strings.TrimSpace(docName)
				}

				doc.Type = articleType
				parentId := 0

				parentIdentify := strings.Replace(strings.Trim(strings.TrimSuffix(strings.TrimPrefix(path, tempPath), info.Name()), "/"), "/", "-", -1)
				if parentIdentify != "" {
					parentIdentify = utils.GenerateMD5Hash(currentTimeStr + parentIdentify)
					if id, ok := docMap[parentIdentify]; ok {
						parentId = id
					}
				}
				if strings.EqualFold(info.Name(), "README.md") {
					global.GVA_LOG.Info(path, "|", info.Name(), "|", parentIdentify, "|", parentId)
				}
				isInsert := false
				// 如果当前文件是README.md，则将内容更新到父级
				if strings.EqualFold(info.Name(), "README.md") && parentId != 0 {
					doc.ID = uint(parentId)
					// global.GVA_LOG.Info(path,"|",parentId)
				} else {
					// global.GVA_LOG.Info(path,"|",parentIdentify)
					doc.ParentID = uint(parentId)
					isInsert = true
				}
				if err := doc.Create(doc); err != nil {
					global.GVA_LOG.Error(doc.Identify, err)
				}

				if isInsert {
					docMap[doc.Identify] = int(doc.ID)
				}
			}
		} else {
			// 如果当前目录下存在Markdown文件，则需要创建此节点
			if utils.HasFileOfExt(path, []string{".md", ".markdown"}) {
				global.GVA_LOG.Info("正在处理 =>", path, info.Name())
				identify := strings.Replace(strings.Trim(strings.TrimPrefix(path, tempPath), "/"), "/", "-", -1)
				if ok, err := regexp.MatchString(`[a-z]+[a-zA-Z0-9_.\-]*$`, identify); !ok || err != nil {
					identify = "import-" + identify
				}
				identify = utils.GenerateMD5Hash(currentTimeStr + identify)

				parentDoc := NewArticle()

				parentDoc.Identify = identify
				parentDoc.Title = "空白文档"
				parentDoc.CollectionID = uint(collectionID)
				parentDoc.CategoryID = uint(categoryID)
				parentDoc.Type = articleType

				parentId := 0

				parentIdentify := strings.TrimSuffix(identify, "-"+info.Name())

				if id, ok := docMap[parentIdentify]; ok {
					parentId = id
				}

				parentDoc.ParentID = uint(parentId)

				if err := parentDoc.Create(parentDoc); err != nil {
					global.GVA_LOG.Error(err)
				}

				docMap[identify] = int(parentDoc.ID)
				// logs.Info(path,"|",parentDoc.DocumentId,"|",identify,"|",info.Name(),"|",parentIdentify)
			}
		}

		return nil
	})
	if err != nil {
		global.GVA_LOG.Error("导入项目异常 => ", err)
	}
	global.GVA_LOG.Info("项目导入完毕 => ")

	return nil
}

func (a *Article) MarkdownToHTML(markdownStr string) (string, string, *FrontMatter, error) {
	// 定义匹配Front Matter的正则表达式
	frontMatterRegex := regexp.MustCompile(`^\S*---\n([\s\S]*?)\n---`)

	// 查找匹配的Front Matter
	var frontMatter *FrontMatter

	matches := frontMatterRegex.FindSubmatch([]byte(markdownStr))
	if len(matches) >= 2 {
		if err := yaml.Unmarshal(matches[0], &frontMatter); err != nil {
			global.GVA_LOG.Error(fmt.Sprintf("无法解析Front Matter: %v", err))
		}
		markdownStr = markdownStr[len(matches[0]):]
	}

	// 自定义解析器
	htmlContent, err := utils.MarkdownToHTML(markdownStr)
	if err != nil {
		return "", "", nil, err
	}

	return htmlContent, markdownStr, frontMatter, nil
}
