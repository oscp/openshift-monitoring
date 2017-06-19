# General idea
We at [@SchweizerischeBundesbahnen](https://github.com/SchweizerischeBundesbahnen) have lots of productive apps running in our OpenShift environment. So we try really hard to avoid any downtime. 
So we test new things (versions/config and so on) in our test environment. As our test environment runs way less pods & traffic we created this tool to check all important OpenShift components under pressure, especially during a change.

Furthermore the daemon now also has a standalone mode. It runs checks based on a http call. So you can monitor all those things from an external monitoring system.

# Screenshot
![Image of the UI](https://github.com/oscp/openshift-monitoring/raw/master/img/screenshot.png)

# Components
- UI: The UI to controll everything
- Hub: The backend of the UI and the daemons
- Daemon: Deploy them as DaemonSet & manually on master & nodes

# Modes & Daemon Types
### Modes
- HUB = Use the hub as control instance. Hub triggers checks on daemons asynchronously
- STANDALONE = Daemon runs on its own and exposes a webserver to run the checks

### Daemon-Types
- NODE = On a Node as systemd-service
- MASTER = On a master as systemd-service
- STORAGE = On glusterfs server as systemd-service
- POD = Runs inside a docker container

# Checks
### Hub mode
| TYPE   | CHECK                            | 
|--------|----------------------------------| 
| MASTER | Master-API check                 | 
| MASTER | ETCD health check                | 
| MASTER | DNS via kubernetes               | 
| MASTER | DNS via dnsmasq                  | 
| MASTER | HTTP check via service           | 
| MASTER | HTTP check via ha-proxy          | 
| NODE   | Master-API check                 | 
| NODE   | DNS via kubernetes               | 
| NODE   | DNS via dnsmasq                  | 
| NODE   | HTTP check via service           | 
| NODE   | HTTP check via ha-proxy          | 
| POD    | Master-API check                 | 
| POD    | DNS via kubernetes               | 
| POD    | DNS via Node > dnsmasq           | 
| POD    | SDN over http via service check  | 
| POD    | SDN over http via ha-proxy check | 

### Standalone mode
| TYPE    | URL           | CHECK                                                   | 
|---------|---------------|---------------------------------------------------------| 
| ALL     | /fast         | Fast endpoint for http-ping                             | 
| ALL     | /slow         | Slow endpoint for slow http-ping                        | 
| NODE    | /checks/minor | Checks if the dockerpool is > 80%                       | 
|         |               | Checks ntpd synchronization status                      | 
|         |               | Checks if http access via service is ok       | 
| NODE    | /checks/major | Checks if the dockerpool is > 90%                       | 
|         |               | Check if dns is ok via kubernetes & dnsmasq             | 
| MASTER  | /checks/minor | Checks ntpd synchronization status                      | 
|         |               | Checks if external system is reachable                  | 
|         |               | Checks if hawcular is healthy                           | 
|         |               | Checks if ha-proxy has a high restart count             | 
|         |               | Checks if all projects have limits & quotas             | 
|         |               | Checks if logging pods are healthy                      |
|         |               | Checks if http access via service is ok       |
| MASTER  | /checks/major | Checks if output of 'oc get nodes' is fine              | 
|         |               | Checks if etcd cluster is healthy                       | 
|         |               | Checks if docker registry is healthy                    | 
|         |               | Checks if all routers are healthy                       | 
|         |               | Checks if local master api is healthy                   | 
|         |               | Check if dns is ok via kubernetes & dnsmasq             |
| STORAGE | /checks/minor | Checks if open-files count is higher than 200'000 files | 
|         |               | Checks every lvs-pool size. Is the value above 80%?     | 
|         |               | Checks every VG has at least 10% free storage           | 
| STORAGE | /checks/major | Checks if output of gstatus is 'healthy'                | 
|         |               | Checks every lvs-pool size. Is the value above 90%?     | 
|         |               | Checks every VG has at least 5% free storage            | 

# Config parameters
## Hub
**NAME**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----
UI\_ADDR|The address & port where the UI should be hosted|10.10.10.1:80
RPC\_ADDR|The address & port where the hub should be hosted|10.10.10.1:2600
MASTER\_API\_URLS|Names or IPs of your masters with the API port|https://master1:8443
DAEMON\_PUBLIC\_URL|Public url of your daemon|http://daemon.yourdefault.route.com
ETCD\_IPS|Names or IPs where to call your etcd hosts|https://localhost:2379
ETCD\_CERT\_PATH|Optional config of alternative etcd certificates path. This is used during certificate renew process of OpenShift to do checks with the old certificates. If this fails the default path will be checked as well|/etc/etcd/old/

## Daemon
#### Hub mode
**NAME**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----
HUB\_ADDRESS|Address & port of the hub|localhost:2600
DAEMON\_TYPE|Type of the daemon out of [MASTER|NODE
POD\_NAMESPACE|The namespace if the daemon runs inside a pod in OpenShift|ose-mon-a

#### Standalone mode
**NAME**|**DAEMON TYPE**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----|-----
WITH\_HUB|ALL|Disable communication with hub|false
DAEMON\_TYPE|ALL|Type of the daemon out of [MASTER|NODE
SERVER\_ADDRESS|ALL|The address & port where the webserver runs|localhost:2600
POD\_NAMESPACE|NODE|The namespace if the daemon runs inside a pod in OpenShift|ose-mon-a
EXTERNAL\_SYSTEM\_URL|MASTER|URL of an external system to call via http to check external connection|www.google.ch
HAWCULAR\_SVC\_IP|MASTER|Ip of the hawcular service|10.10.10.1
ETCD\_IPS|MASTER|Ips of the etcd hosts with protocol & port|https://192.168.125.241:2379,https://192.168.125.244:2379
REGISTRY\_SVC\_IP|MASTER|Ip of the registry service|10.10.10.1
ROUTER\_IPS|MASTER|Ips of the routers services|10.10.10.1,10.10.10.2
PROJECTS\_WITHOUT\_LIMITS|MASTER|Number of system projects that have no limits & quotas|4
IS\_GLUSTER\_SERVER|STORAGE|Boolean value of the node is a gluster server|true/false

# Installation
### OpenShift
```bash
oc new-project ose-mon-a
oc new-project ose-mon-b
oc new-project ose-mon-c

# Join projects a <> c
oc adm pod-network join-projects --to=ose-mon-a ose-mon-c

# Use the template install/ose-mon-template.yaml
# Do this for each project a,b,c
oc project ose-mon-a

# HUB-Mode: IMAGE_SPEC = If you want to use our image use "oscp/openshift-monitoring:version"
oc process -f ose-mon-template.yaml -v DAEMON_PUBLIC_ROUTE=xxx,DS_HUB_ADDRESS=xxx,IMAGE_SPEC=xxx | oc create -f -

# Standalone-Mode:
oc process -f ose-mon-standalone-template.yaml -v DAEMON_PUBLIC_ROUTE=daemon-ose-mon-b.your-route.com IMAGE_SPEC=oscp/openshift-monitoring:xxxx | oc create -f -
```

### Master nodes
```bash
mkdir -p /opt/ose-mon

# Download and unpack from releases or build it yourself (https://github.com/oscp/openshift-monitoring/releases)

chmod +x /opt/ose-mon/hub /opt/ose-mon/daemon

# Add your params to the service definition files
cp /opt/ose-mon/ose-mon-hub.service  /etc/systemd/system/ose-mon-hub.service
cp /opt/ose-mon/ose-mon-daemon.service  /etc/systemd/system/ose-mon-daemon.service

systemctl start ose-mon-hub.service
systemctl enable ose-mon-hub.service

systemctl start ose-mon-daemon.service
systemctl enable ose-mon-daemon.service
```

### Install the UI
```bash
cd /opt/ose-mon
mkdir static

# The UI is included in the download above
```

### Worker / storage nodes
- Do the same as above, just without the hub


