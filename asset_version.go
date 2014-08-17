package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
)

const (
	databaseAddress string = "localhost"
	databaseName    string = "test"
	collectionName  string = "assetVersions"
)

// Asset contains the SOLR SlotID for the given details
type Asset struct {
	SlotID       int    `bson:"slotID,omitempty"`
	AssetType    int    `bson:"assetType"`
	DivisionCode string `bson:"divisionCode"`
	Environment  string `bson:"environment,omitempty"`
}

func main() {
	mgoSession, err := mgo.Dial(databaseAddress)
	if err != nil {
		panic(err)
	}
	defer mgoSession.Close()
	// put some data in there, just to play around
	insertFakeAssets(mgoSession.DB(databaseName).C(collectionName))

	// Initialize HTTP server
	r := gin.Default()
	// create endpoint
	r.GET("/asset_version/:divisionCode/:assetType/slot", func(c *gin.Context) {
		divisionCode := c.Params.ByName("divisionCode")
		assetType := toInt(c.Params.ByName("assetType"))
		environment := getQueryStringParameter(c.Request, "environment")

		query := &Asset{
			DivisionCode: divisionCode,
			AssetType:    assetType,
			Environment:  environment,
		}

		session := mgoSession.Copy()
		defer session.Close() // defer runs when the enclosing func returns
		coll := session.DB(databaseName).C(collectionName)

		asset, err := findAsset(coll, query)

		if err != nil {
			c.String(404, err.Error())
		} else {
			c.JSON(200, asset)
		}

	})

	// Launch server and listen on 0.0.0.0:8080
	r.Run(":8080")
}

func findAsset(coll *mgo.Collection, query *Asset) (*Asset, error) {
	result := Asset{}
	err := coll.Find(query).One(&result)
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func getQueryStringParameter(req *http.Request, parameter string) string {
	values := req.URL.Query()[parameter]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func insertFakeAssets(c *mgo.Collection) {
	assets := []Asset{
		{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Truzzardi",
			Environment:  "preview",
		},
		{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Prada",
			Environment:  "preview",
		},
		{
			SlotID:       1,
			AssetType:    2,
			DivisionCode: "Truzzardi",
			Environment:  "production",
		},
		{
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

func toInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 0)
	return int(i)
}
