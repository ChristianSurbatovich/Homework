package main

import (
	"net"
	"os"
	"log"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"fmt"
	"bufio"
)

func sendMessage(action string, stock string, amount int, price float64, conn net.Conn) bool{
	e, err := conn.Write([]byte(action + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
	log.Println("Sent a message: " + action + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price, 'f', 2, 64))
	log.Println(e)
	log.Println(err)
	if(err != nil){
		log.Println(err)
		return false
	}
	return true
}

func priceList(conn net.Conn, list *[]float64){
	_, err := conn.Write([]byte("PRICE" + " " + "0" + " " + "0" + " " + "0" + " \n"))
	log.Println("Sent a message: " + "PRICE" + " " + "0" + " " + "0" + " " + "0")
	if(err != nil){
		log.Println(err)
	}
	conn.SetReadDeadline(time.Time{})
	prices, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil{
		log.Fatalln(err)
	}

	s := strings.Split(prices," ")
	log.Println("Received a message: " + prices)
	s = s[1:len(s) - 1]
	for i, e := range s{
		(*list)[i], err = strconv.ParseFloat(e,64)
	}
}

func screenDisplay(channel1 chan float64, channel2 chan []float64, channel3 chan []int, stocklist []string){
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	var revenue float64
	var prices []float64
	var amount []int
	for {
		select {
			case revenue = <-channel1:  // update revenue when a new value is available
			case prices = <-channel2:
			case amount = <-channel3:

			case _ = <-ticker.C:
			// every 10 seconds display it
				s := strconv.FormatFloat(revenue, 'f', 2, 64)
				fmt.Println("\n**************************************")
				fmt.Println("This client currently has a revenue of:" + s)
				fmt.Println("Current stock holdings and prices")
				for i, e:= range stocklist{
					fmt.Println(e + " " + strconv.Itoa(amount[i]) + " " + strconv.FormatFloat(prices[i],'f',2,64))
				}
				fmt.Println("\n**************************************")
				//fmt.Println(stocklist)
				//fmt.Println(amount)
				//fmt.Println(prices)
		}
	}

}

func main(){
	idle := true
	var args []string
	var amount, stock int
	stocks := []string{"MSFT", "GOOGL", "AAPL", "GE", "C", "AMD"}
	var address string
	var x, y, z float64
	var err error
	wait := false  // for simulating wait after a transaction succeeds
	waitLonger := false

	oldprice := [6]float64{83.95, 1021.57, 170.58, 17.91, 75.17, 10.66}
	newprice := []float64{83.95, 1021.57, 170.58, 17.91, 75.17, 10.66}
	// amount of each stock this client has
	available := []int{100, 100, 100, 100, 100, 100}
	var delta [6]float64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	revenue := 0.0
	buy := false
	if r.Intn(2) > 0{
		buy = true
	}
	channelr := make(chan float64, 100)
	channela := make(chan []int, 100)
	channelp := make(chan []float64, 100)
	go screenDisplay(channelr, channelp, channela, stocks)

	channelr <- 0.0
	channela <- available
	channelp <- newprice
	if len(os.Args) > 1{
		address = os.Args[1]
	}else{
		address = "45.55.7.24:1732"
	}

	if len(os.Args) > 2{
		args = os.Args[2:]
		log.Println(args)
	}


	switch len(args){
	case 0:
		x = 5.0
		y = -5.0
		z = 50.0
	case 1:
		x, err = strconv.ParseFloat(args[0],64)
		if err != nil{
			log.Println(err)
			x = 5.0
		}
		y = -5.0
		z = 50.0
	case 2:
		x, err = strconv.ParseFloat(args[0],64)
		if err != nil{
			log.Println(err)
			x = 5.0
		}
		y, err = strconv.ParseFloat(args[1],64)
		if err != nil{
			log.Println(err)
			y = -5.0
		}
		z = 50.0
	case 3:
		x, err = strconv.ParseFloat(args[0],64)
		if err != nil{
			log.Println(err)
			x = 5.0
		}
		y, err = strconv.ParseFloat(args[1],64)
		if err != nil{
			log.Println(err)
			y = -5.0
		}
		z, err = strconv.ParseFloat(args[2],64)
		if err != nil{
			log.Println(err)
			z = 60.0
		}
	}

	serverConn, err := net.Dial("tcp",address)
	if err != nil{
		log.Fatalln(err)
	}
	defer serverConn.Close()

	for {

		if idle {

			switch buy{
			case true:
				if float64(r.Intn(100))  < z {
					stock = r.Intn(6)
					amount = r.Intn(10) + 1
					if sendMessage("BUY", stocks[stock], amount, 0, serverConn) {
						// if the message was successfully sent to the server than we wait for a response and indicate that the next action should be a sale
						// otherwise the client will attempt another buy action
						idle = false
					}

				}
				buy = false
			case false:
				// get latest pricelist from the server
				priceList(serverConn,&newprice)
				channelp <- newprice
				// get the % change of each price
				for i, _ := range newprice{
					delta[i] += (newprice[i] - oldprice[i]) / oldprice[i] * 100.0
					oldprice[i] = newprice[i]
				}

				// chose an available stock at random
				stock = r.Intn(6)

				for available[stock] < 1 {
					stock = r.Intn(6)
				}


				if (x <= delta[stock] || delta[stock] <= y) && available[stock] > 0{
					amount = r.Intn(available[stock]) + 1
					if sendMessage("SELL", stocks[stock], amount, oldprice[stock], serverConn){
						idle = false
					}
				}

				buy = true // the next transaction should be a buy


			}


		}else {
			// wait for a response for 5 seconds, after that cancel the transaction and start a new one
			// to prevent getting stuck in a state where everyone is trying to buy/sell different stocks
			for amount > 0 {
				if (buy) {
					if waitLonger{
						serverConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
					}else{
						serverConn.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
					}

				}else {
					if waitLonger{
						serverConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
					}else{
						serverConn.SetReadDeadline(time.Now().Add(40 * time.Millisecond))
					}
				}

				s, err := bufio.NewReader(serverConn).ReadString('\n')
				if (err != nil) {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					//	log.Println(err)
						waitLonger = true
						if(buy){
							sendMessage("CANCELSELL", stocks[stock], amount, 0.0, serverConn)
						}else{
							sendMessage("CANCELBUY", stocks[stock], amount, 0.0, serverConn)
						}
					}else {
						log.Fatalln(err)
					}

				}else {
					waitLonger = false
					s = strings.TrimRight(s,"\n")
					message := strings.Split(s, " ")
					log.Println("Received a message: " + s + "\n")
					if len(message) > 3 {
						for i, e := range stocks{
							if e == message[1]{
								stock = i
							}
						}
						transactionAmount, err := strconv.Atoi(message[2])
						if (err != nil) {
							log.Println(err)
							log.Println(message)
						}
						price, err := strconv.ParseFloat(message[3], 64)
						if err != nil{
							log.Println(err)
						}
						switch message[0] {
						case "BUY":
							revenue -= float64(transactionAmount) * price
							available[stock] += transactionAmount
							amount -= transactionAmount
							channelr <- revenue
							channela <- available
							wait = true

						case "SELL":
							revenue += float64(transactionAmount) * price
							available[stock] -= transactionAmount
							amount -= transactionAmount
							channelr <- revenue
							channela <- available
							delta[stock] = 0.0 // reset the price change for that stock
							wait = true;

						case "CANCELBUY":
							amount -= transactionAmount
						case "CANCELSELL":
							amount -= transactionAmount
						}
					}

				}
			}
			idle = true // ready to buy/sell again
			if wait{
				<-time.NewTimer(time.Second * 1).C
				wait = false
			}
		}


	}


	/*
	recieved messages from server

	cancel stock 0 0
	the transaction failed try same transaction again

	removed stock amount 0
	the transaction was successfully canceled by user, start a new one

	bought stock number price
	the transaction succeeded with payment xxx

	sold stock number price
	the transaction succeeded with revenue xxx

	PRICE
	an updated copy of pricelisting

	 */

}