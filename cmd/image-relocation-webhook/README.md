# Image Relocation MutatingAdmissionWebhook

This webhook supports rewriting kubernetes pods to use relocated image references. The mapping from original
to relocated image references is built by deploying `imagemap` custom resources.

## How to use the webhook

The following was done using a GKE cluster.

* Install Jetstack certificate manager:
```
kubectl create namespace cert-manager && \
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true && \
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.10.0/cert-manager.yaml
```

* If you want the webhook to provide detailed logs, enable debug logging in `config/400-deployment` like this:
```
      ...
      containers:
        - name: ir
          image: github.com/pivotal/image-relocation/cmd/image-relocation-webhook
          args: ["--debug"]
      ...
```

* Install [`ko`](https://github.com/google/ko)

* Configure `ko` by setting the environment variable `$KO_DOCKER_REPO` to a suitable image repository, e.g.:
```
export KO_DOCKER_REPO=gcr.io/my-sandbox
```

* Create various resources and deploy the webhook:
```
ko apply -f config/image-relocation-webhook/
```

* Deploy a sample imagemap custom resource, after editing it to replace `<repo prefix>` with a suitable repository
prefix (e.g. `gcr.io/my-sandbox`) to which you and the cluster have access:
```
# remember to edit config/samples/webhook_v1alpha1_imagemap.yaml
kubectl apply -f config/samples/webhook_v1alpha1_imagemap.yaml
```

* Create a pod, e.g.:
```
kubectl run kubernetes-bootcamp --image=gcr.io/google-samples/kubernetes-bootcamp:v1 --port=8080
```

* Observe the relocated image in the pod, e.g.:
```
kubectl get pod kubernetes-bootcamp-xxx -oyaml
```
Output:
```
...
spec:
  containers:
  - image: <repo prefix>/kubernetes-bootcamp:v1
...
```

Note: the `image` value under `containerStatuses` may not be the relocated value. This is a [known issue](https://github.com/kubernetes/kubernetes/issues/51017) when an image has multiple references referring to it. 

* View the logs from the webhook, e.g.:
```
kubectl logs ir-webhook-xxx -n image-relocation
```

* Now tidy up:
```
kubectl delete deployment kubernetes-bootcamp
kubectl delete -f config/samples/webhook_v1alpha1_imagemap.yaml
ko delete -f config/image-relocation-webhook/
```