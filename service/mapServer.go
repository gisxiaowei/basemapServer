package service

type MapServer struct {
	CurrentVersion            float32          `json:"currentVersion"`
	ServiceDescription        string           `json:"serviceDescription"`
	MapName                   string           `json:"mapName"`
	Description               string           `json:"description"`
	CopyrightText             string           `json:"copyrightText"`
	SupportsDynamicLayers     bool             `json:"supportsDynamicLayers"`
	Layers                    []interface{}    `json:"layers"`
	Tables                    []interface{}    `json:"tables"`
	SpatialReference          SpatialReference `json:"spatialReference"`
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
	MaxRecordCount            int64            `json:"maxRecordCount"`
	MaxImageHeight            int64            `json:"maxImageHeight"`
	MaxImageWidth             int64            `json:"maxImageWidth"`
}

type TileInfo struct {
	Rows               int64            `json:"rows"`
	Cols               int64            `json:"cols"`
	Dpi                int64            `json:"dpi"`
	Format             string           `json:"format"`
	CompressionQuality int64            `json:"compressionQuality"`
	Origin             Point            `json:"origin"`
	SpatialReference   SpatialReference `json:"spatialReference"`
	Lods               []Lod            `json:"lods"`
}

type Point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Extent struct {
	XMin             float64          `json:"xmin"`
	YMin             float64          `json:"ymin"`
	XMax             float64          `json:"xmax"`
	YMax             float64          `json:"ymax"`
	SpatialReference SpatialReference `json:"spatialReference"`
}

type SpatialReference struct {
	Wkid       int64 `json:"wkid"`
	LatestWkid int64 `json:"latestWkid"`
}

type Lod struct {
	Level      int64   `json:"level"`
	Resolution float64 `json:"resolution"`
	Scale      int64   `json:"scale"`
}

type DocumentInfo struct {
	Title                string `json:"title"`
	Author               string `json:"author"`
	Comments             string `json:"comments"`
	Subject              string `json:"subject"`
	Category             string `json:"category"`
	AntialiasingMode     string `json:"antialiasingMode"`
	TextAntialiasingMode string `json:"textAntialiasingMode"`
	Keywords             string `json:"keywords"`
}
