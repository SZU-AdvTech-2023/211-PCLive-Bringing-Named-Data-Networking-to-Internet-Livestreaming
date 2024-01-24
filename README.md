# Build PCLive_reproduction project

```
cd <PCLive-dir>
go build -o run_PCLive
```

# Initial configuration

## Configure the host

run an Simple Realtime Server(SRS) for convertion from RTMP to HLS

push an RTMP video streaming to SRS

then execute the following command in the host:
```
	run_PCLive -d
```

## Configure linux containers

install necessary libraries

```
apt update
apt install nfd
apt install ndnchunks
apt install ndnpeek
apt install ndnping
apt install ndndump
apt install ndn-dissect
```

connect your containers with NDN links

you can execute the following command in a terminal to add a route toward another container

```
nfdc route add /ndn udp://<other-host>
```

# Running

Each of the containers should run the following command to act as a producer:

```
run_PCLive
```

You can obtain data via NDN by executing the following command:

```
run_PCLive -i <interest-name>
```
