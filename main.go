package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gisxiaowei/basemapServer/config"
	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache/arcgisCache10_1"
	"github.com/gisxiaowei/basemapServer/service"
	"github.com/gorilla/mux"
)

var arcgisCache10_1s = make(map[string]arcgisCache10_1.ArcgisCache10_1)

// 请求示例：http://localhost:9000/rest/services/USA/MapServer/tile/2/34/24
func main() {
	var services config.Services
	if _, err := toml.DecodeFile("config.toml", &services); err != nil {
		log.Fatal(err)
	}

	for _, s := range services.Service {
		// 创建ArcGIS缓存对象
		arcgisCache10_1, err := arcgisCache10_1.NewArcgisCache10_1(s.Path)
		if err != nil {
			log.Fatal(err)
		}
		arcgisCache10_1s[s.Name] = arcgisCache10_1
	}

	// 路由
	r := mux.NewRouter()
	// {_:[/]?}表示/可以重复任意次
	r.HandleFunc("/rest/services/{name}/MapServer{_:[/]?}", ArcgisCache10_1Handler)
	r.HandleFunc("/rest/services/{name}/MapServer/tile/{level:[0-9]+}/{row:[0-9]+}/{col:[0-9]+}", ArcgisCache10_1TileHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	// 运行
	log.Fatal(http.ListenAndServe(":9000", r))
}

// MapServer处理函数
func ArcgisCache10_1Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCache10_1s[name]; ok {
		arcgisCache10_1 := arcgisCache10_1s[name]
		cacheInfo := arcgisCache10_1.CacheInfo
		envelope := arcgisCache10_1.Envelope
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

		//w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query()

		// format
		f := query.Get("f")
		var jsonBytes []byte
		var err error
		if strings.ToLower(f) == "pjson" {
			jsonBytes, err = json.MarshalIndent(mapServer, "", "  ")
		} else {
			jsonBytes, err = json.Marshal(mapServer)
		}
		if err != nil {
			log.Fatal(err)
		}
		jsonStr := string(jsonBytes)

		// callback
		callback := query.Get("callback")
		callback = strings.TrimSpace(callback)
		if callback != "" {
			jsonStr = fmt.Sprintf(`%s(%s);`, callback, jsonStr)
		}

		w.Write([]byte(jsonStr))
	} else {
		http.NotFound(w, r)
	}
}

// 瓦片处理函数
func ArcgisCache10_1TileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCache10_1s[name]; ok {
		arcgisCache10_1 := arcgisCache10_1s[name]
		// 级别、行、列号
		level, _ := strconv.ParseInt(vars["level"], 10, 64)
		row, _ := strconv.ParseInt(vars["row"], 10, 64)
		col, _ := strconv.ParseInt(vars["col"], 10, 64)
		bytes, _ := arcgisCache10_1.GetTileBytes(level, row, col)

		// 图片格式
		var suffix string
		if strings.Contains(strings.ToUpper(arcgisCache10_1.CacheInfo.TileImageInfo.CacheTileFormat), "PNG") {
			suffix = "png"
		} else {
			suffix = "jpg"
		}

		w.Header().Set("Content-Type", "image/"+suffix)
		w.Write(bytes)
	} else {
		http.NotFound(w, r)
	}
}
