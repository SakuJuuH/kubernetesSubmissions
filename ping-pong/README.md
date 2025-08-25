# Ping Pong App

First navigate to the `ping-pong` directory:

```shell
cd ping-pong
```

To deploy the application, run:

```shell
kubectl apply -f manifests -f ../log-output/manifests
```

or use the `deploy.sh` script:

App can be accessed in the browser at:

```
http://localhost:8081/pingpong
```