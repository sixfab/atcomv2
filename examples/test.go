package main

import (
	"fmt"

	atcom "github.com/sixfab/atcomv2"
)

func main() {

	at := atcom.NewAtcom(nil, nil)

	port, err := at.DecidePort()

	if err != nil {
		fmt.Println(err)
	}

	args := map[string]interface{}{
		"port":    port["port"],
		"baud":    115200,
		"lineEnd": true,
		"timeout": 5,
	}

	response, err := at.SendAT("ATE1", args)
	fmt.Println(response, err)

	response1, err := at.SendAT("AT+CGDCONT?", args)
	fmt.Println(response1, err)

	response2, err := at.SendAT("AT+COPS?", args)
	fmt.Println(response2, err)

	// response3, err := at.SendAT("AT#ECMD=0", args)
	// fmt.Println(response3, err)

	// response4, err := at.SendAT("AT#ECM=1,0", args)
	// fmt.Println(response4, err)

	// response5, err := at.SendAT("AT#ECM?", args)
	// fmt.Println(response5, err)

	// command that takes very long time to respond
	// response6, err := at.SendAT("AT+COPS=?", args)
	// fmt.Println(response6, err)

}
