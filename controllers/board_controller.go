package controllers

import (
	"math"
	"strconv"

	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/services"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error
	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Invalid read request", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Invalid user id", err.Error())
	}
	board.OwnerPublicID = userID

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Invalid save data", err.Error())
	}
	return utils.Succes(ctx, "Board created successfully", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Invalid read request", err.Error())
	}
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid board id", err.Error())
	}
	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board not found", err.Error())
	}
	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Invalid update data", err.Error())
	}
	return utils.Succes(ctx, "Board updated successfully", board)
}

func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Invalid Parse Data", err.Error())
	}
	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Invalid Add Members", err.Error())
	}
	return utils.Succes(ctx, "Success Add Members", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Invalid Parse Data", err.Error())
	}
	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Invalid Remove Members", err.Error())
	}
	return utils.Succes(ctx, "Success Remove Members", nil)
}

func (c *BoardController) GetMyBoardsPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page , _ := strconv.Atoi(ctx.Query("page", "1"))
	limit , _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.GetAllByUserPaginate(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Invalid Get Data Boards", err.Error())
	}

	meta := utils.PaginationMeta{
		Page : page,
		Limit : limit,
		Total : int(total),
		TotalPage : int(math.Ceil(float64(total) / float64(limit))),
		Filter : filter,
		Sort : sort,
	}
	return utils.SuccesPagination(ctx, "Success Get Data Boards", boards, meta)
}