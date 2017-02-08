# General idea
We at [@SchweizerischeBundesbahnen](https://github.com/SchweizerischeBundesbahnen) have lots of productive apps running in our OpenShift environment. So we try really hard to avoid any downtime. 
So we test new things (versions/config and so on) in our test environment. As our test environment runs way less pods & traffic we created this tool to check all important OpenShift components under pressure, especially during a change.

# Check overview

### Master services
- Check if the master api is available at all times
- Check skydns on the master itself
- Check access to a running app through sdn
- Check etcd health

### Worker nodes
- Check skydns on the master
- Check dnsmasq on the node
- Check access to a running app through sdn
- Check access to a running app via haproxy

### Pods
- Check dns
- Check access to another running app in same project
- Check access to another running app in joined project
- Check access to a running app via haproxy

# Components
- UI: The UI to controll everything
- Hub: The backend of the UI and the daemons
- Daemon: Deploy them as DaemonSet & manually on master & nodes

# Installation

### Config parameters
#### Hub
**NAME**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----
UI\_ADDR|The address & port where the UI should be hosted|10.10.10.1:80
RPC\_ADDR|The address & port where the hub should be hosted|10.10.10.1:2600
MASTER\_API\_URLS|Names or IPs of your masters with the API port|https://master1:8443
DAEMON\_PUBLIC\_URL|Public url of your daemon|http://daemon.yourdefault.route.com
ETCD\_IPS|Names or IPs where to call your etcd hosts|https://localhost:2379
ETCD\_CERT\_PATH|Optional config of alternative etcd certificates path. This is used during certificate renew process of OpenShift to do checks with the old certificates. If this fails the default path will be checked as well|/etc/etcd/old/

#### Daemon
**NAME**|**DESCRIPTION**|**EXAMPLE**
-----|-----|-----
HUB\_ADDRESS|Address & port of the hub|localhost:2600
DAEMON\_TYPE|Type of the daemon out of [MASTER|NODE
POD\_NAMESPACE|The namespace if the daemon runs inside OpenShift|ose-mon-a

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
oc process -f ose-mon-template.yaml -v DAEMON_PUBLIC_ROUTE=xxx,DS_HUB_ADDRESS=xxx,IMAGE_SPEC=xxx | oc create -f -
```

### Master nodes
```bash
mkdir -p /opt/ose-mon

# Download and unpack from releases or build it yourself (https://github.com/oscp/openshift-monitoring/releases)

chmod +x /opt/ose-mon/hub /opt/ose-mon/daemon

# Add your params to the service definition files
ln -s /opt/ose-mon/ose-mon-hub.service  /etc/systemd/system/ose-mon-hub.service
ln -s /opt/ose-mon/ose-mon-daemon.service  /etc/systemd/system/ose-mon-daemon.service

systemctl start ose-mon-hub.service
systemctl start ose-mon-daemon.service
```

### Install the UI
```bash
cd /opt/ose-mon
mkdir static

# The UI is included in the download above
```

### Worker nodes
- Do the same as above, just without the hub


