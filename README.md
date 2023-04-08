# container-info

Application to search pods with given label(s) and return resources information for each container in all pods.

## Build binary and docker image

The project uses goreleaser to build the binary for the application as well as docker image

To build the binary (output in `/dist`):
```shell
goreleaser build --snapshot
```

To build the docker image:
```shell
goreleaser release --snapshot
```

Note: Docker image built and uploaded to dockerhub under `umarhussain/container-info`

## Running the application

The application will only run inside the kubernetes cluster since it will use `serviceaccount` token
added to the pod by the cluster, otherwise it will exit with error.

### Helm chart
Helm chart is provided to install the application inside the cluster. The chart
also creates the `ServiceAccount` for the application and a `ClusterRoleBinding` for this
service account to give view permissions in RBAC.

Install the application with helm chart:
```shell
helm install container-info build/helm/container-info/
```

## Consume application

The REST server api of the application is documented with OpenApi (under `api/`).
For giving label selector to the request the format of the value is similar to `kubectl`.
Example of values for the query:
```
app.kubernetes.io/instance in (app1,app2)
key=value
key!=value
key1 in (app1,app2),key2=work,key3!=dev
```

## Running Test Cases

To run the unit tests of the application run the following command:

```shell
go test -v -cover ./...
```