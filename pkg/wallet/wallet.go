package wallet

import (
	"errors"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
)

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}

var ( 
	ErrAccountNotFound = errors.New("Account not found")
	ErrPhoneRegistred = errors.New("phone already registred")
)

func (service *Service) RegisterAcccount(phone types.Phone) (*types.Account, error) {
	for _, account := range service.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
		
	}

	service.nextAccountID++
	account := &types.Account{
		ID: service.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	service.accounts = append(service.accounts, account)

	return account, nil

}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil 
		}
	}


	return nil, ErrAccountNotFound
	
}