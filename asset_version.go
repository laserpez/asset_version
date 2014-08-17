package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"labix.org/v2/mgo"
)

// Asset contains the SOLR SlotID for the given details
type Asset struct {
	SlotID       int    `bson:"slotID,omitempty"`
	AssetType    int    `bson:"assetType"`
	DivisionCode string `bson:"divisionCode"`
	Environment  string `bson:"environment,omitempty"`
}

func main() {
	mongoSession := getDBSession("localhost")
	defer mongoSession.Close() // defer runs when the enclosing func (in this case, `main`) returns

	coll := mongoSession.DB("test").C("assetVersions")
	insertFakeAssets(coll)

	// Initialize HTTP server
	r := gin.Default()
	// create endpoint
	r.GET("/asset_version/:divisionCode/:assetType/slot", func(c *gin.Context) {
		divisionCode := c.Params.ByName("divisionCode")
		assetType := toInt(c.Params.ByName("assetType"))
		query := Asset{
			DivisionCode: divisionCode,
			AssetType:    assetType,
		}

		environment, err := getQueryStringParameter(c.Request, "environment")
		if err == nil {
			query.Environment = environment
		}

		fmt.Println(query)

		asset, err := findAsset(coll, &query)

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

type errorString struct {
	message string
}

func (s errorString) Error() string {
	return s.message
}

func getQueryStringParameter(req *http.Request, parameter string) (string, error) {
	values := req.URL.Query()[parameter]
	if len(values) == 0 {
		return "", errorString{message: "no such parameter"}
	}
	return values[0], nil
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

func toInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 0)
	return int(i)
}
