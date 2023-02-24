package wallet

import (
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
