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

// TODO: для тех, кто хочет попробовать, можете использовать структуры и методы:
//type manager struct {
//	db  *sql.DB
//	out io.Writer
//	in  io.Reader
//}
//
//func newManagerCLI(db *sql.DB, out io.Writer, in io.Reader) *manager {
//	return &manager{db: db, out: out, in: in}
//}

// Writer, Reader

func main() {
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
	operationsLoop(db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}

func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := handleLogin(db)
		if err != nil {
			fmt.Println("Неправильно введён логин или пароль")
			log.Printf("can't handle login: %v", err)
			return false
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")
			return false
		}
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		err := addClient(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return false // TODO: may be log fatal
		}
	case "2":
		err := updateBalance(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return false
		}
	case "3":
		err := addServices(db)
		if err != nil {
			log.Printf("can't add sale: %v", err)
			return true
		}

	case "4":
		err := addAtm(db)
		if err != nil {
			log.Printf("can't add sale: %v", err)
			return true
		}

	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func handleLogin(db *sql.DB) (ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return false, err
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return false, err
	}

	ok, err = core.LoginForManagers(login, password, db)
	if err != nil {
		return false, err
	}

	return ok, err
}

func addClient(db *sql.DB) (err error) {
	fmt.Println("Введите данные клиента")
	var name string
	fmt.Print("Имя клиента: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	//name = name[:len(name)-2]
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return err
	}
	var password string
	fmt.Print("Пароль клиента: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return err
	}

	var balance uint64
	fmt.Print("Баланс клиента: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}
	var balanceNumber uint64
	fmt.Print("Номер счёт клиента: ")
	_, err = fmt.Scan(&balanceNumber)
	if err != nil {
		return err
	}
	var phoneNumber int64
	fmt.Print("Номер телефон клиента: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}
	err = core.AddClients(core.Client{
		Id:            0,
		Name:          name,
		Login:         login,
		Password:      password,
		Balance:       balance,
		BalanceNumber: balanceNumber,
		PhoneNumber:   phoneNumber,
	}, db)
	if err != nil {
		fmt.Println("Такой логин или пароль уже существует")
		return err
	}
	fmt.Println("Пользователь успешно добавлен!")
	return nil
}
func addAtm(db *sql.DB) (err error) {
	fmt.Println("Введите данные банкомата")
	var name string
	fmt.Print("Имя бокамата: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var street string
	fmt.Print("Где находится банкомат: ")
	_, err = fmt.Scan(&street)
	if err != nil {
		return err
	}

	err = core.AddAtm(core.Atm{
		Id:     0,
		Name:   name,
		Street: street,
	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Банкомат успешно добавлен!")
	return nil
}

func addServices(db *sql.DB) (err error) {
	var name string
	fmt.Print("Название услиги: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}

	var price int64
	fmt.Print("Стоимость услуги: ")
	_, err = fmt.Scan(&price)
	if err != nil {
		return err
	}
	err = core.AddServices(core.Services{
		Id:    0,
		Name:  name,
		Price: price,
	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Услуга успешно добавлена!")
	return nil
}

func updateBalance(db *sql.DB) (err error) {
	fmt.Println("Введите данные клиента")
	var login string
	fmt.Print("Введите логин клиента: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return err
	}

	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.UpdateBalance(core.Client{
		Login:         login,
		Balance:       balance,
	}, db)
	if err != nil {
		return err
	}
	fmt.Println("Счет клиента успешно добавлен!")
	return nil
}
