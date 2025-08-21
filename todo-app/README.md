# Todo app backend

First navigate to the `todo-app` directory:

```shell
cd todo-app
```

To deploy with kubectl:

```shell
kubectl apply -f manifests/backend-deployment.yaml
```

To access the backend, you can create a service:

```shell
kubectl apply -f manifests/backend-service.yaml
```