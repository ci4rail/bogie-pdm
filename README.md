# bogie-pdm
Application to monitor the health of bogie axles of a train


# Test outside edgefarm

On moducop, start nats
```bash
docker run -d --restart always --name nats -p 4222:4222 -p 8222:8222 nats --http_port 8222 -js
```

Create jetstream
```bash
nats -s nats://localhost:4222 stream add test --subjects "bogie metrics" --ack --max-msgs=100000 --max-bytes=1073741824 --max-age=2d --storage file --retention limits --max-msg-size=-1 --discard old --dupe-window="0s" --replicas 1 --max-msgs-per-subject=-1
```

* Copy `cmd/bogie-edge/.bogie-edge-config.yaml` to moducop `/home/root`
* Copy app:

```bash
GOOS=linux GOARCH=arm64 go build -o ../bin/bogie . && scp ../bin/bogie root@192.168.23.159:~
```

Start bogie app on moducop
```
./bogie
```
