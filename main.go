package main

import "net/http"


type Handler struct{}

type Server struct{
	Addr string
	Handler Handler 
}

func (Handler ) ServeHTTP(http.ResponseWriter, *http.Request) {}


func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/", Handler{})
	
}