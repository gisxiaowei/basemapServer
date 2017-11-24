package arcgisCache

import (
	"errors"
)

var (
	ErrUnsupportCacheVersion = errors.New("不支持的缓存版本")
	ErrInvalidLevelRowCol    = errors.New("无效的级别、行、列")
)

// ArcgisCache ArcGIS缓存接口
type ArcgisCache interface {
	GetMapServerJSONString(pretty bool) (string, error)
	GetTileFormat() string
	GetTileBytes(level int64, row int64, col int64) ([]byte, error)
}
