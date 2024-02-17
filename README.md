# P2P-Gossips-Network

A decentralized communication protocol where nodes share information through randomized peer-to-peer exchanges, promoting efficient and fault-tolerant message dissemination.

## Installation

1) Install Go from [GO](https://golang.org/) according to your PC's specifications.
2) Add the path of the Go binary directory to the system variables.
3) To verify if Go is installed correctly, type go version in the terminal.

```bash
go version
```


## Folder Structure 
There are 4 folders in the root directory and 4 files. There is one config folder which contains config.yaml file and there we have specifid the seed port and the ip address. 
To run the seed nodes follow the below path:

1) Go the cmd folder from the root directory
```bash
cd cmd
```
2) Go to the seed folder under the cmd 
```bash
cd seed
```
3) Now run the command
```bash
go run main.go
```
All the seeds are now runing on the specified port and ip addresses and this you will verify in your terminal also.

To run the peer nodes follow the below path:

1) Go the cmd folder from the root directory
```bash
cd cmd
```
2) Go to the seed folder under the cmd 
```bash
cd peer
```
3) Now run the command 
```bash
go run main.go <ip address> <port>
```
Note: Here you have to specify some other ip address and port which are not in the seed nodes ip address and port.
Now you can see that one peer node is connected to (n/2)+1 seed nodes and if you want to connect more peer nodes you can run the above command again but in different terminal.