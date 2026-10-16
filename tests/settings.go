package tests

import (
	"github.com/Izhdenev/Final-Project/pkg/db"
)

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiNTk5NDQ3MWFiYjAxMTEyY
WZjYzE4MTU5ZjZjYzc0YjRmNTExYjk5ODA2ZGE1OWIzY2FmNWE5YzE3M2NhY2ZjNSJ9.yKZCInwMk68u
cbGp2Pf8WEg_m58MHAORXzJOjw-95B4`
var WebDir = "web"

func init() {
	if err := db.Init(DBFile); err != nil {
		panic("failed to init DB: " + err.Error())
	}
}
