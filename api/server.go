package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/morelmiles/school-mgt-system/api/controllers"
	"github.com/morelmiles/school-mgt-system/api/seed"
)

var server = controllers.Server{}

func init() {
	//loads values from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Could not load env file😏")
	}
}

func Run()  {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env file %v", err)
	} else {
		fmt.Println("We are getting the .env file 🎉🙈")
	}
		server.Init( os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

		seed.Load(server.DB)

		server.Run(":8080")
}