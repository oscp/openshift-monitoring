package models

type Type struct {
	Name string
}

var TYPE_NEW_DEAMON = Type{Name : "NEW_DEAMON" }
var TYPE_ALL_DEAMONS = Type{Name: "ALL_DEAMONS"}
var TYPE_DEAMON_LEFT = Type{Name: "DEAMON_LEFT"}