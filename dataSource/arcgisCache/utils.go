package arcgisCache

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache/conf"
	"github.com/gisxiaowei/basemapServer/service"
)

// GetArcgisCache 获取缓存对象
func GetArcgisCache(path string) (ArcgisCache, error) {
	var arcgisCache ArcgisCache
	var err error
	cacheInfo, _ := getCacheInfo(path)
	arr := strings.Split(cacheInfo.Typens, "/")
	if len(arr) > 0 {
		version := arr[len(arr)-1]
		if version < "10.1" {
			err = ErrUnsupportCacheVersion
		} else if version < "10.3" {
			// 10.1，10.2
			var arcgisCache10_1 ArcgisCache10_1
			arcgisCache10_1, err = NewArcgisCache10_1(path)
			arcgisCache = &arcgisCache10_1
		} else {
			// 10.3
			var arcgisCache10_3 ArcgisCache10_3
			arcgisCache10_3, err = NewArcgisCache10_3(path)
			arcgisCache = &arcgisCache10_3
		}
	}

	return arcgisCache, err
}

// getCacheInfo 通过xml获取切片配置信息
func getCacheInfo(path string) (conf.CacheInfo, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf(`%s/conf.xml`, path))
	if err != nil {
		return conf.CacheInfo{}, err
	}
	var cacheInfo conf.CacheInfo
	err = xml.Unmarshal(content, &cacheInfo)
	if err != nil {
		return conf.CacheInfo{}, err
	}
	return cacheInfo, nil
}

// getEnvelope 通过cdi获取Envelope配置信息
func getEnvelope(path string) (conf.EnvelopeN, error) {
	content, err := ioutil.ReadFile(fmt.Sprintf(`%s/conf.cdi`, path))
	if err != nil {
		return conf.EnvelopeN{}, err
	}
	var envelope conf.EnvelopeN
	err = xml.Unmarshal(content, &envelope)
	if err != nil {
		return conf.EnvelopeN{}, err
	}
	return envelope, nil
}

// bytesToInt64 将从低位到高位存储的byte数组转为int64
func bytesToInt64(bytes []byte) int64 {
	var result int64
	for i, byte := range bytes {
		result = result | int64(byte)<<uint(i*8)
	}
	return result
}

// 获取MapServer的json字符串
func getMapServerJSONString(cacheInfo conf.CacheInfo, envelope conf.EnvelopeN, pretty bool) (string, error) {
	lods := []service.Lod{}
	for _, lodInfo := range cacheInfo.TileCacheInfo.LODInfos {
		lods = append(lods, service.Lod{
			Level:      lodInfo.LevelID,
			Resolution: lodInfo.Resolution,
			Scale:      lodInfo.Scale,
		})
	}
	mapServer := service.MapServer{
		CurrentVersion:        10.11,
		ServiceDescription:    "",
		MapName:               "Layers",
		Description:           "",
		CopyrightText:         "",
		SupportsDynamicLayers: false,
		Layers:                []interface{}{},
		Tables:                []interface{}{},
		SpatialReference: service.SpatialReference{
			Wkid:       cacheInfo.TileCacheInfo.SpatialReference.WKID,
			LatestWkid: cacheInfo.TileCacheInfo.SpatialReference.LatestWKID,
		},
		SingleFusedMapCache: true,
		TileInfo: service.TileInfo{
			Rows:               cacheInfo.TileCacheInfo.TileRows,
			Cols:               cacheInfo.TileCacheInfo.TileCols,
			Dpi:                cacheInfo.TileCacheInfo.DPI,
			Format:             cacheInfo.TileImageInfo.CacheTileFormat,
			CompressionQuality: cacheInfo.TileImageInfo.CompressionQuality,
			Origin: service.Point{
				X: cacheInfo.TileCacheInfo.TileOrigin.X,
				Y: cacheInfo.TileCacheInfo.TileOrigin.Y,
			},
			SpatialReference: service.SpatialReference{
				Wkid:       cacheInfo.TileCacheInfo.SpatialReference.WKID,
				LatestWkid: cacheInfo.TileCacheInfo.SpatialReference.LatestWKID,
			},
			Lods: lods,
		},
		InitialExtent: service.Extent{
			XMin: envelope.XMin,
			YMin: envelope.YMin,
			XMax: envelope.XMax,
			YMax: envelope.YMax,
			SpatialReference: service.SpatialReference{
				Wkid:       cacheInfo.TileCacheInfo.SpatialReference.WKID,
				LatestWkid: cacheInfo.TileCacheInfo.SpatialReference.LatestWKID,
			},
		},
		FullExtent: service.Extent{
			XMin: envelope.XMin,
			YMin: envelope.YMin,
			XMax: envelope.XMax,
			YMax: envelope.YMax,
			SpatialReference: service.SpatialReference{
				Wkid:       cacheInfo.TileCacheInfo.SpatialReference.WKID,
				LatestWkid: cacheInfo.TileCacheInfo.SpatialReference.LatestWKID,
			},
		},
		MinScale: lods[0].Scale,
		MaxScale: lods[len(lods)-1].Scale,
		Units:    "esriDecimalDegrees",
		SupportedImageFormatTypes: "PNG32,PNG24,PNG,JPG,DIB,TIFF,EMF,PS,PDF,GIF,SVG,SVGZ,BMP",
		DocumentInfo: service.DocumentInfo{
			Title:                "",
			Author:               "",
			Comments:             "",
			Subject:              "",
			Category:             "",
			AntialiasingMode:     "None",
			TextAntialiasingMode: "Force",
			Keywords:             "",
		},
		Capabilities:          "Map,Query,Data",
		SupportedQueryFormats: "JSON, AMF",
		MaxRecordCount:        1000,
		MaxImageHeight:        2048,
		MaxImageWidth:         2048,
	}

	// 转为json
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(mapServer, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(mapServer)
	}
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonBytes)

	return jsonStr, nil
}
