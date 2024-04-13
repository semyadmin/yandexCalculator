package main

import (
	"fmt"
	"sync"
)

func SelectUsers(in, out chan interface{}) {
	defer close(out) // закрываем канал out после того, как все горутины отработали
	wg := sync.WaitGroup{}
	for email := range in { //  пока поступают email-ы
		wg.Add(1)
		go func(email interface{}, out chan<- interface{}) { //  запускаем обработчики
			defer wg.Done()
			select {
			case out <- GetUser(email.(string)): //  делаем запрос в селекте, чтобы не зависали горутины
			}
		}(email, out)
	}
	wg.Wait()
}

func GetUser(email string) interface{} {
	return email
}

func main() {
	emails := make(chan interface{})
	users := make(chan interface{})
	go SelectUsers(emails, users)

	emails <- "a@a.ru"
	emails <- "b@b.ru"
	emails <- "c@c.ru"
	emails <- "d@d.ru"
	emails <- "e@e.ru"
	close(emails)
	for user := range users {
		fmt.Println(user)
	}
}
