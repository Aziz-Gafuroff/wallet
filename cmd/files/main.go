package main

import (
	"fmt"
	"log"

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

	
	// account, err := svc.RegisterAcccount(types.Phone("+992928330001"))
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }

	// err = svc.Deposit(account.ID, 100000000)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// account2, err := svc.RegisterAcccount(types.Phone("+992928330002"))
	// 	if err != nil {
	// 		log.Print(err)
	// 		return
	// 	}

	// 	err = svc.Deposit(account2.ID, 100000000)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return 
	// 	}

	// for i := 0; i < 100; i++ {

	// 	_, err := svc.Pay(account.ID, types.Money(2+i), "mobile")
	// 	if err != nil {
	// 		log.Print(err)
	// 		return
	// 	}

	// }

	// for i := 0; i < 97; i++ {

	// 	_, err := svc.Pay(account2.ID, types.Money(2+i), "mobile")
	// 	if err != nil {
	// 		log.Print(err)
	// 		return
	// 	}

	// }



	// filterPayments, err := svc.FilterPayments(account2.ID, 5)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return
	// 	}
	// log.Print(len(filterPayments))




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


	// err := svc.HistoryToFiles(payments, "data", 10)
	// if err != nil {
	// 	log.Print(err)
	// 	return 
	// }

	// err := svc.Import("data")
	// if err != nil {
	// 	log.Print(err)
	// }

	// account, err := svc.FindAccountByID(1)
	// if err != nil {
	// 	log.Print(err)
	// }

	// log.Printf("Account: %v", account )

	svc.RegisterAcccount("789")
	// srv.RegisterAccount("456")

	svc.Deposit(1, 120000000000000)

	svc.Pay(1, 1236, "mobile")
	svc.Pay(1, 31336, "Home")
	svc.Pay(1, 4656, "mobile")
	svc.Pay(1, 986, "mobile")
	svc.Pay(1, 7736, "mobile")
	svc.Pay(1, 137736, "mobile")
	svc.Pay(1, 1236, "Cafe")
	svc.Pay(1, 332126, "mobile")
	svc.Pay(1, 36133, "mobile")
	svc.Pay(1, 736, "Home")
	svc.Pay(1, 98736, "Home")

	q, err := svc.FilterPaymentsByFn(getfilteredPayment, 4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(q)

}

func getfilteredPayment(payment types.Payment) bool {
	return payment.Category == "mobile"
}