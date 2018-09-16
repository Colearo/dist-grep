# dist-grep

A distributed version of grep command (So called "dist-grep") implemented by CS 425 Group#29,  baohez2(Baohe Zhang) & kechenl3(Kechen Lu). Dist-grep leverages the Socket API of Golang, and follows the classic Client-Server(C/S) architecture, with fairly fault-tolerant, concurrent, configurable  and protable capability.

## Project Info

- Language: Golang 1.11

- Tested Platform: macOS 10.13.6, CentOS 7

- Code Structure:

  --client: grep client side

  --scripts: helper bash scripts to help build, start and manage the git repo remotely on VMs

  --server: grep server side

  --test: unit test cases 

  --utils: internal packages of project for client

## How-to

### Build 

We made some easy-to-use shell scripts to build the program, including the server and client. Follow the command below:

`./scripts/build_all`

Or go to the server or client directory, and run 

`cd ./server`

`go build`



