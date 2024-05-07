# Remediate

Remediate is an open source solution which can trigger pre-defined action on kubernetes resources to resolve alerts triggered in alertmanager.

Can be helpfull to auto-remediate errors which can be recurring without any durable fixes or for specials events (auto triggering HPA increase when it's limit has been reach). For the moment remediate is only scoped at namespace level with role and rolebinding. It is not able to target cluster's scoped resources like node. Will be maybe supported in the future.

It relies heavily on alertmanager to work at this moment. There is a plan to change remediate code base to transform it into an API which could triggered auto-remediation by any alerting application by sending an HTTP request when an alert is firing.

### How to use it

There is a Helm chart available in the repository which can be used to deploy remediate on your Kubernetes cluster in remediate/ directory.

Remediate query the API of alertmanager each x minute(s) to fetch all the alerts firing on your cluster and will check if any of those alerts name are supported by it's configuration. 

If it's the case it will triggered the define action wanted with this alertName and send an information message on the chosen Slack channel to inform the team that the auto-remediation process has been triggered (Slack is not mandatory).
