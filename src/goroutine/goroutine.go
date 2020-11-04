package main

import (
	"fmt"
	"time"
)

	func oneSecond(){
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("1 second")
	}

	func twoSecond(){
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("2 second")
	}

	func threeSecond(){
		time.Sleep(3000 * time.Millisecond)
		fmt.Println("3 second")
	}

func main() {
	fmt.Println(time.Now())
	 threeSecond()
	go oneSecond()
	go twoSecond()
	 fmt.Println(time.Now())
}