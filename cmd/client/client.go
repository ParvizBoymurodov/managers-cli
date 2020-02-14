package main

import (
	"database/sql"
	"fmt"
	"github.com/ParvizBoymurodov/managers-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)


func main() {
	// os.Stdin, os.Stout, os.Stderr, File
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)

	}
	defer func() {
		log.Print("close db")
		if err := db.Close(); err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db: %v", err)
	}

	fmt.Fprintln(os.Stdout, "Добро пожаловать!")
	log.Print("start operations loop")
	operationsLoop(-1, db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}

func operationsLoop(user_id int64, db *sql.DB, commands string, loop func(db *sql.DB, cmd string, user_id int64) (exit bool)) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd), user_id); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string, user_id int64) (exit bool) {
	switch cmd {
	case "1":
		id, ok, err := handleLoginForClient(db)
		if err != nil {
			fmt.Println("Неправильно введён логин или пароль")
			log.Printf("can't handle login: %v", err)
			return false
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")
			//unauthorizedOperationsLoop(db, "1")
			//Graceful shutdown
			return false
		}
		user_id = id
		operationsLoop(user_id, db, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string,user_id int64) (exit bool) {
	switch cmd {
	 // TODO: may be log fatal
	 case "1":
		listBalance, err := core.GetBalanceList(db, user_id)
		if err != nil {
			 log.Printf("can't get all products: %v", err)
			 return true // TODO: may be log fatal
		}
		printClientBalance(listBalance)
		
	case "2":
		err:=transaction(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return false // TODO: may be log fatal
		}
	case "3":
		err:=transactionByBalanceNumber(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return false // TODO: may be log fatal
		}
	case "4":
		serviceList, err := core.GetServices(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printServiceList(serviceList)
	case "5":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printAtm(atms)

	
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}


func printAtm(atms []core.Atm) {
	for _, atm := range atms {
		fmt.Printf(
			"id: %d, name: %s, street:%s\n",
			atm.Id,
			atm.Name,
			atm.Street,
		)
	}
}

func printClientBalance(listBalance []core.Client)  {
	for _, clientAccounts := range listBalance {
		fmt.Printf(
			"id: %d, name: %s, balanceNumber: %d, balance:%d\n",
			clientAccounts.Id,
			clientAccounts.Name,
			clientAccounts.BalanceNumber,
			clientAccounts.Balance,
		)
	}
}

func printServiceList(serviceList []core.Services)  {
	for _, listService := range serviceList {
		fmt.Printf(
			"id: %d, name: %s, price:%d\n",
			listService.Id,
			listService.Name,
			listService.Price,

		)
	}
}

func handleLoginForClient(db *sql.DB) (id int64,ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return -1,false, err
	}

	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return -1,false, err
	}

	id, ok, err = core.Login(login, password, db)
	if err != nil {
		return -1,false, err
	}

	return id,ok, err
}

func transaction(db *sql.DB)(err error)  {

	var myPhoneNumber int64
	fmt.Print("Введите свой номер телефон: ")
	_, err = fmt.Scan(&myPhoneNumber)
	if err != nil {
		return err
	}
	var phoneNumber int64
	fmt.Print("Введите номер телефон клиента: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}
	err = core.CheckByPhoneNumber(phoneNumber, db)
	if err != nil {
		fmt.Println("fddsdf")
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.TransferByPhoneNumber(myPhoneNumber, balance, core.Client{
		Balance:     balance,
		PhoneNumber: phoneNumber,
	}, db)
	if err != nil {
		fmt.Println("Извините у вас мало денег")
		return err
	}else if myPhoneNumber==phoneNumber {
		fmt.Println("Yt yflj nfr")
		return err
	}
		fmt.Println("Денги переведенный успешно переведенный!")
		return nil
}

func transactionByBalanceNumber(db *sql.DB)(err error)  {

	var myBalanceNumber uint64
	fmt.Print("Введите номер своего баланса: ")
	_, err = fmt.Scan(&myBalanceNumber)
	if err != nil {
		return err
	}
	var balanceNumber uint64
	fmt.Print("Введите номер баланс клиента: ")
	_, err = fmt.Scan(&balanceNumber)
	if err != nil {
		return err
	}
	err = core.CheckByBalanceNumber(balanceNumber, db)
	if err !=nil{
		fmt.Println("щшфырыфр")
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}


	err = core.TransferByBalanceNumber(myBalanceNumber, balance, core.Client{
		Balance:     balance,
		BalanceNumber: balanceNumber,
	}, db)
	if err != nil {
		fmt.Println("Извините у вас мало денег")
		return err
	}else if myBalanceNumber==balanceNumber {
		fmt.Println("Yt yflj nfr")
		return err
	}
	fmt.Println("Денги переведенный успешно переведенный!")
	return nil

}