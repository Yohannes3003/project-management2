package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Yohannes3003/project-management2/config"
	"github.com/Yohannes3003/project-management2/models"
	"github.com/Yohannes3003/project-management2/models/types"
	"github.com/Yohannes3003/project-management2/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardServices interface {
	Create(card *models.Card, listPublicID string) error 
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error
	GetByListID(listPublicID string) ([]models.Card, error)	
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
) CardServices {
	return &cardService{cardRepo, listRepo, userRepo}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("List Not Found : %w", err)
	}

	card.ListID = list.InternalID

	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to Create Card : %w", err)
	}

	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).Where("list_internal_id = ?", list.InternalID).First(&position).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			position = models.CardPosition{
				PublicID: uuid.New(),
				ListID: list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}
			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to Create Card Position : %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("Failed to Get Card Position : %w", err)
		}
	} else {
		position.CardOrder = append(position.CardOrder, card.PublicID)
		if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", position.InternalID).Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to Update Card Position : %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Transcation Commit Failed : %w", err)
	}
	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("Card Not Found : %w", err)
	}

	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("List Not Found : %w", err)
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if existingCard.ListID != newList.InternalID {
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ? ", existingCard.ListID).First(&oldPos).Error; err != nil {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPos.InternalID).Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("Failed to Update Old Card Position : %w", err) 
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("Failed to Get Old Card Position : %w", err)
		}
	}

	var newPos models.CardPosition
	res := tx.Where("List_internal_id = ? ", newList.InternalID).First(&newPos)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		newPos = models.CardPosition{
			PublicID: uuid.New(),
			ListID: newList.InternalID,
			CardOrder: types.UUIDArray{existingCard.PublicID},
		}
		if err := tx.Create(&newPos).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to Create Card Position for New List : %w", err)
		}
	} else if res.Error == nil {
		updateOrder := append(newPos.CardOrder, existingCard.PublicID)
		if err := tx.Model(&models.CardPosition{}).Where("internal_id = ? ", newPos.InternalID).Update("card_order", types.UUIDArray(updateOrder)).Error ; err != nil {
			tx.Rollback()
			return fmt.Errorf("Failed to Update New Card Position : %w", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("Failed to Get New Card Position : %w", res.Error)
	}
	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to Update Card : %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Transaction Commit Failed : %w ", err)
	}
	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("List Not Found : %w", err)
	}
	position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get Card Position : %w", err)
	}
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get Card : %w", err)
	}

	if position != nil && len(position.CardOrder) > 0 {
		cards = sortCardByPosition(cards, position.CardOrder)
	}
	return cards, nil
}

func sortCardByPosition(cards []models.Card, order []uuid.UUID) []models.Card {
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}
	
	defaultIndex := len(order)
	sort.SliceStable(cards, func (i, j int) bool {
		idxI, okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}
		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}
		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})
	return cards
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}

