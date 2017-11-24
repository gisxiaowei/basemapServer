package arcgisCache

import (
	"fmt"
	"os"
	"strings"

	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache/conf"
)

// ArcgisCache10_1 ArcGIS10.1缓存
type ArcgisCache10_1 struct {
	Path      string
	CacheInfo conf.CacheInfo
	Envelope  conf.EnvelopeN
}

// NewArcgisCache10_1 根据路径创建一个新的切片解析器
func NewArcgisCache10_1(path string) (ArcgisCache10_1, error) {
	a := ArcgisCache10_1{Path: path}
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
func (a *ArcgisCache10_1) GetMapServerJSONString(pretty bool) (string, error) {
	return getMapServerJSONString(a.CacheInfo, a.Envelope, pretty)
}

// GetTileFormat 获取瓦片格式
func (a *ArcgisCache10_1) GetTileFormat() string {
	return strings.ToLower(a.CacheInfo.TileImageInfo.CacheTileFormat)
}

// GetTileBytes 根据行列号获取切片
func (a *ArcgisCache10_1) GetTileBytes(level int64, row int64, col int64) ([]byte, error) {
	bundleFilePath, recordNumber, err := a.getTileInfo(level, row, col)
	if err != nil {
		return nil, err
	}
	imageOffset, err := a.getImageOffset(bundleFilePath, recordNumber)
	if err != nil {
		return nil, err
	}
	imageData, err := a.getImageData(bundleFilePath, imageOffset)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

// 根据级别、行、列号获取切片信息
func (a *ArcgisCache10_1) getTileInfo(level int64, row int64, col int64) (string, int64, error) {
	packetSize := a.CacheInfo.CacheStorageInfo.PacketSize
	basePath := a.Path
	rowIndex := (row / packetSize) * packetSize
	colIndex := (col / packetSize) * packetSize

	// L：2位十进制；R：4位十六进制；C：4位十六进制
	filepath := fmt.Sprintf(`%s/_alllayers/L%02d/R%04XC%04X`, basePath, level, rowIndex, colIndex)

	// 切片顺序号
	recordNumber := packetSize*(col-colIndex) + (row - rowIndex)
	if recordNumber < 0 {
		return "", 0, ErrInvalidLevelRowCol
	}

	return filepath, recordNumber, nil
}

// 获取切片数据在bundle中的偏移量
func (a *ArcgisCache10_1) getImageOffset(bundleFilePath string, recordNumber int64) (int64, error) {
	var result int64

	// 打开bundlx文件
	f, err := os.Open(fmt.Sprintf(`%s.bundlx`, bundleFilePath))
	if err != nil {
		return result, err
	}
	defer f.Close()

	// bundlex：16字节头 + 81920字节（128 × 128 × 5）偏移量信息 + 16字节尾
	// 偏移tileOffset，找到记录切片位置的索引
	tileOffset := 16 + (recordNumber * 5)
	_, err = f.Seek(tileOffset, 0)
	if err != nil {
		return result, err
	}

	// 读取5个字节，并转为int64，即为切片在bundle中的偏移量
	bytes := make([]byte, 5)
	_, err = f.Read(bytes)
	if err != nil {
		return result, err
	}
	imageOffset := bytesToInt64(bytes)

	return imageOffset, nil
}

// 获取切片数据
func (a *ArcgisCache10_1) getImageData(bundleFilePath string, imageOffset int64) ([]byte, error) {
	var result []byte

	// 打开bundle文件
	f, err := os.Open(fmt.Sprintf(`%s.bundle`, bundleFilePath))
	if err != nil {
		return result, err
	}
	defer f.Close()

	// 偏移imageOffset，找到切片位置索引
	f.Seek(imageOffset, 0)
	if err != nil {
		return result, err
	}

	// 读取4个字节，并转为int64，即为切片数据长度
	bytes := make([]byte, 4)
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
