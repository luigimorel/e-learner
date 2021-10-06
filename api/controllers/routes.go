package controllers

import "github.com/morelmiles/school-mgt-system/api/middleware"

func (s *Server) initializeRoutes() {
	//Home
	s.Router.HandleFunc("/", middleware.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Login route 
	s.Router.HandleFunc("/login", middleware.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Learning track course 
	s.Router.HandleFunc("/learning_track", middleware.SetMiddlewareJSON(s.CreateLearningTrack)).Methods("POST")
	s.Router.HandleFunc("/learning_track", middleware.SetMiddlewareJSON(s.GetLearningTrack)).Methods("GET")
	s.Router.HandleFunc("/learning_track/{id}", middleware.SetMiddlewareJSON(s.GetLearningTracks)).Methods("GET")
	s.Router.HandleFunc("/learning_track/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateLearningTrack))).Methods("PUT")
	s.Router.HandleFunc("/learning_track/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareJSON(s.DeleteLearningTrack))).Methods("DELETE")

	//Course routes
	s.Router.HandleFunc("/learning_track/{id}/courses", middleware.SetMiddlewareJSON(s.CreateCourse)).Methods("POST")
	s.Router.HandleFunc("/learning_track/{id}/courses", middleware.SetMiddlewareJSON(s.GetCourse)).Methods("GET")
	s.Router.HandleFunc("/learning_track/{id}/course/{id}", middleware.SetMiddlewareJSON(s.GetCourses)).Methods("GET")
	 s.Router.HandleFunc("/learning_track/{id}/courses/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateCourse))).Methods("PUT")
	 s.Router.HandleFunc("/learning_track/{id}/courses/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.DeleteCourse))).Methods("DELETE")

	//Module routes 
	s.Router.HandleFunc("/learning_track/{id}/course/{id}/module", middleware.SetMiddlewareJSON(s.CreateModule)).Methods("POST")
	s.Router.HandleFunc("/learning_track/{id}/course/{id}/module", middleware.SetMiddlewareJSON(s.GetModule)).Methods("GET")
	s.Router.HandleFunc("/learning_track/{id}/course/{id}/module", middleware.SetMiddlewareJSON(s.GetModules)).Methods("GET")
	s.Router.HandleFunc("/learning_track/{id}/course/{id}/module/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateModule))).Methods("PUT")
	s.Router.HandleFunc("/learning_track/{id}/course/{id}/module/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.DeleteModule))).Methods("DELETE")

	}