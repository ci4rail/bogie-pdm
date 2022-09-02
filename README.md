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


# Developing with edgefarm

To avoid long docker build time:

Define app manifest so that application will not be started by default:
```
  ...
  components:
    - name: bogie-edge
      type: edge-worker
      properties:
        runtime:
          - bogie-mc
        image: ci4rail/bogie-edge:a4
        name: bogie-edge-container
        #args: ["--config=/config/.bogie-edge-config.yaml"]
        command: ["sh", "-c", "echo Hello && sleep 10000000"]
```


Build static app
```
$ make bogie-edge-static && scp bin/bogie-edge-static root@192.168.24.110:~
```

After applying the manifest:

On edge device
```
docker cp bogie-edge-static 84f3e32205d2:/
```

On devpc:
```
/bogie-edge-static --config=/config/.bogie-edge-config.yaml
```
