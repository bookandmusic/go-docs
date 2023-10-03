package core

import (
	"path"

	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"

	"github.com/bookandmusic/docs/global"
)

// InitSearchEngine 初始化 Wukong 搜索引擎
func InitSearchEngine() *engine.Engine {
	dictFile := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "index/data/dictionary.txt")
	stopTokenFile := path.Join(global.GVA_CONFIG.Server.WorkingDirectory, "index/data/stop_tokens.txt")
	// 创建搜索引擎实例
	searcher := engine.Engine{}

	// 初始化搜索引擎
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: dictFile,
		StopTokenFile:         stopTokenFile,
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.LocationsIndex,
		},
	})

	return &searcher
}
