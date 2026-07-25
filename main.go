package main

import (
	"log"

	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/controllers"
	"github.com/Yohannes3003/project-management2/database/seed"
	"github.com/Yohannes3003/project-management2/repositories"
	"github.com/Yohannes3003/project-management2/routes"
	"github.com/Yohannes3003/project-management2/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	
	app := fiber.New()
	//user
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	//board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	//list
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepository()
	listService := services.NewListService(listRepo, boardRepo, listPosRepo)
	listController := controllers.NewListController(listService)
	
	routes.SetUp(app, userController, boardController, listController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)

	log.Fatal(app.Listen(":" + port))
}