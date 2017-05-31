package main

import (
	"github.com/BurntSushi/toml"
	"github.com/gisxiaowei/basemapServer/config"
	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var arcgisCache10_1s map[string]arcgisCache.ArcgisCache10_1 = make(map[string]arcgisCache.ArcgisCache10_1)

// 请求示例：http://localhost:9000/2/36/28
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
	r.HandleFunc("/services/{name}/{level:[0-9]+}/{row:[0-9]+}/{col:[0-9]+}", ArcgisCache10_1Handler)

	// 运行
	log.Fatal(http.ListenAndServe(":9000", r))
}

func ArcgisCache10_1Handler(w http.ResponseWriter, r *http.Request) {
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
