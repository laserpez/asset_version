package main

import (
	"fmt"
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
	//////////////////////////////////////////////////////////////////////////////
	// Initialize Mongo
	//////////////////////////////////////////////////////////////////////////////
	mongoHost := "localhost"
	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	// close session after everything has been done on it
	defer session.Close()

	// get the right collection in the right DB
	coll := session.DB("test").C("assetVersions")
	session.SetMode(mgo.Monotonic, true)
	seed(coll)

	//////////////////////////////////////////////////////////////////////////////
	// Initialize HTTP server
	//////////////////////////////////////////////////////////////////////////////
	r := gin.Default()
	// create endpoint
	r.GET("/asset_version/:divisionCode/:assetType", func(c *gin.Context) {

		// parse parameters
		divisionCode := c.Params.ByName("divisionCode")
		assetType := toInt(c.Params.ByName("assetType"))
		// TODO: parse the optional GET parameter "environment"

		// form the mongo query
		// note: lowercase keys
		query := bson.M{
			"divisioncode": divisionCode,
			"assettype":    assetType,
		}
		fmt.Println(query)

		// find the result
		result := Asset{}
		err = coll.Find(query).One(&result)
		if err != nil {
			// TODO: return a 404 here if err == not found
			panic(err)
		}

		// serialize and serve
		c.JSON(200, result)
	})

	// Launch server and listen on 0.0.0.0:8080
	r.Run(":8080")
}

func seed(c *mgo.Collection) {
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
	for a := range assets {
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
