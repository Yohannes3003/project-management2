package repositories

import (
	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/models"
)

type BoardRepository interface {
	Create(board *models.Board) error
}

type boardRepository struct {
}

func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}