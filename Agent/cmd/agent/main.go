package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
	grpcclient "github.com/adminsemy/yandexCalculator/Agent/intenal/grpc_client"
)

func main() {
	config := config.New()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	for i := 0; i < config.MaxGoroutines; i++ {
		go func(id int) {
			for {
				client, err := grpcclient.New(context.Background(), config, uint64(id))
				if err != nil {
					break
				}
				err = client.Start()
				if err != nil {
					time.Sleep(1 * time.Second)
				}
			}
		}(i)
	}

	<-done
}
