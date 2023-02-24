package wallet

import (
	"github.com/google/uuid"
	"testing"
)

func TestService_FindAccountByID_nil(t *testing.T) {
	svc := Service{}
	account, err := svc.FindAccountByID(1)
	if account == nil {
		t.Error(err)
	}
}
func TestService_FindAccountByIDFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccounts("9920000001")
	account, err := svc.FindAccountByID(1)
	if err == nil {
		t.Error("Account was find ", account.Phone, nil)
	}
}

func TestService_FindPaymentByID_nil(t *testing.T) {
	svc := Service{}
	account, err := svc.FindPaymentByID("1")
	if account == nil {
		t.Error(err)
	}
}
func TestService_FindPaymentByIDFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccounts("992001")
	svc.Deposit(1, 10)
	svc.Pay(1, 10, "ansar")
	svc.FindPaymentByID(uuid.NewString())
	account, err := svc.FindPaymentByID("1")
	if account != nil {
		t.Error("It is Good", err)
	}
}
