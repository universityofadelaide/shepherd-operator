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

```bash
$ docker pull ghcr.io/universityofadelaide/shepherd-operator:latest
```

## Deploy

This approach can be used for both production deployments and local develop on [Shepherd](https://github.com/universityofadelaide/shepherd)

1. Generate Github Token

Create new Github Personal Access Token with `read:packages` scope.

https://github.com/settings/tokens/new

Update the [secret manifest](config/manager/secret.yml) to contain a base64 string using the following command:

```bash
echo -n <your-github-username>:<TOKEN> | base64
```

2. Create Namespace

```bash
oc apply -f config/manager/namespace.yml
```

4. Apply CustomerResourceDefinitions

```bash
oc apply -f config/crd/bases/

```

5. Apply RBAC Policies

```bash
oc apply -f config/rbac/
```

6. Apply Deployment Manifests

```bash
oc apply -f config/manager/deploy.yml
```

## Local Development

1. Setup Permissions

```bash
CMD=$(crc console --credentials | grep kubeadmin | cut -d"'" -f2); ${CMD} \
&& oc apply -f config/manager/namespace.yml \
&& oc apply -f config/crd/bases/ \
&& oc apply -f config/rbac/ \
&& CMD=$(crc console --credentials | grep developer | cut -d"'" -f2); ${CMD}
```

2. Run the Operator

```
make run
```

## Resources

The codebase is written in Go and uses the Kubebuilder framework.

* [Getting Started with Go](https://github.com/alco/gostart)
* [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

The core logic of this operator is contained in:
- [pkg/controller/backup/backup_controller.go](pkg/controller/backup/backup_controller.go) in the function `Reconcile()`.
- [pkg/controller/restore/restore_controller.go](pkg/controller/backup/restore_controller.go) in the function `Reconcile()`.
