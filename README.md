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

On devpc: k exec???
```
/bogie-edge-static --config=/config/.bogie-edge-config.yaml
```

# Using edgefarm

Repo contains git-crypted secrets. To decrypt them, you need to have the key: https://vault.bitwarden.com/#/vault?itemId=3715f6a1-95b0-4e44-bf22-af8100ed2d05
Store attachment as `../bogie-pdm-git-crypt-key` and run `git-crypt unlock ../bogie-pdm-git-crypt-key`

Here: Devcluster in linode

```
export KUBECONFIG=~/.kube/linode.cfg
kubectl ns bogie
```

Deploy app
```
kubectl apply -f manifests/bogie_app.yaml
```

Delete app (Delete also if config maps are changed!)
```
kubectl delete -f manifests/bogie_app.yaml
```


# Monitoring device

Get grafana password
```
kubectl get secrets -n monitoring grafana -o jsonpath="{.data.admin-password}" | base64 -d
```

Forward grafana port

For WSL users: Won't work in WSL. Do it on windows host.
```
kubectl port-forward -n monitoring svc/grafana 8080:80
```

Open grafana in browser: http://localhost:8080
Login with user `admin` and password from above.


# Get data from nats

```bash
nats context create ngs004 --server connect.ngs.global  --creds edgefarm/manifest/ngs0004_customer.creds
nats context select ngs004
nats stream ls
```
