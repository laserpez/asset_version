package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Asset contains the SOLR SlotID for the given details
type Asset struct {
	SlotID       int
	AssetType    int
	DivisionCode string
	Environment  string
}

func main() {
	mongoSession := getDBSession("localhost")
	defer mongoSession.Close()
	coll := mongoSession.DB("test").C("assetVersions")
	insertFakeAssets(coll)

	// Initialize HTTP server
	r := gin.Default()
	// create endpoint
	r.GET("/asset_version/:divisionCode/:assetType/slot", func(c *gin.Context) {
		divisionCode := c.Params.ByName("divisionCode")
		assetType := toInt(c.Params.ByName("assetType"))
		environment := getQueryStringParameter(c.Request, "environment")
		asset := findAsset(coll, divisionCode, assetType, environment)
		c.JSON(200, asset)
	})

	// Launch server and listen on 0.0.0.0:8080
	r.Run(":8080")
}

func findAsset(coll *mgo.Collection, divisionCode string, assetType int64, environment string) *Asset {
	// FIXME: there probably is a way to make bson from an Asset
	// note: lowercase keys for now
	query := bson.M{
		"divisioncode": divisionCode,
		"assettype":    assetType,
		// TODO: actually use the environment passed in
	}

	result := Asset{}
	err := coll.Find(query).One(&result)
	if err != nil {
		// TODO: return &result, err
		panic(err)
	}
	return &result
}

func getQueryStringParameter(req *http.Request, parameter string) string {
	values := req.URL.Query()[parameter]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// getSession returns a session of mongodb running on host
func getDBSession(host string) *mgo.Session {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	return session
}

func insertFakeAssets(c *mgo.Collection) {
	assets := []Asset{
		Asset{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Truzzardi",
			Environment:  "preview",
		},
		Asset{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Prada",
			Environment:  "preview",
		},
		Asset{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Truzzardi",
			Environment:  "production",
		},
		Asset{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "VinDiesel",
			Environment:  "preview",
		},
	}
	for _, a := range assets {
		err := c.Insert(a)
		if err != nil {
			panic(err)
		}
	}

}

func toInt(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 0)
	return i
}
