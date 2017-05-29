package arcgisCache

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	ErrInvalidLevelRowCol = errors.New("无效的级别、行、列")
)

// ArcGIS10.1缓存
type ArcgisCache10_1 struct {
	Path      string
	CacheInfo CacheInfo
}

// 根据路径创建一个新的瓦片解析器
func NewArcgisCache10_1(path string) (ArcgisCache10_1, error) {
	a := ArcgisCache10_1{Path: path}
	cacheInfo, err := a.getCacheInfo()
	if err != nil {
		return a, err
	}
	a.CacheInfo = cacheInfo
	return a, nil
}

// 根据行列号获取瓦片
func (a *ArcgisCache10_1) GetTileBytes(level int64, row int64, col int64) ([]byte, error) {
	bundleFilePath, tileOffset, err := a.getTileInfo(level, row, col)
	if err != nil {
		return nil, err
	}
	imageOffset, err := a.getImageOffset(bundleFilePath, tileOffset)
	if err != nil {
		return nil, err
	}
	imageData, err := a.getImageData(bundleFilePath, imageOffset)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

// 通过xml获取瓦片配置信息
func (a *ArcgisCache10_1) getCacheInfo() (CacheInfo, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf(`%s\conf.xml`, a.Path))
	if err != nil {
		return CacheInfo{}, err
	}
	var cacheInfo CacheInfo
	err = xml.Unmarshal(content, &cacheInfo)
	if err != nil {
		return CacheInfo{}, err
	}
	return cacheInfo, nil
}

// 根据级别、行、列号获取瓦片信息
func (a *ArcgisCache10_1) getTileInfo(level int64, row int64, col int64) (string, int64, error) {
	packetSize := a.CacheInfo.CacheStorageInfo.PacketSize
	basePath := a.Path
	rowIndex := (row / packetSize) * packetSize
	colIndex := (col / packetSize) * packetSize

	// L：2位十进制；R：4位十六进制；C：4位十六进制
	filepath := fmt.Sprintf(`%s\_alllayers\L%02d\R%04XC%04X`, basePath, level, rowIndex, colIndex)

	recordNumber := ((packetSize * (col - colIndex)) + (row - rowIndex))
	if recordNumber < 0 {
		return "", 0, ErrInvalidLevelRowCol
	}

	// bundlex：16字节头 + 81920字节（128 × 128 × 5）偏移量信息 + 16字节尾
	// 偏移量信息按列存储
	tileOffset := 16 + (recordNumber * 5)
	return filepath, tileOffset, nil
}

// 获取图片数据在bundle中的偏移量
func (a *ArcgisCache10_1) getImageOffset(bundleFilePath string, tileOffset int64) (int64, error) {
	var result int64

	// 打开bundlx文件
	f, err := os.Open(fmt.Sprintf(`%s.bundlx`, bundleFilePath))
	if err != nil {
		return result, err
	}
	defer f.Close()

	// 偏移tileOffset
	_, err = f.Seek(tileOffset, 0)
	if err != nil {
		return result, err
	}

	// 读取5个字节，并转为int64，即为图片在bundle中的偏移量
	bytes := make([]byte, 5)
	_, err = f.Read(bytes)
	if err != nil {
		return result, err
	}
	imageOffset := bytesToInt64(bytes)

	return imageOffset, nil
}

// 获取图片数据
func (a *ArcgisCache10_1) getImageData(bundleFilePath string, imageOffset int64) ([]byte, error) {
	var result []byte

	// 打开bundle文件
	f, err := os.Open(fmt.Sprintf(`%s.bundle`, bundleFilePath))
	if err != nil {
		return result, err
	}
	defer f.Close()

	// 偏移imageOffset
	f.Seek(imageOffset, 0)
	if err != nil {
		return result, err
	}

	// 读取4个字节，并转为int64，即为图片数据长度
	bytes := make([]byte, 4)
	_, err = f.Read(bytes)
	if err != nil {
		return result, err
	}
	imageLength := bytesToInt64(bytes)

	// 读取imageLength字节，即为图片数据
	imageData := make([]byte, imageLength)
	f.Read(imageData)
	if err != nil {
		return result, err
	}

	return imageData, nil
}

// 将从低位到高位存储的byte数组转为int64
func bytesToInt64(bytes []byte) int64 {
	var result int64
	for i, byte := range bytes {
		result = result | int64(byte)<<uint(i*8)
	}
	return result
}
