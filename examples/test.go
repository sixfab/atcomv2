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

	response2, err := at.SendAT("AT+COPS?", args)
	fmt.Println(response2, err)

	args["desired"] = []string{"internet"}
	response1, err := at.SendAT("AT+CGDCONT?", args)
	fmt.Println(response1, err)
}
