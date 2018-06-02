package main

import "fmt"

func checkErr(err error) bool {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	return err != nil
}
