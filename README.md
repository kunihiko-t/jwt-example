# JWT Example

## docker-compose

```
docker-compose up
```

Access to  http://localhost:8080/


## k8s


```
# Edit hostPath of k8s/app-deployment.yaml
minikube start
kubectl apply -f ./k8s
minikube service app --url
#  http://192.168.99.101:32453
minikube stop
```