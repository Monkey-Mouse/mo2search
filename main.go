package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/config"
	bleveHttp "github.com/blevesearch/bleve/v2/http"
	"github.com/gin-gonic/gin"
	_ "github.com/leopku/bleve-gse-tokenizer/v2"
)

var dataDir = "./.bleve"

var Indexes map[string]bleve.Index = make(map[string]bleve.Index)

// CreateOrLoadIndex as name
func CreateOrLoadIndex(name string) {

	dir := path.Join(dataDir, name)
	var index bleve.Index
	var err error
	if _, err = os.Stat(dir); !os.IsNotExist(err) {
		index, err = bleve.Open(dir)
	} else {

		mapping := bleve.NewIndexMapping()
		if err := mapping.AddCustomTokenizer("gse", map[string]interface{}{
			"type":       "gse",
			"user_dicts": "./dict.txt", // <-- MUST specified, otherwise panic would occurred.
		}); err != nil {
			panic(err)
		}
		if err := mapping.AddCustomAnalyzer("gse", map[string]interface{}{
			"type":      "gse",
			"tokenizer": "gse",
		}); err != nil {
			panic(err)
		}
		mapping.DefaultAnalyzer = "gse"

		index, err = bleve.New(dir, mapping)
	}
	if err != nil {
		panic(err)
	}
	Indexes[name] = index
}

func main() {
	err := os.Mkdir(dataDir, os.ModeDir)
	if err != nil {
		log.Println(err.(*os.PathError))
	}
	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Fatalf("error reading data dir: %v", err)
	}

	for _, dirInfo := range dirEntries {
		indexPath := dataDir + string(os.PathSeparator) + dirInfo.Name()

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("not registering %s, skipping", indexPath)
			continue
		}
		CreateOrLoadIndex(dirInfo.Name())
		i := Indexes[dirInfo.Name()]
		if err != nil {
			log.Printf("error opening index %s: %v", indexPath, err)
		} else {
			log.Printf("registered index: %s", dirInfo.Name())
			bleveHttp.RegisterIndexName(dirInfo.Name(), i)
			// set correct name in stats
			i.SetName(dirInfo.Name())
		}
	}
	r := gin.Default()
	li := bleveHttp.NewListIndexesHandler()
	index := bleveHttp.NewDocIndexHandler("blog")
	index.DocIDLookup = func(req *http.Request) string {
		return req.URL.Query().Get("id")
	}
	index.IndexNameLookup = func(req *http.Request) string {
		_, f := path.Split(req.URL.RawPath)
		return f
	}
	search := bleveHttp.NewSearchHandler("blog")
	search.IndexNameLookup = func(req *http.Request) string {
		return req.URL.Query().Get("index")
	}
	del := bleveHttp.NewDocDeleteHandler("blog")
	del.DocIDLookup = func(req *http.Request) string {
		return req.URL.Query().Get("id")
	}
	del.IndexNameLookup = func(req *http.Request) string {
		_, f := path.Split(req.URL.RawPath)
		return f
	}
	api := r.Group("api")
	{
		api.POST("/index", func(ctx *gin.Context) {
			n := ctx.Query("name")
			if _, e := Indexes[n]; e {
				ctx.Status(http.StatusCreated)
				return
			}
			CreateOrLoadIndex(n)
			bleveHttp.RegisterIndexName(n, Indexes[n])
			ctx.Status(http.StatusCreated)
		})
		api.GET("/index", func(ctx *gin.Context) {
			li.ServeHTTP(ctx.Writer, ctx.Request)
		})
		api.PUT("/:index", func(ctx *gin.Context) {
			index.ServeHTTP(ctx.Writer, ctx.Request)
		})
		api.DELETE("/:index", func(ctx *gin.Context) {
			del.ServeHTTP(ctx.Writer, ctx.Request)
		})
		api.POST("/search", func(ctx *gin.Context) {
			search.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	r.Run(":5097")
}
