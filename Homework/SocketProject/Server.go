package main

import (
	"net"
	"log"
	"strconv"
	"strings"
	"sync"
	"math/rand"
	"time"
	"bufio"
	"io"
)

//available := make(map[string]int)
// a map of a stock being sold and the connection selling it
var sellList []clientInfo
var muxSell sync.Mutex
// a map of the stock being bought and the connection trying to buy it
var buyList []clientInfo
var muxBuy sync.Mutex

var salesMade = make(chan int16)

type clientInfo struct{
	client net.Conn
	amount int
	stock string
}




// array of stock price
var pricelist = [6]float64{113.95, 221.57, 170.58, 68.91, 75.17, 30.66}  // starting prices from CNN money
var muxPrice sync.Mutex
var stocks = []string{"MSFT", "GOOGL", "AAPL", "GE", "C", "AMD"}

func getPrice(stock string) float64{
	for i, e := range stocks{
		if e == stock{
			muxPrice.Lock()
			p := pricelist[i]
			muxPrice.Unlock()
			return p
		}
	}
	return 0.0
}
func handleConnection(client net.Conn){
	defer client.Close()
	reader := bufio.NewReader(client)
	log.Println("Received a connection from: " + client.RemoteAddr().String())
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		s, err := reader.ReadString('\n')
		//log.Println("Received from: " + client.RemoteAddr().String() + " " + s)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && (netErr.Timeout()) {
				log.Println(err)
				continue
			}
			if err == io.EOF{
				log.Println(err)
				log.Println(s)
				continue
			}
			// if a client disconnects then remove any stocks they had listed
			log.Println(err)
			muxBuy.Lock()
			for i, e := range buyList{
				if e.client == client{
					buyList = append(buyList[:i],buyList[i+1:]...)
				}
			}
			muxBuy.Unlock()
			muxSell.Lock()
			for i, e := range sellList{
				if e.client == client{
					sellList = append(sellList[:i],sellList[i+1:]...)
				}
			}
			muxSell.Unlock()
			break
		}else{
			if strings.ContainsAny(s, ":/!@#$%^&*()_+-,`~"){
				log.Println("Illegal input")
				break
			}
			//log.Println("Received from: " + client.RemoteAddr().String() + " " + s)
			message := strings.Split(s," ")
			if len(message) > 3 {
				stock := message[1]
				amount, err := strconv.Atoi(message[2])
				if(err != nil){
					log.Println(err)
						log.Println()
				}
				price := getPrice(stock)
				switch message[0] {
				case "BUY":
					// 10% chance for transaction to fail
					if r.Intn(100) < 10{
							// the transaction failed
							_, err := client.Write([]byte("CANCELBUY" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							continue
					}
					var remove []bool
					// first, see if anyone is selling that stock
					muxSell.Lock()
					for i, seller := range sellList{
						// let the client selling the stock know it was sold
						remove = append(remove,false)
						if seller.stock == stock && seller.amount > amount{
							_, err := seller.client.Write([]byte("SELL" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							_, err = client.Write([]byte("BUY" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							sellList[i].amount -= amount
							amount = 0
							salesMade<- 1
							break  // no reason to keep going
						}else if seller.stock == stock && seller.amount <= amount { // if the buyer is buying more stocks than the seller is selling
							_, err := seller.client.Write([]byte("SELL" + " " + stock + " " + strconv.Itoa(seller.amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							_, err = client.Write([]byte("BUY" + " " + stock + " " + strconv.Itoa(seller.amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							amount -= seller.amount
							remove[i] = true
							salesMade<- 1
						}
					}
					shift := 0
					for i, b := range remove{
						if b {
							sellList = append(sellList[:i - shift],sellList[i+1 - shift:]...) //remove the stock from the sell list
							shift ++
						}
						remove[i] = false
					}
					muxSell.Unlock()
					if amount > 0 {
						muxBuy.Lock()
						buyList = append(buyList,clientInfo{
							client,
							amount,
							stock,
						})
						muxBuy.Unlock()
					}
					//log.Print(buyList)
					// if not then add it to the list of stocks people are buying


				case "SELL":
					if r.Intn(100) < 10 {
						// the transaction failed
						_, err := client.Write([]byte("CANCELSELL" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
						if err != nil{
							log.Println(err)
							break
						}
						continue
					}
					var remove []bool
					// first, see if anyone is buying that stock
					muxBuy.Lock()
					for i, buyer := range buyList{
						remove = append(remove,false)
						if buyer.stock == stock && buyer.amount > amount{
							_, err := buyer.client.Write([]byte("BUY" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							// let the buyer know they succsefully bought the stock
							_, err = client.Write([]byte("SELL" + " " + stock + " " + strconv.Itoa(amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							buyList[i].amount -= amount
							amount = 0
							salesMade<- 1
							break
						}else if buyer.stock == stock && buyer.amount <= amount { // if the buyer is buying more stocks than the seller is selling
							_, err := buyer.client.Write([]byte("BUY" + " " + stock + " " + strconv.Itoa(buyer.amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							// let the buyer know they succsefully bought the stock
							_, err = client.Write([]byte("SELL" + " " + stock + " " + strconv.Itoa(buyer.amount) + " " + strconv.FormatFloat(price, 'f', 2, 64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							amount -= buyer.amount
							remove[i] = true
							salesMade<- 1
						}
					}
					shift := 0
					for i, b := range remove{
						if b {
							buyList = append(buyList[:i - shift],buyList[i+1 - shift:]...) //remove the stock from the sell list
							shift++
						}
						remove[i] = false
					}
					muxBuy.Unlock()

					if amount > 0 {
						muxSell.Lock()
						sellList = append(sellList,clientInfo{
							client,
							amount,
							stock,
						})
						muxSell.Unlock()
					}
					// if not then add it to the list of stocks people are selling

				case "CANCELBUY":
					muxBuy.Lock()
					for i, buyer := range buyList{
						if buyer.stock == stock && buyer.client == client{
							_, err := client.Write([]byte("CANCELBUY" + " " + buyer.stock + " " + strconv.Itoa(buyer.amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							buyList = append(buyList[:i],buyList[i+1:]...) //remove the stock from the buy list
							break
						}
					}
					muxBuy.Unlock()
				case "CANCELSELL": // same thing as cancel buy just looking in the sellers list
					muxSell.Lock()
					for i, seller := range sellList{
						if seller.stock == stock && seller.client == client{
							_, err := client.Write([]byte("CANCELSELL" + " " + seller.stock + " " + strconv.Itoa(seller.amount) + " " + strconv.FormatFloat(price,'f',2,64) + " \n"))
							if err != nil{
								log.Println(err)
								break
							}
							sellList = append(sellList[:i],sellList[i+1:]...) //remove the stock from the sell list
							break
						}
					}
					muxSell.Unlock()
				case "PRICE":
					var prices []byte
					prices = append(prices,"PRICE"...)
					muxPrice.Lock()
					for _, e := range pricelist{
						prices = append(prices," " + strconv.FormatFloat(e,'f',2,64)...)
					}
					muxPrice.Unlock()
					prices = append(prices," \n"...)
					_, err := client.Write(prices)
					if err != nil{
						log.Println(err)
						break
					}
				}
			}

		}

	}
}

func priceHandler(){
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for _ = range ticker.C{
		for i, e := range pricelist{
			// ajdust each price by -5% to +5%
			muxPrice.Lock()
			pricelist[i] = e * (0.95 + (float64(r.Intn(11))/100.0))
			muxPrice.Unlock()
		}
	}

}

func listenFunc(){
	listener, err := net.Listen("tcp", ":1732")
	if err != nil {
		log.Println(err)
		log.Fatalln("Failed to start server")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
		}else {
			go handleConnection(conn)
		}
	}
}


func main(){
	// start the listener as a goroutine which will run concurrently
	go listenFunc()
	// the number of each stock available
	go priceHandler()
	var sales int
	ticker := time.NewTicker(1 * time.Minute)
	sales = 0
	for {
		select{
			case _ = <-salesMade:
				sales++
			case _ = <-ticker.C:
				log.Println(strconv.Itoa(sales) + " transactions completed")
		}
	}

}