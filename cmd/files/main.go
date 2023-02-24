package main

import (
	"log"

	"github.com/Aziz-Gafuroff/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	
	// account, err := svc.RegisterAcccount("+992928330099")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Deposit(account.ID, 10)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// account, err = svc.RegisterAcccount("+992928330000")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Deposit(account.ID, 10)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// err = svc.ExportToFile("data/massage.txt")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// log.Printf("%v", account)

	err := svc.ImpoToFile("data/massage.txt")
	if err != nil {
		log.Print(err)
		return
	}

	account, err := svc.FindAccountByID(1)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("%v\n", account)

}
