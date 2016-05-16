package main

import (
	"fmt"
	"github.com/actionpay/postmanq/common"
	"github.com/actionpay/postmanq/logger"
	"log"
	"net"
	"net/smtp"
	"runtime"
)

func main() {
	common.DefaultWorkersCount = runtime.NumCPU()
	logger.Inst()

	logger.By("localhost").Info("start!")
	tcpAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort("localhost", "0"))
	if err == nil {
		logger.By("localhost").Info("resolve tcp addr localhost")
		dialer := &net.Dialer{
			LocalAddr: tcpAddr,
			DualStack: true,
			//Timeout:   time.Second * 30,
		}
		hostname := net.JoinHostPort("localhost", "2225")
		connection, err := dialer.Dial("tcp", hostname)
		if err == nil {
			logger.By("localhost").Info("dial localhost:2225")
			c, err := smtp.NewClient(connection, "example.com")
			//c, err := smtp.Dial(hostname)
			//if err != nil {
			//	log.Fatal(err)
			//}
			if err == nil {
				// Set the sender and recipient first
				if err := c.Mail("sender@example.org"); err != nil {
					log.Fatal("Mail", " ", err)
				}
				if err := c.Rcpt("recipient@example.net"); err != nil {
					log.Fatal("Rcpt", " ", err)
				}

				// Send the email body.
				wc, err := c.Data()
				if err != nil {
					log.Fatal("Data", " ", err)
				}
				_, err = fmt.Fprintf(wc, "This is the email body")
				if err != nil {
					log.Fatal("Fprintf", " ", err)
				}
				err = wc.Close()
				if err != nil {
					log.Fatal("Close", " ", err)
				}

				// Send the QUIT command and close the connection.
				err = c.Quit()
				if err != nil {
					log.Fatal("Quit", " ", err)
				}

				log.Println("success!")
			} else {
				logger.By("localhost").Info("can't create client")
			}
		} else {
			logger.By("localhost").Info("can't dial localhost:2225")
		}
	} else {
		logger.By("localhost").Info("can't resolve tcp addr localhost")
	}
}