package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gisxiaowei/basemapServer/config"
	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache"
	"github.com/gorilla/mux"
)

var arcgisCache10_1s map[string]arcgisCache.ArcgisCache10_1 = make(map[string]arcgisCache.ArcgisCache10_1)

// 请求示例：http://localhost:9000/rest/services/USA/MapServer/tile/2/34/24
func main() {
	var services config.Services
	if _, err := toml.DecodeFile("config.toml", &services); err != nil {
		log.Fatal(err)
	}

	for _, s := range services.Service {
		// 创建ArcGIS缓存对象
		arcgisCache10_1, err := arcgisCache.NewArcgisCache10_1(s.Path)
		if err != nil {
			log.Fatal(err)
		}
		arcgisCache10_1s[s.Name] = arcgisCache10_1
	}

	// 路由
	r := mux.NewRouter()
	r.HandleFunc("/rest/services/{name}/MapServer", ArcgisCache10_1Handler)
	r.HandleFunc("/rest/services/{name}/MapServer/tile/{level:[0-9]+}/{row:[0-9]+}/{col:[0-9]+}", ArcgisCache10_1TileHandler)

	// 运行
	log.Fatal(http.ListenAndServe(":9000", r))
}

// MapServer处理函数
func ArcgisCache10_1Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCache10_1s[name]; ok {
		//arcgisCache10_1 := arcgisCache10_1s[name]
		mapServer := arcgisCache.MapServer{
			CurrentVersion:        10.11,
			ServiceDescription:    "",
			MapName:               "Layers",
			Description:           "",
			CopyrightText:         "",
			SupportsDynamicLayers: false,
			Layers:                nil,
			Tables:                nil,
			/*spatialReference          SpatialReference `json:"spatialReference"`
			SingleFusedMapCache       bool             `json:"singleFusedMapCache"`
			TileInfo                  TileInfo         `json:"tileInfo"`
			InitialExtent             Extent           `json:"initialExtent"`
			FullExtent                Extent           `json:"fullExtent"`
			MinScale                  int64            `json:"minScale"`
			MaxScale                  int64            `json:"maxScale"`
			Units                     string           `json:"units"`
			SupportedImageFormatTypes string           `json:"supportedImageFormatTypes"`
			DocumentInfo              DocumentInfo     `json:"documentInfo"`
			Capabilities              string           `json:"capabilities"`
			SupportedQueryFormats     string           `json:"supportedQueryFormats"`
			MaxRecordCount            int32            `json:"maxRecordCount"`
			MaxImageHeight            int32            `json:"maxImageHeight"`
			MaxImageWidth             int32            `json:"maxImageWidth"`*/
		}

		w.Header().Set("Content-Type", "text/plain")
		jsonStr, _ := json.Marshal(mapServer)
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
