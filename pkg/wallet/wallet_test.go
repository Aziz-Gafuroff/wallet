package wallet

import (
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