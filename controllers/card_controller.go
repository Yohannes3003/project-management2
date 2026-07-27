package controllers

import (
	"time"

	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/services"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardServices
}

func NewCardController(s services.CardServices) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate time.Time `json:"due_date"`
		Position int `json:"position"`
	}

	var req CreateCardRequest
	if err :=  ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to Get Data", err.Error())
	}

	card := &models.Card {
		Title: req.Title,
		Description: req.Description,
		DueDate: &req.DueDate,
		Position: req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed Create Card", err.Error())
	}
	
	return utils.Succes(ctx, "Succes Create Card", card )
}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type UpdateCardRequest struct {
		ListPublicID string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate  *time.Time `json:"due_date"`
		Position int `json:"position"`
	}

	var req UpdateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed Parse Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	card := &models.Card {
		Title: req.Title,
		Description: req.Description,
		DueDate: req.DueDate,
		Position: req.Position,
		PublicID: uuid.MustParse(publicID),
	}

	if err := c.service.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed Update Data", err.Error())
	}

	return utils.Succes(ctx, "Success Update Data", card)
}

func (c *CardController) DeleteCard (ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card Not Found", err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Failed Delete Data", err.Error())
	}

	return utils.Succes(ctx, "Success Delete Card", publicID)
}

func (c *CardController) GetListCard(ctx *fiber.Ctx) error {
	listID := ctx.Params("list_id")
	if _, err := uuid.Parse(listID); err != nil {
		return utils.BadRequest(ctx, "ID List Not Valid", err.Error())
	}

	cards, err := c.service.GetByListID(listID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to Get Data", err.Error())
	}

	return utils.Succes(ctx, "Succes Get Data Card", cards)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	card, err := c.service.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Error While Retrieving Data", err.Error())
	}
	if card == nil {
		return utils.NotFound(ctx, "Data Card Not Found", err.Error())
	}

	return utils.Succes(ctx, "Succes Get Data", card)
}