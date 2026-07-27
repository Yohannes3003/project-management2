package controllers

import (
	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/services"
	"github.com/Yohannes3003/project-management2/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Failed Read Request", err.Error())
	}
	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Failed Create List", err.Error())
	}
	return utils.Succes(ctx, "Succes Create List", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list  := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Failed Parse Data", err.Error())
	}

	if _, err := uuid.Parse(publicID) ; err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Not Found", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Failed Update List", err.Error())
	}

	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Not Found", err.Error())
	}

	return utils.Succes(ctx, "Succes Create List", updatedList)

}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx)  error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID) ; err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Lists Not Found", err.Error())
	}

	return utils.Succes(ctx, "Succes Get Data", lists)

}

func (c *ListController) DeleteLits(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id") 
	if _, err := uuid.Parse(publicID) ; err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Not Found", err.Error())
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Failed Delete List", err.Error())
	} 
	return utils.Succes(ctx, "Succes Remove List", publicID)
}

func (c *ListController) UpdateListPosition(ctx *fiber.Ctx) error {
	boardID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardID) ; err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}
	var positionUUID []uuid.UUID
	if err := ctx.BodyParser(&positionUUID); err != nil {
		//if failed try parse array of string
		var positionString []string
		if err := ctx.BodyParser(&positionString); err != nil {
			return utils.BadRequest(ctx, "Invalid Position Format", err.Error())
		}
		for _, s := range positionString {
			u, err := uuid.Parse(s)
			if err != nil {
				return utils.BadRequest(ctx, "Invalid Position Format", err.Error())
			}
			positionUUID = append(positionUUID, u)
		}
	}
	if err := c.service.UpdatePosition(boardID, positionUUID); err != nil {
		return utils.InternalServerError(ctx, "Failed Update List", err.Error())
	}
	return utils.Succes(ctx, "Succes Update Position List", nil)
}