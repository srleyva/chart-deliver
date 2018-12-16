# Chart-Deliver

Aimed to provide a easy snapable component for the deployment of K8s applications into a pipeline. This abstracts the responsibility of dealing with k8s specs and helm charts into the tool. These chart templates can then be administered in a different repo and can be module as one chart template can be used to deploy multiple applications with the `values.yaml` being used to customize deployments.
