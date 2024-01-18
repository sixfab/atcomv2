# atcomv2 
AT communication library and cli tool for cellular modules.

# What does it do?
This library is intended to be used with cellular modules that support AT commands. It provides a simple interface to send and receive AT commands and parse the responses. It also provides a cli tool to send AT commands to the module and receive the responses on the terminal. Both library and cli tool have auto detection of the serial port of the supported cellular modules.

# Supported modules
Listed in [modems.go](https://github.com/sixfab/atcomv2/blob/master/modems.go) file in the library.

# Installation
```
go get github.com/sixfab/atcomv2
```

# Usage
## Library
Run the example code.

```
cd examples
go run test.go
```

## CLI Tool
Build the cli tool.

```
cd pkg/cli
go build -o atcom
```

Run the cli tool. 

List the available commands.
```
./atcom -h
```

Send AT command to the module.
```
./atcom AT
```

Send AT command to the module and wait for the desired response for 5 seconds.
```
./atcom AT+CREG? -d "+CREG: 0,1" -t 5
```
