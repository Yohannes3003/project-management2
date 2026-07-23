package routes

import (
	"log"

	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/controllers"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
)

func SetUp(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController) {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)

	// JWT Protected routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func (c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error Unauthorized", err.Error())
		},
	}))

	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination) // api/v1/users/page
	userGroup.Get("/:id", uc.GetUser) // api/v1/users/:id
	userGroup.Put("/:id", uc.UpdateUser) // api/v1/users/:id
	userGroup.Delete("/:id", uc.DeleteUser) // api/v1/users/:id

	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard) // api/v1/boards

}