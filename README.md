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

# Daemon Types
### Daemon-Types
- NODE = On a Node as systemd-service
- MASTER = On a master as systemd-service
- POD = Runs inside a docker container

# Checks
The following checks are run:
**TYPE**|**CHECK**
-----|-----|-----
| MASTER | Master-API check                 
| MASTER | ETCD health check                
| MASTER | DNS via kubernetes                
| MASTER | DNS via dnsmasq                   
| MASTER | HTTP check via service            
| MASTER | HTTP check via ha-proxy           
| NODE   | Master-API check                  
| NODE   | DNS via kubernetes                
| NODE   | DNS via dnsmasq                   
| NODE   | HTTP check via service           
| NODE   | HTTP check via ha-proxy           
| POD    | Master-API check                  
| POD    | DNS via kubernetes                
| POD    | DNS via Node > dnsmasq            
| POD    | SDN over http via service check   
| POD    | SDN over http via ha-proxy check  



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
**NAME**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----
HUB\_ADDRESS|Address & port of the hub|localhost:2600
DAEMON\_TYPE|Type of the daemon out of [MASTER|NODE
POD\_NAMESPACE|The namespace if the daemon runs inside a pod in OpenShift|ose-mon-a

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

# IMAGE_SPEC = If you want to use our image use "oscp/openshift-monitoring:version"
oc process -f ose-mon-template.yaml -p DAEMON_PUBLIC_ROUTE=xxx -p DS_HUB_ADDRESS=xxx -p IMAGE_SPEC=xxx | oc create -f -
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


