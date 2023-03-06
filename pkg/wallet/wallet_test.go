package wallet

import (
	"fmt"
	"github.com/RakhimovAns/wallet/pkg/types"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

var defaultTestAccount = testAccount{
	phone:   "9920000001",
	Balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

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

func TestService_Reject_Success(t *testing.T) {
	s := &Service{}
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(nil)
		return
	}
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by ID = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentsStatusFail {
		t.Errorf("Reject(): status didn't changed ,payment= %v", savedPayment)
		return
	}
	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find acccount by ID error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.Balance {
		t.Errorf("Reject(): Balance didn't changed account = %v", savedAccount)
		return
	}
}

func TestService_FindPaymentByID_Success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentsByID(): error = %v", err)
		return
	}
	if reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment return = %v", err)
		return
	}
}
func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentsByID(): must return error not nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound,return =%v", err)
		return
	}
}
func TestService_Repeat_Fail(t *testing.T) {
	s := Service{}
	account, _, err := s.addAccount(defaultTestAccount)
	if err == nil {
		t.Error(err)
		return
	}
	payment, err := s.Pay(account.ID, 9, "auto")
	if err == nil {
		t.Error("Repeat():  err shouldn't be nil")
		return
	}
	pay, err := s.Repeat(payment.ID)
	if reflect.DeepEqual(pay, payment) == false {
		t.Errorf("They should be  equal ,return =%v", err)
		return
	}
}
func TestService_Repeat_Success(t *testing.T) {
	s := Service{}
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment, err := s.Pay(account.ID, 9, "auto")
	if err != nil {
		t.Error("Repeat(): should be nil ")
		return
	}
	pay, err := s.Repeat(payment.ID)
	if reflect.DeepEqual(pay, payment) == true {
		t.Errorf("They should be not equal,return = %v", err)
		return
	}
}

func TestService_PayFromFavorite_rules(t *testing.T) {
	s := Service{}
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment, err := s.FavoritePayment(uuid.New().String(), "megafon")
	fmt.Println(payment.ID)
}
