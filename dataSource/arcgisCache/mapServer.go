package arcgisCache

type MapServer struct {
	CurrentVersion            float32          `json:"currentVersion"`
	ServiceDescription        string           `json:"serviceDescription"`
	MapName                   string           `json:"mapName"`
	Description               string           `json:"description"`
	CopyrightText             string           `json:"copyrightText"`
	SupportsDynamicLayers     bool             `json:"supportsDynamicLayers"`
	Layers                    []interface{}    `json:"layers"`
	Tables                    []interface{}    `json:"tables"`
	spatialReference          SpatialReference `json:"spatialReference"`
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
	MaxImageWidth             int32            `json:"maxImageWidth"`
}

type TileInfo struct {
	Rows               int32            `json:"rows"`
	Cols               int32            `json:"cols"`
	Dpi                int32            `json:"dpi"`
	Format             string           `json:"format"`
	CompressionQuality int32            `json:"compressionQuality"`
	Origin             Point            `json:"origin"`
	SpatialReference   SpatialReference `json:"spatialReference"`
	Lods               []Lod            `json:"lods"`
}

type Point struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Extent struct {
	Xmin             float64          `json:"xmin"`
	Ymin             float64          `json:"ymin"`
	Xmax             float64          `json:"xmax"`
	Ymax             float64          `json:"ymax"`
	SpatialReference SpatialReference2 `json:"spatialReference"`
}

type SpatialReference2 struct {
	Wkid       int32 `json:"wkid"`
	LatestWkid int32 `json:"latestWkid"`
}

type Lod struct {
	Level      int32   `json:"level"`
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
