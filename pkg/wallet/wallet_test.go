package wallet

import (
	"reflect"
	"testing"
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