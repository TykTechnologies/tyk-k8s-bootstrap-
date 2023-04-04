# Tyk-K8S-bootstrap

This is a standalone app meant to help with bootstrap and deletion of the tyk-stack when installed
via tyk-helm-charts.

**Note:** 
<br>This app is needed only for Tyk [Self-managed](https://tyk.io/docs/tyk-on-premises/) deployment!
<br>[Tyk OSS](https://tyk.io/docs/apim/open-source/) doesn't have a special bootstrap and in [Tyk Cloud](https://tyk.io/docs/tyk-cloud/) it is done for you (being a SaaS).

## What it does?

### 1. Tyk post delpoyment bootstrapping
a. Creates a basic organization with the values specified in the env vars
via the tyk-helm charts
<br>
b. Creates a user to access the dashboard (values determined as above)
<br>
c. Bootstraps tyk-portal with a mock page (only if enabled in tyk-helm-charts)
<br>
d. Creates the secret required for the tyk-operator to work (only if enabled in tyk-helm-charts)



### 2. Tyk pre deletion hook
a. Ensures that no failed jobs are still running by deleting them (as they prevented
<br>
b. clean uninstallation of the helm charts)
<br>
c. Also detects and deletes an existing tyk-operator-secret on helm charts uninstallation

Required RBAC roles for the app to work inside the k8s cluster:
- delete
- list


### Useful debug/test tips/commands:

If you want to create a k8s kind cluster that also has a local repository where
you can push the images generated by the Makefile just run the "local_registry.sh" script.
After, the commands below help you with building and pushing the containers to the local repository.


```bash
(rm bin/bootstrapapp-post || true) && make build-bootstrap-post && docker build -t localhost:5001/bootstrap-tyk-post:$bsVers -f ./.container/image/bootstrap-post/Dockerfile ./bin && docker push localhost:5001/bootstrap-tyk-post:$bsVers
```
```bash
(rm bin/bootstrapapp-pre-delete) && make build-bootstrap-pre-delete && docker build -t localhost:5001/bootstrap-tyk-pre-delete:$bsVers -f ./.container/image/bootstrap-pre-delete/Dockerfile ./bin & docker push localhost:5001/bootstrap-tyk-pre-delete:$bsVers
```

The "hack" folder comes with a job (job.yaml) that can be applied directly together
with the role.yaml (which contains the ServiceAccount and ClusterRoleBinding) 
into a namespace running an empty tyk stack for debugging/development purposes.
