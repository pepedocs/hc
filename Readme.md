# Hybrid Cloud CLI 
A CLI for locally provisioning a container that provides a management environment for managing an OpenShift-based hybrid cloud.


```$ hc --help
A CLI for locally provisioning a container that provides a management 
	environment for managing an OpenShift-based hybrid cloud.

Usage:
  hc [command]

Available Commands:
  build            Builds the hc image
  clusterLogin     Logs in to an hybrid-cloud OpenShift cluster.
  completion       Generate the autocompletion script for the specified shell
  currentCluster   Shows the current cluster where a user is logged in.
  currentNamespace Shows OpenShift's current context namespace given an OpenShift user.
  help             Help about any command
  login            Runs the hc container and logs into OCM and optionally to a cluster

Flags:
      --config string   config file (default is $HOME/.hc.yaml)
  -d, --debug           verbose logging
  -h, --help            help for hc
  -t, --toggle          Help message for toggle
```