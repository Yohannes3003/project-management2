package controllers

import (
	"math"
	"strconv"

	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/services"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
} 

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)

	return utils.Succes(ctx, "Register Success", userResp)
}

func (c *UserController) Login (ctx *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	user, err := c.service.Login(body.Email, body.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken , _ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, &user)
	return utils.Succes(ctx, "Login Succes", fiber.Map{
		"access_token" : token,
		"refresh_token" : refreshToken,
		"user" : userResp,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicId(id)
	if err != nil {
		return utils.NotFound(ctx, "Data not found", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil {
		return utils.BadRequest(ctx, "Internal server error", err.Error())
	}
	return utils.Succes(ctx, "Data berhasil ditemuka", userResp)
} 

func (c *UserController) GetUserPagination (ctx *fiber.Ctx) error {
	// /users/page?page=1&limit=10&sort=-id&filter=rafael
	//100/10 data = 10 page
	page, _ := strconv.Atoi(ctx.Query("page","1"))
	limit, _ := strconv.Atoi(ctx.Query("limit","10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort","")

	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil data", err.Error())
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)

	meta := utils.PaginationMeta {
		Page: page,
		Limit : limit,
		Total: int(total),
		TotalPage: int(math.Ceil(float64(total)/float64(limit))),
		Filter: filter,
		Sort: sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data user not found", userResp, meta)
	}

	return utils.SuccesPagination(ctx, "Data user found", userResp, meta)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to parse data", err.Error())
	}
	user.PublicID = publicID
	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to update user", err.Error())
	}

	userUpdated, err := c.service.GetByPublicId(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to get data", err.Error())
	}
	var userResp models.UserResponse
	err = copier.Copy(&userResp, &userUpdated)
	if err != nil {
		return utils.InternalServerError(ctx, "Internal server error", err.Error())
	}
	return utils.Succes(ctx, "User updated successfully", userResp)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.service.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Failed to delete user", err.Error())
	}
	return utils.Succes(ctx, "User deleted successfully", id)
}