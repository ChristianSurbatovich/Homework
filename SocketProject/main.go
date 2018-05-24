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

var startTime = time.Now() // the time the program was started for showing how long the program has running
// simple function for sending a message to a connection, the message takes the form ACTION STOCK AMOUNT PRICE with a newline character
// designating the end of the message
func sendMessage(action string, stock string, amount int, price float64, conn net.Conn) bool{
	_, err := conn.Write([]byte(action + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
	if(err != nil){
		log.Println(err)
		return false
	}
	return true
}

// gets a list of prices from the server and puts them into the array it was given
func priceList(conn net.Conn, r *bufio.Reader, list *[6]float64){
	_, err := conn.Write([]byte("PRICE" + " " + "0" + " " + "0" + " " + "0" + " \n"))
	if(err != nil){
		log.Println(err)
	}
	conn.SetReadDeadline(time.Now().Add(450 * time.Millisecond)) // reset the timeout on the connection
	prices, err := r.ReadString('\n')
	if err != nil{
		log.Fatalln(err)
	}

	s := strings.Split(prices," ") // split the message up by spaces
	s = s[1:len(s) - 1] // get rid of "PRICE" at the front so just the actual prices are left
	for i, e := range s{
		(*list)[i], err = strconv.ParseFloat(e,64) // convert to a float and store in array
	}
}

// function to periodically display some statistics of the program on the screen
func screenDisplay(channel1 chan float64, channel2 chan [6]float64, channel3 chan []int, stocklist []string){
	ticker := time.NewTicker(time.Second * 10) // for displaying every 10 seconds
	defer ticker.Stop()
	var revenue float64
	var prices [6]float64
	var amount []int
	var assetValue float64
	for {
		select { // wait until one of the following actions is available and then do it
		case revenue = <-channel1:  // update revenue when a new value is available
		case prices = <-channel2:   // update the list of prices when a new price list is available
		case amount = <-channel3:   // update the number of stocks held when it changes

		case _ = <-ticker.C:        // every 10 seconds display the information
			newTime := time.Now().Sub(startTime)
			assetValue = 0.0
			fmt.Println("\n**************************************")
			fmt.Printf("This client has been running for %02d:%02d:%02d\n",int(newTime.Seconds() / 3600),int(newTime.Seconds() / 60),int(newTime.Seconds()) % 60)
			fmt.Println("Current stock holdings and prices")
			for i, e:= range stocklist{
				fmt.Println(e + " " + strconv.Itoa(amount[i]) + " " + strconv.FormatFloat(prices[i],'f',2,64))
				assetValue += float64(amount[i]) * prices[i]
			}
			fmt.Printf("Asset value = %0.2f		Balance = %0.2f		Net Worth = %0.2f\n",assetValue,revenue,assetValue + revenue)
			fmt.Println("\n**************************************")
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
	var oldprice  [6]float64
	var newprice  [6]float64
	// amount of each stock this client has
	available := []int{100, 100, 100, 100, 100, 100}
	var delta [6]float64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	revenue := 0.0
	// choose whether the first transaction will be a buy or sell
	buy := false
	if r.Intn(2) > 0{
		buy = true
	}

	// Get the command line arguments for server address, X, Y, and Z or assign the default values
	if len(os.Args) > 1{
		address = os.Args[1]
	}else{
		address = "45.55.7.24:1732"
	}

	if len(os.Args) > 2{
		args = os.Args[2:]
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


	// open the tcp socket
	serverConn, err := net.Dial("tcp",address)
	if err != nil{
		log.Fatalln(err)
	}
	defer serverConn.Close()  //
	reader := bufio.NewReader(serverConn)

	priceList(serverConn, reader, &oldprice) // fill in the price list with initial values
	priceList(serverConn, reader, &oldprice)
	// channels for communicating with the display function
	channelr := make(chan float64, 100)
	channela := make(chan []int, 100)
	channelp := make(chan [6]float64, 100)
	go screenDisplay(channelr, channelp, channela, stocks) // start the display function as a new Goroutine
	// send the display function some intial values
	channelr <- 0.0
	channela <- available
	channelp <- newprice

	for {

		if idle {
			// try and make a buy or sell decision
			switch buy{
			case true:
				if float64(r.Intn(100))  < z {
					stock = r.Intn(6)
					amount = r.Intn(10) + 1
					if sendMessage("BUY", stocks[stock], amount, 0, serverConn) {
						// if the message was sent then change idle to false so the client will wait for a response that the sale has gone through
						// otherwise the client will attempt a sell action
						idle = false
					}

				}
				buy = false
			case false:
				// get latest pricelist from the server
				priceList(serverConn, reader, &newprice)
				channelp <- newprice
				// calculate the % change of each price
				for i, _ := range newprice{
					delta[i] += (newprice[i] - oldprice[i]) / oldprice[i] * 100.0
					oldprice[i] = newprice[i]
				}

				// check if the client has any stocks to sell
				anyAvailable := false
				for _, e := range available{
					if e > 0 {
						anyAvailable = true
					}
				}
				// if it does than pick one at random
				if anyAvailable{
					for available[stock] < 1 {
						stock = r.Intn(6)
					}
				}else{ // if there are no stocks to sell than try to buy some
					buy = true
					continue
				}

				// check if the chosen stock price has increased by atleast x or decreased by atleast y
				if (x <= delta[stock] || delta[stock] <= y) && available[stock] > 0{
					amount = r.Intn(available[stock]) + 1
					if sendMessage("SELL", stocks[stock], amount, oldprice[stock], serverConn){
						idle = false
					}
				}

				buy = true // the next transaction should be a buy


			}


		}else {
			// read a response from the serve, if the read operation times out then cancel the transaction and start a new one
			// to prevent getting stuck in a state where everyone is trying to buy/sell different stocks and no matches can be made
			// because I am doing the bonus where clients have to buy the stock from a different client
			for amount > 0 {
				if (buy) {
					// set a deadline for the read operation, if nothing is read by the deadline than it will return a timeout error
					serverConn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))


				}else {
					serverConn.SetReadDeadline(time.Now().Add(30 * time.Millisecond))

				}

				s, err := reader.ReadString('\n') // read a string from the server in the form ACTION STOCK AMOUNT PRICE
				if (err != nil) {
					// if it is a timeout error it means no one was selling/buying the same stock so cancel the current transaction
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						//	log.Println(err)
						if(buy){
							sendMessage("CANCELSELL", stocks[stock], amount, 0.0, serverConn)
						}else{
							sendMessage("CANCELBUY", stocks[stock], amount, 0.0, serverConn)
						}
					}else {
						log.Fatalln(err)
					}

				}else {
					// seperate the string by spaces
					s = strings.TrimRight(s,"\n")
					message := strings.Split(s, " ")
					if message[0] == "PRICE"{
						break
					}
					//log.Println("Received a message: " + s + "\n")
					if len(message) > 3 {
						// getting index of the stock that was bought or sold
						for i, e := range stocks{
							if e == message[1]{
								stock = i
							}
						}
						// the number of shares
						transactionAmount, err := strconv.Atoi(message[2])
						if (err != nil) {
							log.Println(err)
							log.Println(message)
						}
						// the price per share
						price, err := strconv.ParseFloat(message[3], 64)
						if err != nil{
							log.Println(err)
						}
						switch message[0] {
						case "BUY":
							// subtract the cost of the transaction from revenue, increase the amount of available stock
							// also keep track of how many shares were actually bought in the transaction compared to how many the client was trying to buy
							// because the other client might have been selling fewer shares so the buy request might get matched up to different sellers
							revenue -= float64(transactionAmount) * price
							available[stock] += transactionAmount
							amount -= transactionAmount  // the remaining amount of stocks this client is trying to purchas
							channelr <- revenue // update the display function
							channela <- available
							wait = true  // short wait after buying/selling

						case "SELL":
							// same as for a buy transactioin but revenue is increased and stock amount decreased
							// also the % price change for the stock that was sold is reset
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
			if wait{ // short waiting period if a transaction was made
				<-time.NewTimer(time.Second * 1).C
				wait = false
			}
		}


	}
}