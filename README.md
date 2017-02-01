# General idea
On our OpenShift Environment we have lots of productive apps running, so we can't have any downtime. 
So we try hard to test new versions & config in our test environment. As our test env runs way less pods & traffic we created this tool to check all components under pressure, especially during an upgrade-test.

# Check overview

### Master services
- Check if the master api is available on all times
- Check skydns on the master itself
- Check acces to a running app through sdn

### Worker nodes
- Check skydns on the master
- Check dnsmasq on the node
- Check acces to a running app through sdn
- Check acces to a running app via outside <> haproxy

### Pods
- Check dns
- Check access to another running app in same project
- Check access to another running app in joined project
- Check acces to a running app via outside <> haproxy
- Check access to something outside openshift
- Check access to something on the internet

# Components
- UI
- UI Backend: WS-Gateway
- Daemon: Deploy on master/node/various pods

# Installation

### OpenShift
```bash
oc new-project ose-mon-a
oc new-project ose-mon-b
oc new-project ose-mon-c

# Join projects a <> c
oc adm pod-network join-projects --to=ose-mon-a ose-mon-c

# Install daemonset


```

### Master nodes
TODO: Install hub / ui
