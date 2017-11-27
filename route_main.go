package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gisxiaowei/basemapServer/dataSource/arcgisCache"
	"github.com/gisxiaowei/basemapServer/service"
	"github.com/gorilla/mux"
)

// RootHandler 根目录处理函数
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/rest/services", http.StatusFound)
}

// ServicesDirectoryHandler 服务目录处理函数
func ServicesDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	// query
	query := r.URL.Query()

	// format
	f := strings.TrimSpace(strings.ToLower(query.Get("f")))
	if f == "" || f == "html" { // html
		templates := template.Must(template.ParseFiles("templates/servicesDirectory.html"))
		err := templates.ExecuteTemplate(w, "servicesDirectory", arcgisCaches)
		if err != nil {
			log.Fatalln("模板出错", err)
		}
	} else if f == "json" || f == "pjson" { // json
		pretty := f == "pjson"
		jsonStr, err := getServicesDirectoryJSONString(arcgisCaches, pretty)
		if err != nil {
			log.Fatal(err)
		}

		// callback
		callback := query.Get("callback")
		callback = strings.TrimSpace(callback)
		if callback != "" {
			jsonStr = fmt.Sprintf(`%s(%s);`, callback, jsonStr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonStr))
	} else {
		w.WriteHeader(400)
		templates := template.Must(template.ParseFiles("templates/error.html"))
		err := templates.ExecuteTemplate(w, "error", service.Error{Message: "不支持此格式", Code: 400})
		if err != nil {
			log.Fatalln("模板出错", err)
		}
	}
}

// 获取服务目录对象json字符串
func getServicesDirectoryJSONString(arcgisCaches map[string]arcgisCache.ArcgisCache, pretty bool) (string, error) {
	services := []service.Service{}
	for key := range arcgisCaches {
		services = append(services, service.Service{
			Name: key,
			Type: "MapServer",
		})
	}
	servicesDirectory := service.ServicesDirectory{
		CurrentVersion: 10.11,
		Folders:        []interface{}{},
		Services:       services,
	}

	// 转为json
	var jsonBytes []byte
	var err error
	if pretty {
		jsonBytes, err = json.MarshalIndent(servicesDirectory, "", "  ")
	} else {
		jsonBytes, err = json.Marshal(servicesDirectory)
	}
	if err != nil {
		return "", err
	}
	jsonStr := string(jsonBytes)

	return jsonStr, nil
}

// ArcgisCacheMapServerHandler MapServer处理函数
func ArcgisCacheMapServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// 服务名
	name, _ := vars["name"]
	if _, ok := arcgisCaches[name]; ok {
		arcgisCache := arcgisCaches[name]

		// query
		query := r.URL.Query()

		// format
		f := strings.TrimSpace(strings.ToLower(query.Get("f")))
		if f == "" || f == "html" { // html
			templates := template.Must(template.ParseFiles("templates/mapServer.html"))
			err := templates.ExecuteTemplate(w, "mapServer", name)
			if err != nil {
				log.Fatalln("模板出错", err)
			}
		} else if f == "json" || f == "pjson" { // json
			pretty := f == "pjson"
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

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(jsonStr))
		} else if f == "jsapi" { // jsapi
			templates := template.Must(template.ParseFiles("templates/jsapi.html"))
			err := templates.ExecuteTemplate(w, "jsapi", name)
			if err != nil {
				log.Fatalln("模板出错", err)
			}
		} else {
			w.WriteHeader(400)
			templates := template.Must(template.ParseFiles("templates/error.html"))
			err := templates.ExecuteTemplate(w, "error", service.Error{Message: "不支持此格式", Code: 400})
			if err != nil {
				log.Fatalln("模板出错", err)
			}
		}

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
