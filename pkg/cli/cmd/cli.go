package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	atcom "github.com/sixfab/atcomv2"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", atcom.Version)
	},
}

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect the modem port",
	Long:  `Detect the modem port`,
	Run: func(cmd *cobra.Command, args []string) {

		allFlag := cmd.Flag("all").Value.String()
		vidFlag := cmd.Flag("vid").Value.String()
		pidFlag := cmd.Flag("pid").Value.String()
		portFlag := cmd.Flag("port").Value.String()

		at := atcom.NewAtcom(nil, nil, nil)

		modem, err := at.DecidePort()

		if err != nil {
			fmt.Println(err)
		}

		if allFlag == "true" {
			for key, value := range modem {
				fmt.Println(key + ":" + value)
			}
		}

		if vidFlag == "true" {
			fmt.Println(modem["vendor_id"])
		}

		if pidFlag == "true" {
			fmt.Println(modem["product_id"])
		}

		if portFlag == "true" {
			fmt.Println(modem["port"])
		}

		if allFlag == "false" && vidFlag == "false" && pidFlag == "false" && portFlag == "false" {
			fmt.Println(modem["port"])
		}
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atcomv2cli [AT] [flags]",
	Short: "AT Command CLI",
	Long:  `AT Command CLI for communicating cellular modems with AT commands.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		command := args[0]
		port := cmd.Flag("port").Value.String()
		baud := cmd.Flag("baud").Value.String()
		desired := cmd.Flag("desired").Value.String()
		fault := cmd.Flag("fault").Value.String()
		timeout := cmd.Flag("timeout").Value.String()
		lineend := cmd.Flag("lineend").Value.String()
		verbose := cmd.Flag("verbose").Value.String()

		// convert parameters to suitable format with library
		baudInt, _ := strconv.Atoi(baud)
		timeoutInt, _ := strconv.Atoi(timeout)
		lineendBool, _ := strconv.ParseBool(lineend)

		desiredSlice := []string{}
		faultSlice := []string{}

		if desired != "" {
			words := strings.Split(desired, "/")
			desiredSlice = append(desiredSlice, words...)
		} else {
			desiredSlice = nil
		}

		if fault != "" {
			words := strings.Split(fault, "/")
			faultSlice = append(faultSlice, words...)
		} else {
			faultSlice = nil
		}

		// prepare command arguments
		comArgs := map[string]interface{}{
			"port":    port,
			"baud":    baudInt,
			"desired": desiredSlice,
			"fault":   faultSlice,
			"timeout": timeoutInt,
			"lineEnd": lineendBool,
		}

		// if verbose mode is true, print parameters
		if verbose == "true" {
			fmt.Println("--------------------------------------")
			fmt.Println("Command: ", command)
			fmt.Println("Port: ", port)
			fmt.Println("Baud: ", baud)
			fmt.Println("Desired: ", desiredSlice)
			fmt.Println("Fault: ", faultSlice)
			fmt.Println("Timeout: ", timeout)
			fmt.Println("Verbose: ", verbose)
			fmt.Println("--------------------------------------")
			fmt.Println("")
		}

		at := atcom.NewAtcom(nil, nil, nil)

		response, err := at.SendAT(command, comArgs)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, res := range response {
			fmt.Println(res)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.atcomv2cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("port", "p", "/dev/ttyUSB2", "port name")
	rootCmd.Flags().IntP("baud", "b", 115200, "baud rate")
	rootCmd.Flags().StringP("desired", "d", "", "desired responses - separate your multiple words with /")
	rootCmd.Flags().StringP("fault", "f", "", "fault responses - separate your multiple words with /")
	rootCmd.Flags().IntP("timeout", "t", 5, "timeout duration in seconds")
	rootCmd.Flags().BoolP("lineend", "l", true, "line end")
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose mode")
	rootCmd.Flags().StringP("version", "V", "", "version")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(detectCmd)

	detectCmd.Flags().BoolP("all", "a", false, "all modem attributes")
	detectCmd.Flags().BoolP("vid", "v", false, "vendor id")
	detectCmd.Flags().BoolP("pid", "i", false, "product id")
	detectCmd.Flags().BoolP("port", "p", false, "serial port")
}
