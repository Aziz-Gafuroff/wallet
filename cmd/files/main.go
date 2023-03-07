package main

import (
	"log"
	"strconv"

	"github.com/Aziz-Gafuroff/wallet/pkg/types"
	"github.com/Aziz-Gafuroff/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	
	// account, err := svc.RegisterAcccount("+992928330099")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Deposit(account.ID, 1000)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// payment, err := svc.Pay(account.ID, 9, "auto")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// favorite, err := svc.FavoritePayment(payment.ID, "auto")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// _, err = svc.PayFromFavorite(favorite.ID)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// account, err = svc.RegisterAcccount("+992928330000")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Deposit(account.ID, 1000)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// payment, err = svc.Pay(account.ID, 8, "mobile")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// favorite, err = svc.FavoritePayment(payment.ID, "mobile")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// _, err = svc.PayFromFavorite(favorite.ID)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Export("data")
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

	// err := svc.ImpoToFile("data/massage.txt")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// account, err := svc.FindAccountByID(1)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// log.Printf("%v\n", account)

	var payments []types.Payment

	for i := 0; i < 100; i++ {
		account, err := svc.RegisterAcccount(types.Phone("+99292833000"+strconv.Itoa(i)))
		if err != nil {
			log.Print(err)
			return
		}

		err = svc.Deposit(account.ID, 1000)
		if err != nil {
			log.Print(err)
			return 
		}

		payment, err := svc.Pay(account.ID, types.Money(8+i), "mobile")
		if err != nil {
			log.Print(err)
			return
		}

		payments  = append(payments, *payment)
	}

	// account, err := svc.FindAccountByID(1)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// log.Printf("%v\n", account)

	// payment, err := svc.FindPaymentByID("644abb7e-8446-4f87-8a52-ed329360c25a")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// log.Printf("%v\n", payment)

	// favorite, err := svc.FindFavoriteByID("ab918391-918c-4d4e-a2b0-a5c644be3933")
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// log.Printf("%v\n", favorite)


	err := svc.HistoryToFiles(payments, "data", 10)
	if err != nil {
		log.Print(err)
		return 
	}

	// err := svc.Import("data")
	// if err != nil {
	// 	log.Print(err)
	// }

	// account, err := svc.FindAccountByID(1)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Printf("Account: %v", account )
}
