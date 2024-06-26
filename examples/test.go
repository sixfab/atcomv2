package main

import (
	"fmt"

	atcom "github.com/sixfab/atcomv2"
)

func main() {

	at := atcom.NewAtcom(nil, nil)

	_, err := at.DecidePort()

	if err != nil {
		fmt.Println(err)
		return
	}

	com := atcom.NewATCommand("AT+CGDCONT?")
	com.Desired = []string{"super"}
	com = at.SendAT(com)
	com.GetMeaningfulPart("+CGDCONT: ")

	if com.Error != nil {
		fmt.Println(com.Error)
	}

	fmt.Println("Command: ", com.Command)
	fmt.Println("Response: ", com.Response)
	fmt.Println("Processed: ", com.Processed)
	fmt.Println("Error: ", com.Error)
	fmt.Println("Desired: ", com.Desired)
	fmt.Println("Fault: ", com.Fault)
	fmt.Println("Timeout: ", com.Timeout)
	fmt.Println("LineEnd: ", com.LineEnd)
}
