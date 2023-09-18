# ngctl

Nebula Operator Command Line Tool

# Features

- deploy Nebula Graph studio to connect to Nebula Graph cluster
- deploy Nebula Graph console to connect to Nebula Graph cluster
- show the version of the local ngctl and Nebula Operator installed in the target cluster.
- list all installed Nebula Graph clusters
- specify the Nebula Graph cluster which the current ngctl command operates on
- get information of selected Nebula Graph cluster
- get the details of Nebula Graph cluster components

# Quick Start

download the latest release from GitHub release page, and install it.

## ngctl studio

deploy Nebula Graph studio

- install nebula studio

```text
install nebula graph studio.

Usage:
  ngctl studio install [flags]

Flags:
  -h, --help             help for install
      --image string     image of the nebula graph studio (default "vesoft/nebula-graph-studio:v3.7.0")
      --nodePort int32   nodePort of the nebula graph studio (default 30180)

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
      --name string         name of the nebula graph studio (default "studio")
      --namespace string    namespace of the nebula graph studio (default "default")
```

example:

```text
>> ngctl studio install
2023/09/10 16:08:26 Resource Deployment studio is created
2023/09/10 16:08:26 Resource Service studio is created
```

- uninstall nebula studio

```text
uninstall nebula graph studio.

Usage:
  ngctl studio uninstall [flags]

Flags:
  -h, --help   help for uninstall

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
      --name string         name of the nebula graph studio (default "studio")
      --namespace string    namespace of the nebula graph studio (default "default")
```

```text
>> ngctl studio uninstall
2023/09/10 16:08:39 Resource Deployment studio is removed
2023/09/10 16:08:39 Resource Service studio is removed
```

### options

| option       | description                                                         |
|--------------|---------------------------------------------------------------------|
| --kubeconfig | specify the path of the kubernetes config file                      |
| --name       | specify the name of nebula graph  studio deployment                 |
| --namespace  | specify the namespace of nebula graph  studio deployment            |
| --image      | specify the container image of nebula graph  studio deployment      |
| --nodePort   | specify the NodePort service port of nebula graph studio deployment |

## ngctl console

deploy Nebula Graph console and connect to Nebula Graph cluster

```text
nebula console client for nebula graph.

Usage:
  ngctl console [flags]

Flags:
      --enable_ssl                    connect to NebulaGraph using SSL encryption and two-way authentication.
  -e, --eval string                   set the nGQL statement in string type.
  -f, --file string                   set the path of the file that stores nGQL statements.
  -h, --help                          help for console
      --image string                  image of the nebula graph console (default "vesoft/nebula-console:v3.5")
  -p, --password string               set the password of the NebulaGraph account.
  -n, --pod_name string               set the name of the console pod.  (default "nebula-console")
      --ssl_cert_path string          specify the path of the SSL public key certificate.
      --ssl_private_key_path string   specify the path of the SSL key.
      --ssl_root_ca_path string       specify the path of the CA root certificate.
  -t, --timeout int32                 set the connection timeout in milliseconds.  (default 120)
  -u, --user string                   set the username of the NebulaGraph account.  (default "root")

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl console -u root  -p nebula
2023/09/08 09:09:48 Resource ConfigMap nebula-console is already created
2023/09/08 09:09:48 console pod is already ready, skip init pod
2023/09/08 09:09:48 console pod is ready

Welcome!

(root@nebula) [(none)]> 

Bye root!
Fri, 08 Sep 2023 09:09:49 UTC
```

### options

| option                 | shortcut | description                                                             |
|------------------------|----------|-------------------------------------------------------------------------|
| --image                |          | specify the container image of nebula graph  studio deployment          |
| --kubeconfig           |          | specify the path of the kubernetes config file                          |
| --name                 |          | specify the name of nebula graph  studio deployment                     |
| --namespace            |          | specify the namespace of nebula graph  studio deployment                |
| --enable_ssl           |          | connect to NebulaGraph using SSL encryption and two-way authentication. |
| --ssl_cert_path        |          | specify the path of the SSL public key certificate.                     |
| --ssl_private_key_path |          | specify the path of the SSL key.                                        |
| --ssl_root_ca_path     |          | specify the path of the CA root certificate.                            |
| --user                 | -u       | set the username of the NebulaGraph account.                            |
| --password             | -p       | set the password of the NebulaGraph account.                            |
| --eval                 | -e       | set the nGQL statement in string type.                                  |
| --file                 | -f       | set the path of the file that stores nGQL statements.                   |
| --timeout              | -t       | set the connection timeout in milliseconds.                             |
| --pod_name             | -n       | set the name of the console pod.                                        |

## ngctl use

change the current context to the specified cluster

```text
Specify a Nebula Graph cluster to use

Usage:
  ngctl use [flags]

Flags:
  -h, --help               help for use
      --namespace string   namespace of the nebula graph cluster (default "default")

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl use nebula
2023/09/10 16:15:59 use nebula graph cluster nebula in namespace default
```

### options

| option       | shortcut | description                                    |
|--------------|----------|------------------------------------------------|
| --namespace  |          | specify the namespace of clusters              |
| --kubeconfig |          | specify the path of the kubernetes config file |

## ngctl get

component can be one of the following:
metad, graphd, storaged, volume

```text
get component of nebula graph cluster.

Usage:
  ngctl get [flags]

Flags:
  -A, --all-namespaces   if set, list the nebula graph clusters across all namespaces
  -h, --help             help for get

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl get storaged
+-------------------+-------+---------+--------+------+----------+--------------------+--------------+
| NAME              | READY | STATUS  | MEMORY | CPU  | RESTARTS | AGE                | NODE         |
+-------------------+-------+---------+--------+------+----------+--------------------+--------------+
| nebula-storaged-0 | true  | Running | 100Mi  | 100m |        3 | 81h2m25.303223762s | 192.168.49.2 |
| nebula-storaged-1 | true  | Running | 100Mi  | 100m |        3 | 81h2m25.303225996s | 192.168.49.2 |
| nebula-storaged-2 | true  | Running | 100Mi  | 100m |        3 | 81h2m25.303227228s | 192.168.49.2 |
+-------------------+-------+---------+--------+------+----------+--------------------+--------------+
```

### options

| option           | shortcut | description                                    |
|------------------|----------|------------------------------------------------|
| --all-namespaces | -A       | get component of all namespaces                |
| --namespace      |          | specify the namespace of clusters              |
| --kubeconfig     |          | specify the path of the kubernetes config file |

## ngctl info

get information of selected Nebula Graph cluster

```text
information of nebula graph clusters.

Usage:
  ngctl info [flags]

Flags:
  -h, --help   help for info

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl info
+-------------------+-------------------------------+
| Name              | nebula                        |
| Namespace         | default                       |
| CreationTimestamp | 2023-09-07 07:14:58 +0000 UTC |
+-------------------+-------------------------------+
2023/09/10 16:18:06 Overview:
+----------+---------+-------+---------+-----+--------+------------+-----------+---------+
|          | PHASE   | READY | DESIRED | CPU | MEMORY | DATAVOLUME | LOGVOLUME | VERSION |
+----------+---------+-------+---------+-----+--------+------------+-----------+---------+
| Metad    | Running |     1 |       1 | 1   | 1Gi    | 5Gi        | 1Gi       | v3.4.0  |
| Storaged | Running |     3 |       3 | 1   | 1Gi    | 10Gi       | 1Gi       | v3.4.0  |
| Graphd   | Running |     1 |       1 | 1   | 1Gi    |            | 1Gi       | v3.4.0  |
+----------+---------+-------+---------+-----+--------+------------+-----------+---------+
2023/09/10 16:18:06 Endpoints:
+-----------+--------+-----------+----------------------------------------------------------+
| COMPONENT | NAME   | TYPE      | ENDPOINT                                                 |
+-----------+--------+-----------+----------------------------------------------------------+
| graphd    | thrift | NodePort  | 192.168.49.2:32046                                       |
| graphd    | http   | NodePort  | 192.168.49.2:32298                                       |
| graphd    | http2  | NodePort  | 192.168.49.2:31008                                       |
| metad     | thrift | ClusterIP | nebula-metad-headless.default.svc.cluster.local:9559     |
| metad     | http   | ClusterIP | nebula-metad-headless.default.svc.cluster.local:19559    |
| metad     | http2  | ClusterIP | nebula-metad-headless.default.svc.cluster.local:19560    |
| storaged  | thrift | ClusterIP | nebula-storaged-headless.default.svc.cluster.local:9779  |
| storaged  | http   | ClusterIP | nebula-storaged-headless.default.svc.cluster.local:19779 |
| storaged  | http2  | ClusterIP | nebula-storaged-headless.default.svc.cluster.local:19780 |
| storaged  | admin  | ClusterIP | nebula-storaged-headless.default.svc.cluster.local:9778  |
| graphd    | thrift | ClusterIP | nebula-graphd-svc.default.svc.cluster.local:9669         |
| graphd    | http   | ClusterIP | nebula-graphd-svc.default.svc.cluster.local:19669        |
| graphd    | http2  | ClusterIP | nebula-graphd-svc.default.svc.cluster.local:19670        |
+-----------+--------+-----------+----------------------------------------------------------+
```

### options

| option       | shortcut | description                                    |
|--------------|----------|------------------------------------------------|
| --namespace  |          | specify the namespace of clusters              |
| --kubeconfig |          | specify the path of the kubernetes config file |

## ngctl version

get the version of ngctl

```text
show the version of ngctl and nebula operator.

Usage:
  ngctl version [flags]

Flags:
  -h, --help                        help for version
      --operator-namespace string   namespace of nebula operator (default "nebula-operator-system")

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl version
2023/09/10 16:23:44 ngctl Version: Nebula Operator Command Line Tool,V-0.0.1 [GitSha: a2842efe28adb6e2665cba05780553e9660ac33c GitRef: main]
2023/09/10 16:23:44 Nebula Operator Version: vesoft/nebula-operator:v1.4.2
```

### options

| option               | shortcut | description                                    |
|----------------------|----------|------------------------------------------------|
| --operator-namespace |          | specify the namespace of nebula operator       |
| --kubeconfig         |          | specify the path of the kubernetes config file |

## ngctl list

list all clusters

```text
list all installed nebula graph clusters.

Usage:
  ngctl list [flags]

Flags:
  -A, --all-namespaces     if set, list the nebula graph clusters across all namespaces
  -h, --help               help for list
      --namespace string   namespace of the nebula graph cluster (default "default")

Global Flags:
      --kubeconfig string   path of the kubernetes config file (default "~/.kube/config")
```

example:

```text
>> ngctl list
+-----------+--------+--------+-------+----------+
| NAMESPACE | NAME   | GRAPHD | METAD | STORAGED |
+-----------+--------+--------+-------+----------+
| default   | nebula | 1/1    | 1/1   | 3/3      |
+-----------+--------+--------+-------+----------+
```

### options

| option           | shortcut | description                       |
|------------------|----------|-----------------------------------|
| --all-namespaces | -A       | get component of all namespaces   |
| --namespace      |          | specify the namespace of clusters |

# License

ngctl is licensed under the Apache License 2.0.