package utils

import "fmt"

func Recover() {
	if r := recover(); r != nil {
		// log error
		fmt.Println(r)
	}
}
