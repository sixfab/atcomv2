/*
Created by: Yasin Kaya (selengalp), yasinkaya.121@gmail.com, 2023

Copyright (c) 2023 Sixfab Inc.
*/
package main

import (
	"fmt"
)

func main() {
	//cmd.Execute()

	args := make(map[string]interface{})
	args["port"] = "/dev/ttyUSB3"
	args["baud"] = 115200
	args["desired"] = []string{"+CREG: 0,1", "+CREG: 0,5"}
	args["fault"] = []string{"+CREG: 0,2", "+CREG: 0,3", "+CREG: 0,4"}
	args["lineEnd"] = true
	args["timeout"] = 5

	res, err := sendAT("AT+CREG?", args)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
