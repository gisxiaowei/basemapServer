package arcgisCache

import (
	"fmt"
	"os"
	"strings"

	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache/conf"
)

// ArcgisCache10_3 ArcGIS10.3缓存
type ArcgisCache10_3 struct {
	Path      string
	CacheInfo conf.CacheInfo
	Envelope  conf.EnvelopeN
}

// NewArcgisCache10_3 根据路径创建一个新的切片解析器
func NewArcgisCache10_3(path string) (ArcgisCache10_3, error) {
	a := ArcgisCache10_3{Path: path}
	cacheInfo, err := getCacheInfo(path)
	if err != nil {
		return a, err
	}
	a.CacheInfo = cacheInfo

	envelope, err := getEnvelope(path)
	if err != nil {
		return a, err
	}
	a.Envelope = envelope
	return a, nil
}

// GetMapServerJSONString 获取MapServer的json字符串
func (a *ArcgisCache10_3) GetMapServerJSONString(pretty bool) (string, error) {
	return getMapServerJSONString(a.CacheInfo, a.Envelope, pretty)
}

// GetTileFormat 获取瓦片格式
func (a *ArcgisCache10_3) GetTileFormat() string {
	return strings.ToLower(a.CacheInfo.TileImageInfo.CacheTileFormat)
}

// GetTileBytes 根据行列号获取切片
func (a *ArcgisCache10_3) GetTileBytes(level int64, row int64, col int64) ([]byte, error) {
	bundleFilePath, recordNumber, err := a.getTileInfo(level, row, col)
	if err != nil {
		return nil, err
	}
	imageData, err := a.getImageData(bundleFilePath, recordNumber)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

// 根据级别、行、列号获取切片信息
func (a *ArcgisCache10_3) getTileInfo(level int64, row int64, col int64) (string, int64, error) {
	packetSize := a.CacheInfo.CacheStorageInfo.PacketSize
	basePath := a.Path
	rowIndex := (row / packetSize) * packetSize
	colIndex := (col / packetSize) * packetSize

	// L：2位十进制；R：4位十六进制；C：4位十六进制
	filepath := fmt.Sprintf(`%s/_alllayers/L%02d/R%04XC%04X`, basePath, level, rowIndex, colIndex)

	// 切片顺序号
	recordNumber := packetSize*(row-rowIndex) + (col - colIndex)
	if recordNumber < 0 {
		return "", 0, ErrInvalidLevelRowCol
	}

	return filepath, recordNumber, nil
}

// 获取切片数据（bundleFilePath：文件路径，recordNumber：切片顺序号）
func (a *ArcgisCache10_3) getImageData(bundleFilePath string, recordNumber int64) ([]byte, error) {
	var result []byte

	// 打开bundle文件
	f, err := os.Open(fmt.Sprintf(`%s.bundle`, bundleFilePath))
	if err != nil {
		return result, err
	}
	defer f.Close()

	// 偏移tileOffset，找到切片位置索引
	tileOffset := 64 + (recordNumber * 8)
	_, err = f.Seek(tileOffset, 0)
	if err != nil {
		return result, err
	}

	// 读取4个字节，并转为int64，即为切片位置偏移量
	bytes := make([]byte, 4)
	_, err = f.Read(bytes)
	if err != nil {
		return result, err
	}
	imageOffset := bytesToInt64(bytes)

	// 偏移imageOffset-4，找到切片数据长度索引
	f.Seek(imageOffset-4, 0)
	if err != nil {
		return result, err
	}

	// 读取4个字节，并转为int64，即为切片数据长度
	bytes = make([]byte, 4)
	_, err = f.Read(bytes)
	if err != nil {
		return result, err
	}
	imageLength := bytesToInt64(bytes)

	// 读取imageLength字节，即为切片数据
	imageData := make([]byte, imageLength)
	f.Read(imageData)
	if err != nil {
		return result, err
	}

	return imageData, nil
}
