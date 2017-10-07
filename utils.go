package main

import "fmt"

func Check(e error) {
	if e != nil {
		fmt.Println("Got error:", e)
	}
}
