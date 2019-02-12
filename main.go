package main

//! Golang(API) =>Echo + Mondgo Atlas

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//* Model
type Profile struct {
	Username     string	   `json:"username" bson:"username"`
	Password     string    	   `json:"password" bson:"password"`
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Avatar       string        `json:"avatar" bson:"avatar"`
	Name         string        `json:"name" bson:"name"`
	Nameth       string        `json:"nameth" bson:"nameth"`
	Nickname     string        `json:"nickname" bson:"nickname"`
	Nicknameth   string        `json:"nicknameth" bson:"nicknameth"`
	Telno        string        `json:"telno" bson:"telno"`
	Email        string        `json:"email" bson:"email"`
	Team         string        `json:"team" bson:"team"`
	Company      string        `json:"company" bson:"company"`
	Oslogo       string        `json:"oslogo" bson:"oslogo"`
	Mobileoslogo string        `json:"mobileoslogo" bson:"mobileoslogo"`
	Locate       string        `json:"locate" bson:"locate"`
	Time         string        `json:"time" bson:"time"`
}

const (
	mongo_host = "mongodb://admin:muyon@teamming-shard-00-00-odfpd.mongodb.net:27017,teamming-shard-00-01-odfpd.mongodb.net:27017,teamming-shard-00-02-odfpd.mongodb.net:27017/test?&replicaSet=Teamming-shard-0&authSource=admin"
)

func main() {
	//+ Echo instance
	e := echo.New()

	//+ Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	//+ Route =>handler
	//* Hi!
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to api teammate")
	})
	//* Query all data
	e.GET("/read", Getdata)

	//+ Start server
	e.Logger.Fatal(e.Start(getPort()))
}
func getPort() string {
	var port = os.Getenv("PORT") // ----> (A)
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port // ----> (B)
}

func Getdata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongo_host)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()

	s := session.DB("teamming").C("detail")
	var profiles []Profile
	err = s.Find(bson.M{}).Limit(100).All(&profiles)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}

	return c.JSON(http.StatusOK, profiles)

}
