package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ServicesDirectoryHandler 服务目录处理函数
func ServicesDirectoryHandler(w http.ResponseWriter, r *http.Request) {

	templates := template.Must(template.ParseFiles("templates/servicesDirectory.html"))
	err := templates.ExecuteTemplate(w, "servicesDirectory", arcgisCaches)
	if err != nil {
		log.Fatalln("模板出错", err)
	}
}

// ArcgisCacheMapServerHandler MapServer处理函数
func ArcgisCacheMapServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCaches[name]; ok {
		arcgisCache := arcgisCaches[name]

		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query()

		// format
		f := query.Get("f")
		pretty := strings.ToLower(f) == "pjson"
		jsonStr, err := arcgisCache.GetMapServerJSONString(pretty)
		if err != nil {
			log.Fatal(err)
		}

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

// ArcgisCacheTileHandler 瓦片处理函数
func ArcgisCacheTileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCaches[name]; ok {
		arcgisCache := arcgisCaches[name]
		// 级别、行、列号
		level, _ := strconv.ParseInt(vars["level"], 10, 64)
		row, _ := strconv.ParseInt(vars["row"], 10, 64)
		col, _ := strconv.ParseInt(vars["col"], 10, 64)
		bytes, _ := arcgisCache.GetTileBytes(level, row, col)

		suffix := arcgisCache.GetTileFormat()
		w.Header().Set("Content-Type", "image/"+suffix)
		w.Write(bytes)
	} else {
		http.NotFound(w, r)
	}
}
