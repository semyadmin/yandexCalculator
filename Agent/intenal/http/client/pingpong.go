package client

import (
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/adminsemy/yandexCalculator/Agent/intenal/config"
)

func Ping(conf *config.Config) {
	go func() {
		goroutines := strconv.FormatInt(int64(conf.MaxGoroutines), 10)
		for {
			address := conf.Host + ":" + conf.Port
			conn, err := net.Dial("tcp", address)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			_, err = conn.Write([]byte("ping"))
			if err != nil {
				conn.Close()
				continue
			}
			buf := make([]byte, 512)
			n, err := conn.Read(buf)
			if !errors.Is(io.EOF, err) && err != nil {
				conn.Close()
				break
			}
			if string(buf[:n]) != "pong" {
				conn.Close()
				continue
			}
			for {
				working := strconv.FormatInt(conf.WorkGoroutines.Load(), 10)
				_, err = conn.Write([]byte(goroutines + " " + working))
				if err != nil {
					break
				}
			}
			conn.Close()
		}
	}()
}
