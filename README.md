# dist-grep

A distributed version of grep command (So called "dist-grep") implemented by CS 425 Group#29,  baohez2(Baohe Zhang) & kechenl3(Kechen Lu). Dist-grep leverages the Socket API of Golang, and follows the classic Client-Server(C/S) architecture, with fairly fault-tolerant, concurrent, configurable  and protable capability.

## Project Info

- Language: Golang 1.11

- Tested Platform: macOS 10.13.6, CentOS 7

- Code Structure:

  ​		--client: grep client side

  ​		--scripts: helper bash scripts to help build, start and manage the git repo remotely on VMs

  ​		--server: grep server side

  ​		--test: unit test cases 

  ​		--utils: internal packages of project for client

## How-to

### Build 

We made some easy-to-use shell scripts to build the program, including the server and client. Follow the command below:

```bash
./scripts/build_all
```

Or go to the server or client directory, and run `go build`

After this, under server/client directory, it has the binary executable file.

### Install and Run

First, we need to deploy the project on each VM. After **ssh** into a VM console, we can use git clone:

```bash
git clone git@gitlab.engr.illinois.edu:kechenl3/dist-grep.git`
```

By the way, `ssh-copyid` is a good way to get away from annoy password typing in. Before we run the server or client, we will have to configure the `config.json` , setting up correct VM addresses and Local VM Information. By default, server port is `5555` and log path is in `/var/log/dist-grep/vm?.log` , ? refers to 0-10. Please remember put the log file in configured path. 

So far, we are able to deploy the server on each VM. Go to **server** directory, and run `./server` . Or we can run `nohup ./server&` to safely put the server process in background. Actually we have simple script to up all servers on just one single machine or your laptop. Run like:

```bash
cd ./scripts/
./start_server_remote.sh
```

If we have something update in git repo, we could just run below to pull repo in all VMs.

```bash
cd ./scripts/
./update_vm_repo.sh
```

After we deploy the servers on each machine, we could run our dist-grep command. It resides in client directory. Run the `./client` executable file and you can use all the arguments like run `grep` in Unix/Linux. An example like below. -E and -c is just like the original grep command, for the extended reg exp pattern and show the matched lines count. We internally add the `-Hn` argument to show filename and matched line number, with a clear query results feedback.

```bash
./client -E YOUR_PATTERN -c(optional)
```

If we want to crash a server, just run a helper script, for example, crash VM9. By the way, you need modify the VM address pattern to run it successly and make sure `ssh-copyid` first.

```bash
cd ./scripts/
./kill_server_remote.sh '09'
```

### Grep Query Example

For example, if we query a infrequent pattern like `1.2{6}` , command running and output like below.

```bash
cd ~/go/src/dist-grep/client
./client -E 1.2{6} -c
Connected with fa18-cs425-g29-08.cs.illinois.edu:5555
Connected with fa18-cs425-g29-03.cs.illinois.edu:5555
Failed to connect fa18-cs425-g29-09.cs.illinois.edu:5555
Connected with fa18-cs425-g29-05.cs.illinois.edu:5555
Connected with fa18-cs425-g29-01.cs.illinois.edu:5555
Connected with fa18-cs425-g29-06.cs.illinois.edu:5555
Connected with fa18-cs425-g29-04.cs.illinois.edu:5555
Connected with fa18-cs425-g29-10.cs.illinois.edu:5555
Connected with fa18-cs425-g29-07.cs.illinois.edu:5555
Connected with fa18-cs425-g29-02.cs.illinois.edu:5555
/var/log/dist-grep/vm8.log:0
/var/log/dist-grep/vm5.log:2
/var/log/dist-grep/vm3.log:0
/var/log/dist-grep/vm6.log:0
/var/log/dist-grep/vm4.log:0
/var/log/dist-grep/vm2.log:1
/var/log/dist-grep/vm1.log:0
/var/log/dist-grep/vm7.log:1
/var/log/dist-grep/vm10.log:0
Total Connected VMs: 9
Total Counts: 4
Total Time: 0.333 seconds
```

### Unit Test

We set up four test cases to test our program, they are query of regular expression pattern, frequent pattern, infrequent pattern and frequent pattern(one server crashed). We get the locally greped file for the whole complete log first, put them in the `golden_test` as correct queried log datasets.

Then we could look into the `test/log_static` directory. `go build ` first and run it`./log_static`. The output like below could verify the program correctness. We use `diff` command internally in the program to decide if there are differences between two logs.

```bash
./log_static 
Test mode starts...
Connected with fa18-cs425-g29-04.cs.illinois.edu:5555
Connected with fa18-cs425-g29-03.cs.illinois.edu:5555
Connected with fa18-cs425-g29-10.cs.illinois.edu:5555
Connected with fa18-cs425-g29-05.cs.illinois.edu:5555
Connected with fa18-cs425-g29-01.cs.illinois.edu:5555
Connected with fa18-cs425-g29-02.cs.illinois.edu:5555
Connected with fa18-cs425-g29-09.cs.illinois.edu:5555
Connected with fa18-cs425-g29-08.cs.illinois.edu:5555
Connected with fa18-cs425-g29-06.cs.illinois.edu:5555
Connected with fa18-cs425-g29-07.cs.illinois.edu:5555
Total Connected VMs: 10
Total Counts: 52
Total Time: 0.357 seconds
Test Passed for Regular Expression Pattern: -t -E 1.2{5}

Test mode starts...
Connected with fa18-cs425-g29-01.cs.illinois.edu:5555
Connected with fa18-cs425-g29-10.cs.illinois.edu:5555
Connected with fa18-cs425-g29-03.cs.illinois.edu:5555
Connected with fa18-cs425-g29-06.cs.illinois.edu:5555
Connected with fa18-cs425-g29-08.cs.illinois.edu:5555
Connected with fa18-cs425-g29-02.cs.illinois.edu:5555
Connected with fa18-cs425-g29-04.cs.illinois.edu:5555
Connected with fa18-cs425-g29-09.cs.illinois.edu:5555
Connected with fa18-cs425-g29-05.cs.illinois.edu:5555
Connected with fa18-cs425-g29-07.cs.illinois.edu:5555
Total Connected VMs: 10
Total Counts: 6
Total Time: 0.292 seconds
Test Passed for Infrequent Pattern: -t -E 1.1{6}

Test mode starts...
Connected with fa18-cs425-g29-10.cs.illinois.edu:5555
Connected with fa18-cs425-g29-06.cs.illinois.edu:5555
Connected with fa18-cs425-g29-01.cs.illinois.edu:5555
Connected with fa18-cs425-g29-02.cs.illinois.edu:5555
Connected with fa18-cs425-g29-03.cs.illinois.edu:5555
Connected with fa18-cs425-g29-05.cs.illinois.edu:5555
Connected with fa18-cs425-g29-04.cs.illinois.edu:5555
Connected with fa18-cs425-g29-07.cs.illinois.edu:5555
Connected with fa18-cs425-g29-09.cs.illinois.edu:5555
Connected with fa18-cs425-g29-08.cs.illinois.edu:5555
Total Connected VMs: 10
Total Counts: 180047
Total Time: 24.300 seconds
Test Passed for Frequent Pattern: -t -E 2.11

Test mode starts...
Failed to connect fa18-cs425-g29-09.cs.illinois.edu:5555
Connected with fa18-cs425-g29-08.cs.illinois.edu:5555
Connected with fa18-cs425-g29-03.cs.illinois.edu:5555
Connected with fa18-cs425-g29-10.cs.illinois.edu:5555
Connected with fa18-cs425-g29-01.cs.illinois.edu:5555
Connected with fa18-cs425-g29-04.cs.illinois.edu:5555
Connected with fa18-cs425-g29-06.cs.illinois.edu:5555
Connected with fa18-cs425-g29-07.cs.illinois.edu:5555
Connected with fa18-cs425-g29-02.cs.illinois.edu:5555
Connected with fa18-cs425-g29-05.cs.illinois.edu:5555
Total Connected VMs: 9
Total Counts: 161963
Total Time: 23.598 seconds
Test Passed for Crashed One Server Pattern: -t -E 2.11
```













