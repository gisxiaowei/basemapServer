package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gisxiaowei/basemapServer/config"
	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache"
	"github.com/gorilla/mux"
)

var arcgisCaches = make(map[string]arcgisCache.ArcgisCache)

// 请求示例：http://localhost:9000/rest/services/SampleWorldCities10.1/MapServer/tile/0/2/2
func main() {
	var config config.Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}

	for _, s := range config.Services {
		// 创建ArcGIS缓存对象
		arcgisCache, err := arcgisCache.GetArcgisCache(s.Path)
		if err != nil {
			log.Fatal(err)
		}
		arcgisCaches[s.Name] = arcgisCache
	}

	// 路由
	r := mux.NewRouter()
	// 静态文件
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

	// {_:[/]?}表示/可以重复任意次
	r.HandleFunc("/rest/services{_:[/]?}", ServicesDirectoryHandler)
	r.HandleFunc("/rest/services/{name}/MapServer{_:[/]?}", ArcgisCacheMapServerHandler)
	r.HandleFunc("/rest/services/{name}/MapServer/tile/{level:[0-9]+}/{row:[0-9]+}/{col:[0-9]+}", ArcgisCacheTileHandler)

	// 运行
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.Server.Port), r))
}
