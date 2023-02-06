package wallet

import (
	"errors"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}

var ( 
	ErrAccountNotFound = errors.New("Account not found")
	ErrPhoneRegistred = errors.New("phone already registred")
	ErrAmountMustBePositive = errors.New("Amount must be positive")
	ErrNotEnoughBalance = errors.New("Balance is not enough")
	ErrPaymentNotFound = errors.New("Payment not found")
)

func (s *Service) RegisterAcccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

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


func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymaneID := uuid.New().String()
	payment := &types.Payment{
		ID: paymaneID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusOk,
	}

	s.payments = append(s.payments, payment)
	
	return payment, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil 
		}
	}


	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}
	
	account.Balance = amount

	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	payment, err = s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
