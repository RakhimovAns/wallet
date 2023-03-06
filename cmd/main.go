package main

import (
	"fmt"
	"github.com/RakhimovAns/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := wallet.Service{}

	/*	_, err := svc.RegisterAccounts("+992000001")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = svc.RegisterAccounts("+992000002")
		if err != nil {
			fmt.Println(err)
			return
		}*/
	/*err = svc.ExportToFile("pkg/wallet/accounts.txt")
	if err != nil {
		log.Print(err)
		return
	}*/
	err := svc.ImportFromFile("pkg/wallet/accounts.txt")
	if err != nil {
		log.Print(err)
		return
	}
	for _, account := range svc.Accounts {
		fmt.Println(account.Phone)
	}
}
