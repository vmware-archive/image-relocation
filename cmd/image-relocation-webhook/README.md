# Image Relocation MutatingAdmissionWebhook

## Deploying

The following was done using a GKE cluster.

* Install Jetstack certificate manager:
```
kubectl create namespace cert-manager && \
kubectl label namespace cert-manager certmanager.k8s.io/disable-validation=true && \
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.10.0/cert-manager.yaml
```

* Install [`ko`](https://github.com/google/ko)

* Configure `ko` by setting the environment variable `$KO_DOCKER_REPO` to a suitable image repository, e.g.:
```
export KO_DOCKER_REPO=gcr.io/my-sandbox
```

* Create a namespace and a service account and then deploy the webhook:
```
kubectl apply -f config/image-relocation-webhook/100-namespace.yaml && \
kubectl apply -f config/image-relocation-webhook/200-rbac.yaml && \
ko apply -f config/image-relocation-webhook/
```

* Create a pod, e.g.:
```
kubectl run kubernetes-bootcamp --image=gcr.io/google-samples/kubernetes-bootcamp:v1 --port=8080
```

* Grab the logs from the webhook, e.g.:
```
kubectl logs ir-webhook-xxx -n image-relocation
```

* Observe the hook being called, assuming debug logging is enabled in the webhook deployment:
```
      ...
      containers:
        - name: ir
          image: github.com/pivotal/image-relocation/cmd/image-relocation-webhook
          args: ["--debug"]
      ...
```

Here's an example log:
```
{"level":"info","ts":1569423113.8658712,"logger":"image-relocation.handlers","msg":"req: {\"uid\":\"038beb9a-dfa4-11e9-aeb1-42010a8402aa\",\"kind\":{\"group\":\"\",\"version\":\"v1\",\"kind\":\"Pod\"},\"resource\":{\"group\":\"\",\"version\":\"v1\",\"resource\":\"pods\"},\"namespace\":\"default\",\"operation\":\"CREATE\",\"userInfo\":{\"username\":\"system:serviceaccount:kube-system:replicaset-controller\",\"uid\":\"23537f3f-d304-11e9-a808-42010a840094\",\"groups\":[\"system:serviceaccounts\",\"system:serviceaccounts:kube-system\",\"system:authenticated\"]},\"object\":{\"kind\":\"Pod\",\"apiVersion\":\"v1\",\"metadata\":{\"generateName\":\"kubernetes-bootcamp-6bf84cb898-\",\"creationTimestamp\":null,\"labels\":{\"pod-template-hash\":\"6bf84cb898\",\"run\":\"kubernetes-bootcamp\"},\"annotations\":{\"kubernetes.io/limit-ranger\":\"LimitRanger plugin set: cpu request for container kubernetes-bootcamp\"},\"ownerReferences\":[{\"apiVersion\":\"apps/v1\",\"kind\":\"ReplicaSet\",\"name\":\"kubernetes-bootcamp-6bf84cb898\",\"uid\":\"0387ff13-dfa4-11e9-aeb1-42010a8402aa\",\"controller\":true,\"blockOwnerDeletion\":true}]},\"spec\":{\"volumes\":[{\"name\":\"default-token-vxtxr\",\"secret\":{\"secretName\":\"default-token-vxtxr\"}}],\"containers\":[{\"name\":\"kubernetes-bootcamp\",\"image\":\"gcr.io/google-samples/kubernetes-bootcamp:v1\",\"ports\":[{\"containerPort\":8080,\"protocol\":\"TCP\"}],\"resources\":{\"requests\":{\"cpu\":\"100m\"}},\"volumeMounts\":[{\"name\":\"default-token-vxtxr\",\"readOnly\":true,\"mountPath\":\"/var/run/secrets/kubernetes.io/serviceaccount\"}],\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"IfNotPresent\"}],\"restartPolicy\":\"Always\",\"terminationGracePeriodSeconds\":30,\"dnsPolicy\":\"ClusterFirst\",\"serviceAccountName\":\"default\",\"serviceAccount\":\"default\",\"securityContext\":{},\"schedulerName\":\"default-scheduler\",\"tolerations\":[{\"key\":\"node.kubernetes.io/not-ready\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300},{\"key\":\"node.kubernetes.io/unreachable\",\"operator\":\"Exists\",\"effect\":\"NoExecute\",\"tolerationSeconds\":300}],\"priority\":0,\"enableServiceLinks\":true},\"status\":{}},\"oldObject\":null,\"dryRun\":false} resp: {\"Patches\":null,\"uid\":\"\",\"allowed\":true,\"status\":{\"metadata\":{},\"reason\":\"no change - not implemented\",\"code\":200}}"}
```

which contains the following request:
```
{ 
   "uid":"038beb9a-dfa4-11e9-aeb1-42010a8402aa",
   "kind":{ 
      "group":"",
      "version":"v1",
      "kind":"Pod"
   },
   "resource":{ 
      "group":"",
      "version":"v1",
      "resource":"pods"
   },
   "namespace":"default",
   "operation":"CREATE",
   "userInfo":{ 
      "username":"system:serviceaccount:kube-system:replicaset-controller",
      "uid":"23537f3f-d304-11e9-a808-42010a840094",
      "groups":[ 
         "system:serviceaccounts",
         "system:serviceaccounts:kube-system",
         "system:authenticated"
      ]
   },
   "object":{ 
      "kind":"Pod",
      "apiVersion":"v1",
      "metadata":{ 
         "generateName":"kubernetes-bootcamp-6bf84cb898-",
         "creationTimestamp":null,
         "labels":{ 
            "pod-template-hash":"6bf84cb898",
            "run":"kubernetes-bootcamp"
         },
         "annotations":{ 
            "kubernetes.io/limit-ranger":"LimitRanger plugin set: cpu request for container kubernetes-bootcamp"
         },
         "ownerReferences":[ 
            { 
               "apiVersion":"apps/v1",
               "kind":"ReplicaSet",
               "name":"kubernetes-bootcamp-6bf84cb898",
               "uid":"0387ff13-dfa4-11e9-aeb1-42010a8402aa",
               "controller":true,
               "blockOwnerDeletion":true
            }
         ]
      },
      "spec":{ 
         "volumes":[ 
            { 
               "name":"default-token-vxtxr",
               "secret":{ 
                  "secretName":"default-token-vxtxr"
               }
            }
         ],
         "containers":[ 
            { 
               "name":"kubernetes-bootcamp",
               "image":"gcr.io/google-samples/kubernetes-bootcamp:v1",
               "ports":[ 
                  { 
                     "containerPort":8080,
                     "protocol":"TCP"
                  }
               ],
               "resources":{ 
                  "requests":{ 
                     "cpu":"100m"
                  }
               },
               "volumeMounts":[ 
                  { 
                     "name":"default-token-vxtxr",
                     "readOnly":true,
                     "mountPath":"/var/run/secrets/kubernetes.io/serviceaccount"
                  }
               ],
               "terminationMessagePath":"/dev/termination-log",
               "terminationMessagePolicy":"File",
               "imagePullPolicy":"IfNotPresent"
            }
         ],
         "restartPolicy":"Always",
         "terminationGracePeriodSeconds":30,
         "dnsPolicy":"ClusterFirst",
         "serviceAccountName":"default",
         "serviceAccount":"default",
         "securityContext":{ 

         },
         "schedulerName":"default-scheduler",
         "tolerations":[ 
            { 
               "key":"node.kubernetes.io/not-ready",
               "operator":"Exists",
               "effect":"NoExecute",
               "tolerationSeconds":300
            },
            { 
               "key":"node.kubernetes.io/unreachable",
               "operator":"Exists",
               "effect":"NoExecute",
               "tolerationSeconds":300
            }
         ],
         "priority":0,
         "enableServiceLinks":true
      },
      "status":{ 

      }
   },
   "oldObject":null,
   "dryRun":false
}
```

and response:
```
{ 
   "Patches":null,
   "uid":"",
   "allowed":true,
   "status":{ 
      "metadata":{ 

      },
      "reason":"no change - not implemented",
      "code":200
   }
}
```