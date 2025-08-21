# Todo app backend

First navigate to the `todo-app` directory:

```shell
cd todo-app
```

To deploy with kubectl:

```shell
kubectl apply -f manifests/backend-deployment.yaml
```

To access the backend, you can use port forwarding:

```shell
kubectl port-forward <todo-backend-pod-name> 3000:3000
```