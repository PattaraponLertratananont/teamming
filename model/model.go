package model

import "github.com/globalsign/mgo/bson"

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
