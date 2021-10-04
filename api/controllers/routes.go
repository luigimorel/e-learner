package controllers

import "github.com/morelmiles/school-mgt-system/api/middleware"

func (s *Server) initializeRoutes() {
	//Home
	s.Router.HandleFunc("/", middleware.SetMiddlewareJSON(s.Home)).Methods("GET")
}