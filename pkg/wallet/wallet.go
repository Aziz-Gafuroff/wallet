package wallet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

var (
	ErrAccountNotFound      = errors.New("Account not found")
	ErrPhoneRegistred       = errors.New("phone already registred")
	ErrAmountMustBePositive = errors.New("Amount must be positive")
	ErrNotEnoughBalance     = errors.New("Balance is not enough")
	ErrPaymentNotFound      = errors.New("Payment not found")
	ErrFavoriteNotFound     = errors.New("Favorite payment not found")
)

func (s *Service) RegisterAcccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistred
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
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
		ID:        paymaneID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusOk,
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

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {

			return favorite, nil
		}

	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, ErrFavoriteNotFound
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for _, item := range s.accounts {
		accountStr := strconv.Itoa(int(item.ID)) + ";" + string(item.Phone) + ";" + strconv.Itoa(int(item.Balance)) + "|"
		log.Print(accountStr)
		_, err = file.Write([]byte(accountStr))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

func (s *Service) ImpoToFile(path string) error {
	file, err := os.Open("data/massage.txt")
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 10)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")
	for _, item := range data {
		acc := strings.Split(item, ";")
		if len(acc) != 3 {
			continue
		}
		accountID, err := strconv.ParseInt(acc[0], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		balance, err := strconv.Atoi(acc[2])
		if err != nil {
			log.Print(err)
			return err
		}

		s.accounts = append(s.accounts, &types.Account{
			ID:      int64(accountID),
			Phone:   types.Phone(acc[1]),
			Balance: types.Money(balance),
		})
	}

	return nil
}

func (s *Service) Export(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		}
	}

	file, err := os.OpenFile(dir + "/accounts.dump", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()
	
	for _, item := range s.accounts {
		accountStr := fmt.Sprintf("%d;%s;%d\n",item.ID, item.Phone, item.Balance)
		_, err = file.Write([]byte(accountStr))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	file, err = os.Create(dir + "/payments.dump")
	if err != nil {
		return err
	}
	
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()
	
	for _, item := range s.payments {
		paymentSrt := fmt.Sprintf("%s;%d;%d;%s;%s\n",item.ID, item.AccountID, item.Amount, item.Category, item.Status)
		_, err = file.Write([]byte(paymentSrt))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	file, err = os.Create(dir + "/favorites.dump")
	if err != nil {
		return err
	}
	
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()
	
	for _, item := range s.favorites {
		favoriteStr := fmt.Sprintf("%s;%d;%s;%d;%s\n",item.ID, item.AccountID, item.Name, item.Amount, item.Category)
		_, err = file.Write([]byte(favoriteStr))
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}

func (s *Service)Import(dir string) error {
	
	file, err := os.Open(dir + "/accounts.dump")
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	
	reader := bufio.NewReader(file)
	
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		line = strings.ReplaceAll(line, "\n", "")

		acc := strings.Split(line, ";")
		if len(acc) != 3 {
			continue
		}
		accountID, err := strconv.ParseInt(acc[0], 10, 64)
		if err != nil {
			return err
		}

		balance, err := strconv.Atoi(acc[2])
		if err != nil {
			return err
		}

		s.accounts = append(s.accounts, &types.Account{
			ID:      int64(accountID),
			Phone:   types.Phone(acc[1]),
			Balance: types.Money(balance),
		})
	}

	file1, err := os.Open(dir + "/payments.dump")
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		err := file1.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	
	reader = bufio.NewReader(file1)
	
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}

		line = strings.ReplaceAll(line, "\n", "")

		acc := strings.Split(line, ";")
		if len(acc) != 5 {
			continue
		}
		accountID, err := strconv.ParseInt(acc[1], 10, 64)
		if err != nil {
			log.Print(err)
			return err
		}

		balance, err := strconv.Atoi(acc[2])
		if err != nil {
			log.Print(err)
			return err
		}

		s.payments = append(s.payments, &types.Payment{
			ID:      acc[0],
			AccountID: accountID,
			Amount: types.Money(balance),
			Category: types.PaymentCategory(acc[3]),
			Status: types.PaymentStatus(acc[4]),
		})
	}

	file2, err := os.Open(dir + "/favorites.dump")
	if err != nil {
		log.Print(err)
		return err
	}

	defer func() {
		err := file2.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	
	reader = bufio.NewReader(file2)
	
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}

		line = strings.ReplaceAll(line, "\n", "")

		acc := strings.Split(line, ";")
		if len(acc) != 5 {
			continue
		}
		accountID, err := strconv.ParseInt(acc[1], 10, 64)
		if err != nil {
			return err
		}

		balance, err := strconv.Atoi(acc[3])
		if err != nil {
			log.Print(err)
			return err
		}

		s.favorites = append(s.favorites, &types.Favorite{
			ID:      acc[0],
			AccountID: accountID,
			Name: acc[2],
			Amount: types.Money(balance),
			Category: types.PaymentCategory(acc[4]),
		})
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	newPayments := make([]types.Payment, 0)
	
	_, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}
	
	for _, payment := range s.payments {
		if payment.AccountID == accountID {
			newPayments = append(newPayments, *payment)
		}
	}

	if len(newPayments) < 1 {
		return nil, ErrPaymentNotFound
	}


	return newPayments, nil
}

func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	fileNames := "payments"
	if len(payments) > records{
		countParts := 0
		parts := len(payments) / records
		for i := 1; i <= parts; i++ {
			saveHistoryToFile(payments[countParts:i*records], dir, fmt.Sprintf("%s%d",fileNames,i))
			countParts = i * records
		}
	} else {
		saveHistoryToFile(payments, dir, fileNames)
	}
	return nil
}

func saveHistoryToFile(payments []types.Payment, dir string, fileName string) error {
	
	file, err := os.OpenFile(dir + "/"+fileName + ".dump", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = cerr
		}
	}()
	
	for _, payment := range payments {
		paymentStr := fmt.Sprintf("%s;%d;%d;%s;%s\n",payment.ID, payment.AccountID, payment.Amount, payment.Category, payment.Status)
		_, err = file.Write([]byte(paymentStr))
		if err != nil {
			log.Print(err)
			return err
		}
	}
	
	return nil
	
}


