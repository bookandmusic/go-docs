package core

import (
	"os"
	"path"

	bleve "github.com/blevesearch/bleve/v2"
	_ "github.com/ttys3/gojieba-bleve/v2"
	_ "github.com/yanyiwu/gojieba"

	"github.com/bookandmusic/docs/global"
)

func InitBleveIndex() bleve.Index {
	indexPath := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "bleve.index")
	os.RemoveAll(indexPath)

	DICT_DIR := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "lib/gojieba/dict")
	DICT_PATH := path.Join(DICT_DIR, "jieba.dict.utf8")
	HMM_PATH := path.Join(DICT_DIR, "hmm_model.utf8")
	USER_DICT_PATH := path.Join(DICT_DIR, "user.dict.utf8")
	IDF_PATH := path.Join(DICT_DIR, "idf.utf8")
	STOP_WORDS_PATH := path.Join(DICT_DIR, "stop_words.utf8")

	indexMapping := bleve.NewIndexMapping()
	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":     DICT_PATH,
			"hmmpath":      HMM_PATH,
			"userdictpath": USER_DICT_PATH,
			"idf":          IDF_PATH,
			"stop_words":   STOP_WORDS_PATH,
			"type":         "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	err = indexMapping.AddCustomAnalyzer("gojieba",
		map[string]interface{}{
			"type":      "gojieba",
			"tokenizer": "gojieba",
		},
	)
	if err != nil {
		panic(err)
	}
	indexMapping.DefaultAnalyzer = "gojieba"

	index, err := bleve.New(indexPath, indexMapping)
	if err != nil {
		panic(err)
	}
	return index
}
