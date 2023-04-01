// Code generated by datamapper.
// https://github.com/underbek/datamapper

// Package mapper is a generated datamapper package.
package mapper

import (
	"errors"
	"fmt"

	"github.com/underbek/datamapper/_test_data/mapper/with_dash_and_pointers/convertors"
	db "github.com/underbek/datamapper/_test_data/mapper/with_dash_and_pointers/dao"
	"github.com/underbek/datamapper/_test_data/mapper/with_dash_and_pointers/domain"
	"github.com/underbek/datamapper/_test_data/mapper/with_dash_and_pointers/domain/user"
	"github.com/underbek/datamapper/converts"
)

// ConvertDomainOrderToDbOrderData convert *domain.Order by tag map to db.OrderData by tag db
func ConvertDomainOrderToDbOrderData(from *domain.Order) (db.OrderData, error) {
	if from == nil {
		return db.OrderData{}, errors.New("Order is nil")
	}

	if from.OrderID == nil {
		return db.OrderData{}, errors.New("cannot convert *domain.Order.OrderID -> db.OrderData.Order.ID, field is nil")
	}

	fromOrderID, err := converts.ConvertStringToSigned[int64](*from.OrderID)
	if err != nil {
		return db.OrderData{}, fmt.Errorf("convert Order.OrderID -> OrderData.Order.ID failed: %w", err)
	}

	fromAdditions := make([]db.Additional, 0, len(from.Additions))
	for _, item := range from.Additions {
		fromAdditions = append(fromAdditions, convertors.ConvertDomainAdditionalToDaoAdditional(item))
	}

	if from.User == nil {
		return db.OrderData{}, errors.New("Order.User is nil")
	}

	fromUserID, err := converts.ConvertStringToSigned[int64](from.User.ID)
	if err != nil {
		return db.OrderData{}, fmt.Errorf("convert Order.User.ID -> OrderData.UserData.ID failed: %w", err)
	}

	if from.User.UserTimes == nil {
		return db.OrderData{}, errors.New("Order.User.UserTimes is nil")
	}

	return db.OrderData{
		Order: &db.Order{
			ID:        fromOrderID,
			UUID:      from.OrderUUID,
			Additions: fromAdditions,
		},
		UserData: &db.User{
			ID:        fromUserID,
			CreatedAt: from.User.UserTimes.CreatedAt,
		},
		Urls: db.OrderUrls{
			SiteUrl:     from.SiteUrl,
			RedirectUrl: from.RedirectUrl,
		},
	}, nil
}

// ConvertDbOrderDataToDomainOrder convert db.OrderData by tag db to *domain.Order by tag map
func ConvertDbOrderDataToDomainOrder(from db.OrderData) (*domain.Order, error) {
	if from.Order == nil {
		return nil, errors.New("OrderData.Order is nil")
	}

	fromOrderID := converts.ConvertNumericToString(from.Order.ID)

	if from.UserData == nil {
		return nil, errors.New("OrderData.UserData is nil")
	}

	fromOrderAdditions := make([]domain.Additional, 0, len(from.Order.Additions))
	for _, item := range from.Order.Additions {
		fromOrderAdditions = append(fromOrderAdditions, convertors.ConvertDaoAdditionalToDomainAdditional(item))
	}

	return &domain.Order{
		OrderID:     &fromOrderID,
		OrderUUID:   from.Order.UUID,
		SiteUrl:     from.Urls.SiteUrl,
		RedirectUrl: from.Urls.RedirectUrl,
		Additions:   fromOrderAdditions,
		User: &user.User{
			ID: converts.ConvertNumericToString(from.UserData.ID),
			UserTimes: &user.Times{
				CreatedAt: from.UserData.CreatedAt,
			},
		},
	}, nil
}
