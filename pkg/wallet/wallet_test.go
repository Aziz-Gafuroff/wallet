package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
)

func TestFindAccountByID_nil(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAcccount("+992928330099")
	if err != nil {
		t.Error(err)
		return
	}

	result, err := svc.FindAccountByID(account.ID)

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(result, account) {
		t.Errorf("result: %v, expected: %v", result, account)
	}


}

func TestFindAccountByID_empty(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindAccountByID(10)

	if err != ErrAccountNotFound {
		t.Error(err)
		return
	}

}

func TestFindPaymentByID_nil(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAcccount("+992928330099")
	if err != nil {
		t.Error(err)
		return
	}
	
	account.Balance = 6_000_00
	payment, err := svc.Pay(account.ID, 5_000_00, "food")

	if err != nil {
		t.Error(err)
		return
	}

	result, err := svc.FindPaymentByID(payment.ID)

	if !reflect.DeepEqual(result, payment) {
		t.Errorf("result: %v, expected: %v", result, payment)
	}


}

func TestFindPaymenttByID_empty(t *testing.T) {
	svc := &Service{}

	_, err := svc.FindPaymentByID("")

	if err != ErrPaymentNotFound {
		t.Error(err)
		return
	}

}

func TestReject_succes(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAcccount("+992928330099")
	if err != nil {
		t.Error(err)
		return
	}
	
	account.Balance = 6_000_00
	payment, err := svc.Pay(account.ID, 5_000_00, "food")

	if err != nil {
		t.Error(err)
		return
	}

	 err = svc.Reject(payment.ID)

	if payment.Status != types.PaymentStatusFail {
		t.Errorf("Payment not rejected")
		return 
	}
	
	if err != nil {
		t.Error(err)
		return
	}


}


func TestReject_fail(t *testing.T) {
	svc := &Service{}

	 err := svc.Reject("")
	
	if err == nil {
		t.Error(err)
		return
	}


}

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone: "+992928330099",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	// регистрируем пользователя

	account, err := s.RegisterAcccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	//пополнение счета
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	} 

	//выполняем платежи
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		//тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}
func TestRepeat_succsess(t *testing.T) {
	tsvc := newTestService()

	_, payments, err := tsvc.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	payment, err := tsvc.Repeat(payments[0].ID)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	if reflect.DeepEqual(payments[0], payment) {
		t.Errorf("result: %v, expexted: %v", payments[0], payment)
		return
	}
}