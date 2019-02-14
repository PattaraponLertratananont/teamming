package main

//! Golang(API) =>Echo + Mondgo Atlas

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nfnt/resize"
)

// Profile is struct to get detail form data
type Profile struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Username     string        `json:"username" bson:"username"`
	Password     string        `json:"password" bson:"password"`
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

// ImageProfile is struct to get image form data
type ImageProfile struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username,omitempty"`
	Avatar   string        `json:"avatar" bson:"avatar,omitempty"`
}

const (
	mongoHost = "mongodb://admin:muyon@teamming-shard-00-00-odfpd.mongodb.net:27017,teamming-shard-00-01-odfpd.mongodb.net:27017,teamming-shard-00-02-odfpd.mongodb.net:27017/test?&replicaSet=Teamming-shard-0&authSource=admin"
)

const (
	dbname     = "teammate"
	collection = "detail"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to api teammate")
	})
	// Query all data
	e.GET("/read", Getdata)
	e.POST("/post", Postdata)
	e.GET("/image/:username", GetImage)
	e.POST("/image/upload", UploadImage)

	// Start server
	e.Logger.Fatal(e.Start(getPort()))
}
func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "1323"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port
}

// Getdata using for get data form mongo atlas
func Getdata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongoHost)

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

// Postdata using for post data to mongo atlas
func Postdata(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongoHost)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()

	var profiles Profile
	err = c.Bind(&profiles) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	if profiles.Name == "" || profiles.Telno == "" { //* Forbid blank name and telno.
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}

	err = session.DB(dbname).C(collection).Insert(profiles) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Insert mongo.")
	}
	return c.JSON(http.StatusCreated, "Post successfully!") //* Done!
}

// GetImage using for get image form mongo atlas
func GetImage(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongoHost)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()
	//image get
	username := c.Param("username")
	//img, err := h.FindByUsername(username)
	fmt.Println(username)
	var imagepro *ImageProfile
	s := session.DB(dbname).C("image")                                                            //* Choose Database and Collection
	err = s.Find(bson.M{"username": string(username)}).Select(bson.M{"avatar": 1}).One(&imagepro) //* Delve all data
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}
	img2html := "<html><body><img src=\"data:image/png;base64," + imagepro.Avatar + "\" /></body></html>"
	return c.HTML(http.StatusOK, img2html) //* Done (return data to API)
}

// UploadImage using for post image to mongo atlas
func UploadImage(c echo.Context) (err error) {
	tlsConfig := &tls.Config{}
	dialInfo, err := mgo.ParseURL(mongoHost)

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		fmt.Println("Can't connect database:", err)
		return c.String(http.StatusInternalServerError, "Oh!, Can't connect database.")
	}
	defer session.Close()
	// image upload
	var imgprofile Profile
	imgprofile.Username = c.FormValue("username")
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	reader := bufio.NewReader(src)
	content, _ := ioutil.ReadAll(reader)
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		imageResized := resize.Resize(250, 250, img, resize.Lanczos3)
		buf := new(bytes.Buffer)
		err = png.Encode(buf, imageResized)
		imageBit := buf.Bytes()
		photoBase64 := base64.StdEncoding.EncodeToString(imageBit)
		fmt.Println(photoBase64)
		imgprofile.Avatar = photoBase64
		if imgprofile.Username == "" { //* Forbid blank name and telno.
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
		}
		err = session.DB(dbname).C("image").Update(bson.M{"username": imgprofile.Username}, bson.M{"avatar": imgprofile.Avatar}) //* Choose database, collection and insert data
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error! <=from Insert mongo.")
		}
		return c.JSON(http.StatusCreated, imgprofile)
	}
	imageResized := resize.Resize(250, 250, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, imageResized, nil)
	imageBit := buf.Bytes()
	photoBase64 := base64.StdEncoding.EncodeToString(imageBit)
	fmt.Println(photoBase64)
	imgprofile.Avatar = photoBase64
	if imgprofile.Username == "" { //* Forbid blank name and telno.
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}
	err = session.DB(dbname).C("image").Insert(imgprofile) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Insert mongo.")
	}
	return c.JSON(http.StatusCreated, imgprofile) //* Done!
}
