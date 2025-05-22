package main

import (
	"eth-practice/tx"
	"fmt"
)

func main() {

	tx.BlockSubscribe()
	fmt.Println("Hello World")
}
