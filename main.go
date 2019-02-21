package main

//! Golang(API) =>Echo + Mondgo Atlas

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	Date         string        `json:"date" bson:"date"`
}

const (
	// mongoHost = "mongodb://127.0.0.1:27017"
	mongoHost = "mongodb://admin:muyon@teamming-shard-00-00-odfpd.mongodb.net:27017,teamming-shard-00-01-odfpd.mongodb.net:27017,teamming-shard-00-02-odfpd.mongodb.net:27017/test?&replicaSet=Teamming-shard-0&authSource=admin"
)

const (
	dbname     = "teammate"
	collection = "profile"
	// dbname     = "user"
	// collection = "users"
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

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
	e.POST("/register", Postdata)
	e.PUT("/updateavatar", UploadImage)
	e.GET("/image/:username", GetImage)
	e.PUT("/checkin", UpdateTimeAndLocation)
	e.PUT("/telno", UpdateTelNumber)
	e.PUT("/email", UpdateEmail)
	e.PUT("/team", UpdateTeam)
	e.GET("/sort", SortDateAndTime)
	e.PUT("/forgetpassword", RandomCode)

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

	var profiles []Profile
	err = session.DB(dbname).C(collection).Find(bson.M{}).Sort("name").All(&profiles)
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
	profiles.ID = bson.NewObjectId()
	err = c.Bind(&profiles) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Insert(profiles) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Insert mongo.")
	}
	return c.JSON(http.StatusCreated, "Post successfully!") //* Done!
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
	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	dataimg := imgprofile.Avatar
	bs64tostr, _ := base64.StdEncoding.DecodeString(string(dataimg))
	img, _, err := image.Decode(bytes.NewReader(bs64tostr))
	if err != nil {
		m := imaging.Thumbnail(img, 500, 500, imaging.Lanczos)
		buf := new(bytes.Buffer)
		err = png.Encode(buf, m)
		imageBit := buf.Bytes()
		photoBase64 := base64.StdEncoding.EncodeToString(imageBit)
		//fmt.Println(photoBase64)
		if imgprofile.Username == "" { //* Forbid blank name and telno.
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
		}
		err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"avatar": string(photoBase64)}}) //* Choose database, collection and insert data
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error! <=from Update mongo.")
		}
		return c.JSON(http.StatusCreated, "Update successfully!")
	}
	m := imaging.Thumbnail(img, 500, 500, imaging.Lanczos)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	imageBit := buf.Bytes()
	photoBase64 := base64.StdEncoding.EncodeToString(imageBit)
	//fmt.Println(photoBase64)
	if imgprofile.Username == "" { //* Forbid blank name and telno.
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"avatar": string(photoBase64)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// GetImage using for get image in mongo atlas
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
	username := c.Param("username")
	var imagepro Profile
	err = session.DB(dbname).C(collection).Find(bson.M{"username": string(username)}).Select(bson.M{"avatar": 1}).One(&imagepro) //* Delve all data
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}
	img2html := "<html><body><img src=\"data:image/png;base64," + imagepro.Avatar + "\" /></body></html>"
	return c.HTML(http.StatusOK, img2html) //* Done (return data to API)
}

// UpdateTimeAndLocation using for update time and location in mongo atlas
func UpdateTimeAndLocation(c echo.Context) (err error) {
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

	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"time": string(imgprofile.Time)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Time in mongo.")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"locate": string(imgprofile.Locate)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Location in mongo.")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"date": string(imgprofile.Date)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Date in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// UpdateTelNumber using for update TelNumber in mongo atlas
func UpdateTelNumber(c echo.Context) (err error) {
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

	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"telno": string(imgprofile.Telno)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update telno in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// UpdateEmail using for update Email in mongo atlas
func UpdateEmail(c echo.Context) (err error) {
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

	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"email": string(imgprofile.Email)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Email in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// UpdateTeam using for update Team in mongo atlas
func UpdateTeam(c echo.Context) (err error) {
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

	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"team": string(imgprofile.Team)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Team in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// UpdatePassword using for update password in mongo atlas
func UpdatePassword(c echo.Context) (err error) {
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

	var imgprofile Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(imgprofile.Username)}, bson.M{"$set": bson.M{"password": string(imgprofile.Password)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Password in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!
}

// SortDateAndTime using for get date and time in mongo atlas
func SortDateAndTime(c echo.Context) (err error) {
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

	var profiles []Profile
	err = session.DB(dbname).C(collection).Find(bson.M{}).Sort("-date", "-time").All(&profiles)
	if err != nil {
		fmt.Println("Error query mongo:", err)
	}

	return c.JSON(http.StatusOK, "Update successfully!")

}

func RandomCode(c echo.Context) (err error) {
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
	rancode := make([]byte, 8)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := 8-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			rancode[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	err = session.DB(dbname).C(collection).Update(bson.M{"username": string(profiles.Username)}, bson.M{"$set": bson.M{"password": string(rancode)}}) //* Choose database, collection and insert data
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error! <=from Update Password in mongo.")
	}
	return c.JSON(http.StatusCreated, "Update successfully!") //* Done!

}
