package main

import (
	"fmt"
	"github.com/RakhimovAns/wallet/pkg/types"
	"github.com/RakhimovAns/wallet/pkg/wallet"
	"reflect"
)

func main() {
	svc := wallet.Service{}
	account, err := svc.RegisterAccounts("+9920000001")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrPhoneNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}
	payment, err := svc.Pay(account.ID, 9, "auto")
	var pay *types.Payment
	pay, _ = svc.Repeat(payment.ID)
	if reflect.DeepEqual(pay, payment) == false {
		fmt.Println("It is not working")
	} else {
		fmt.Println("it is working")
	}
}
