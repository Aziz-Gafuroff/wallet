package wallet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

type Progress struct {
	Part int
	Result types.Money
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

func (s *Service) SumPayments(goroutines int) types.Money {
	mu := sync.Mutex{}
	var sumAll types.Money
	wg := sync.WaitGroup{}
	part := len(s.payments) / goroutines
	countParts := 0

	if goroutines <= 1 {
		sumAll = saveSumPayments(s.payments)
	} else	{
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			go func(j int) {
				mu.Lock()
				defer wg.Done()
				sum := saveSumPayments(s.payments[countParts:countParts + part])
				countParts += part
				sumAll += sum
				mu.Unlock()
			}(i)
		}
		wg.Wait()
	}
	return sumAll
}

func saveSumPayments(payments []*types.Payment) types.Money {
	sum := types.Money(0)
	if len(payments) == 0 {
		return 0
	}

	for _, payment := range payments {
		sum += payment.Amount
	}

	return sum
}

func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {
	mu := sync.Mutex{}
	var filterPayments []types.Payment
	wg := sync.WaitGroup{}
	
	partLimit := int(math.Ceil(float64(len(s.payments))/float64(goroutines)))
	lenPayments := len(s.payments)
	part := []int{}
	for i := 0; i < goroutines; i++ {
		if lenPayments <= partLimit {
			part = append(part,lenPayments)
			continue
		}
		part = append(part, partLimit)
		lenPayments -= partLimit
	}

	countParts := 0

	if goroutines <= 1 {
		filterPayments = partFilterPayments(s.payments, accountID)
	} else	{
		for i := 1; i <= goroutines; i++ {
			wg.Add(1)
			go func(j int) {
				mu.Lock()
				defer wg.Done()
				
				newPayments := partFilterPayments(s.payments[countParts:countParts + part[j-1]], accountID)
				countParts += part[j-1]
				filterPayments = append(filterPayments, newPayments...)
				mu.Unlock()
			}(i)
		}
		wg.Wait()
	}
	return filterPayments, nil
}

func partFilterPayments(payments []*types.Payment, accountID int64) []types.Payment {
	var newPayments []types.Payment
	for _, payment := range payments {
		if payment.AccountID == accountID {
			newPayments = append(newPayments, *payment)
		}
		
	}
	return newPayments
}

func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int,) ([]types.Payment, error) {
	if goroutines < 2 {
		return filterByFn(filter, s.payments), nil
	}
	
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	
	lengthPayment := len(s.payments)
	part := partsPayment(lengthPayment, goroutines)
	
	var filteredPayments []types.Payment

	for i := 1; i <= goroutines; i++ {
		start, finish := i*part-part, i*part

		if start > len(s.payments)-1 {
			break
		}

		wg.Add(1)
		go func(start, finish int) {
			var fp []types.Payment
			defer wg.Done()
			if finish < lengthPayment {
				fp = filterByFn(filter, s.payments[start:finish])
			} else {
				fp = filterByFn(filter, s.payments[start: ])
			}
			mu.Lock()
			filteredPayments = append(filteredPayments, fp...)
			mu.Unlock()
		}(start, finish)
	}

	
	
	wg.Wait()

	return filteredPayments, nil
}

func partslengthPayment(lengthPayment, goroutines int) {
	panic("unimplemented")
}

func filterByFn(p func(payment types.Payment) bool, payments []*types.Payment) []types.Payment {
	var filPayment []types.Payment

	for _, payment := range payments {
		if p(*payment) {
			filPayment = append(filPayment, *payment)
		}
		
	}
	return filPayment
}

func partsPayment(sliceLengthpayment int, parts int) int {
	if sliceLengthpayment%parts == 0 {
		return sliceLengthpayment / parts
	}
	return sliceLengthpayment / parts + 1
}

func (s *Service) SumPaymentsWithProgress() <-chan Progress {
	var wg sync.WaitGroup
	const elemInPart = 40

	chunkSum := make(chan Progress)
	
	partsNum := parts(len(s.payments), elemInPart)
	wg.Add(partsNum)

	for i := 0; i < partsNum; i++ {
		start, finish := i*elemInPart, (i+1)*elemInPart
		if finish > len(s.payments) {
			finish = len(s.payments)
		}

		j := i
		go func(ch chan<- Progress, data []*types.Payment) {
			defer wg.Done()
			progress := Progress {
				Part: j + 1,
				Result: sum(data),
			}

			ch <- progress
		}(chunkSum, s.payments[start:finish])
	}

	go func() {
		defer close(chunkSum)
		wg.Wait()
	}()

	return chunkSum
}

func sum(payments []*types.Payment) types.Money {
	var tl types.Money

	if len(payments) == 0 {
		return 0
	}

	for _, payment := range payments {
		tl += payment.Amount
		
	}
	return tl

}

func parts(sliceLength int, m int) int {
	if sliceLength == 0 {
		return sliceLength / m
	}

	return sliceLength / m + 1
}