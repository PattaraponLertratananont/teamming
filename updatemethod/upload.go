package updatemethod

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/kachawut/api_test/teamming/model"
	"github.com/labstack/echo"
)

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
	var imgprofile model.Profile
	err = c.Bind(&imgprofile) //* Receive data from Body(API)
	if err != nil {
		log.Println("Error: from c.Bind()")
		return c.String(http.StatusInternalServerError, "Error: from c.Bind()")
	}
	dataimg := imgprofile.Avatar
	bs64tostr, _ := base64.StdEncoding.DecodeString(string(dataimg))
	img, _, err := image.Decode(bytes.NewReader(bs64tostr))
	if err != nil {
		m := imaging.CropAnchor(img, 500, 500, imaging.Center)
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
	m := imaging.CropAnchor(img, 500, 500, imaging.Center)
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

	var imgprofile model.Profile
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

	var imgprofile model.Profile
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

	var imgprofile model.Profile
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

	var imgprofile model.Profile
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

	var imgprofile model.Profile
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
