package routes

import "github.com/tedsuo/rata"

const (
	Env   = "ENV"
	Hello = "HELLO"
	Exit  = "EXIT"
	Index = "INDEX"
	Port  = "PORT"
)

var Routes = rata.Routes{
	{Path: "/", Method: "GET", Name: Hello},
	{Path: "/env", Method: "GET", Name: Env},
	{Path: "/exit", Method: "GET", Name: Exit},
	{Path: "/index", Method: "GET", Name: Index},
	{Path: "/port", Method: "GET", Name: Port},
}
