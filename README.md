# Shepherd Operator

This project provides kubernetes operators which control backing up and restoring an environment.

It is written in Go using the [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) framework.

_Note: Kubebuilder version 2 is not compatible with the version of OpenShift that Shepherd runs on_

## Usage

### Via Shepherd

Shepherd has integration with this operator. It allows administrators to create new `Backup` and `Restore` objects via the Drupal UI.

### Via kube manifests

Backup/Restores can be created by creating a new object with a manifest similar to the example below. In this example:

* The `site` and `environment` labels allow Shepherd to query these objects for display in the UI.
* `volumes` holds a unique name and `claimName` for each PVC that is going to be backed up.
* `mysql` holds a unique name and `secret` for each database that is going to be backed up. The `secret.keys` property is used by the operator to mount environment variables into the containers running the sql dump/restores in order to connect to the database.

```
apiVersion: extension.shepherd/v1
kind: Backup
metadata:
  name: node-123-backup-xyz
  labels:
    site: 456
    environment: 123
spec:
  volumes:
    shared:
      claimName: node-123-shared
  mysql:
    default:
      secret:
        name: node-123
        keys:
          username: DATABASE_USER
          password: DATABASE_PASSWORD
          database: DATABASE_NAME
          hostname: DATABASE_HOST
          port: DATABASE_PORT
status:
  startTime: '2018-11-21T00:16:23Z'
  completionTime: '2018-11-21T00:16:43Z'
  resticId: abcd969xcz
  phase: New|InProgress|Failed|Completed
```

## Containers

TODO

## Cluster Setup

1. Configure your namespace as required. Used for RBAC rules and defaults to 'myproject' (for local development).
    ```
    export NAMESPACE=shepherd-dev
    make kustomize
    ```

2. Install the CRD.
    ```
    make install
    ```
3. Review the RBAC and manager manifests, ensure you are happy with rbac / jobs etc..
    ```
    less config/rbac/rbac_role.yaml
    ```
4. Install RBAC rules and manager.
    ```
    kubectl apply -f config/rbac
    ```
5. Configure RBAC rules for accounts which should have access to create Backup/Restore objects (i.e. shepherd service account)
    ```
    oc create clusterrole shepherd-backups --verb=get,list,create,update,delete --resource=backups,restores
    oc adm policy add-cluster-role-to-user shepherd-backups --serviceaccount=shepherd
    ```

## Development

The codebase is written in Go and uses the Kubebuilder framework. 

* [Getting Started with Go](https://github.com/alco/gostart)
* [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

The core logic of this operator is contained in:
- [pkg/controller/backup/backup_controller.go](pkg/controller/backup/backup_controller.go) in the function `Reconcile()`.
- [pkg/controller/restore/restore_controller.go](pkg/controller/backup/restore_controller.go) in the function `Reconcile()`.

### Getting Started

To get started developing this operator, ensure you the following prerequisites:

* A minishift VM running locally
* Go >=1.10 installed
* An IDE such as VSCode or Goland is recommended 👍

1. Clone the repo to `$GOPATH/src/github.com/universityofadelaide/shepherd-operator`
1. Run `make install` to set up the CRD on minishift
1. Run `make run` to compile the local workspace and run the operator. Keep it running for the next few steps.
1. Backup an environment via the Shepherd UI.
1. Run `oc get jobs` to check out the jobs that were created.
