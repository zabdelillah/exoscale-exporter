# exoscale-exporter
## Exoscale Exporter for Prometheus, written in Go
Based on the prometheus and egoscale v3 APIs, this tool is intended to convert platform-level statistics from the Exoscale provider into any prometheus instance.

## Getting Started
To launch the compiled binary, API Keys from Exoscale are required. The same environment variables used by the Exoscale CLI and Terraform can be used, which are:
- `EXOSCALE_API_KEY`
- `EXOSCALE_API_SECRET`

The `listen` flag can be used to drive what port to listen on. By default, port `:9999` will
be used.

Upon performing a hit to the `/metrics` endpoint, the Exoscale API will be queried and formatted for Prometheus. There is no caching directly within this binary.

## Exported Basic Statistics
- [x] API Keys
- [x] Instances (except data-center information)
- [x] Account Balance
- [x] SKS Clusters
- [x] Block Storage Volumes
- [x] SOS Buckets
- [ ] DBaaS
- [x] DNS
- [x] Security Groups
- [x] Elastic IPs
- [x] Load Balancers
- [x] VPCs
- [x] SSH Keys
- [x] Anti-Affinity Groups

## Repository To-Dos
- [ ] Complete Unit Tests
- [-] Automatic builds
- [-] CLI Flags & Configuration File

## Metrics
|   |   |
|---|---|
|**IAM Keys**||
|exoscale_iam_key{"key", "name", "role"}||
|exoscale_iam_key_count||
|**Instances**||
|exoscale_instance_up{"id", "name", "family", "size", "zone"}||
|exoscale_instance_cpus{"id", "name", "family", "size", "zone"}||
|exoscale_instance_gpus{"id", "name", "family", "size", "zone"}||
|exoscale_instance_memory{"id", "name", "family", "size", "zone"}||
|exoscale_instances_count||
|exoscale_instance_pool_up"id", "name"}||
|exoscale_instance_pool_size"id", "name"}||
|exoscale_instance_pool_count||
|**Organization**||
|exoscale_organization_balance{"organization_id", "organization_name"}||
|exoscale_organization_usage{"organization_id", "organization_name"}||
|**SKS**||
|exoscale_sks_cluster_up{"id", "name", "level", "version"}||
|exoscale_sks_cluster_count||
|exoscale_sks_cluster_size||
|exoscale_sks_nodepool_up{"id", "name", "version"}||
|exoscale_sks_nodepool_size{"id", "name", "version"}||
|exoscale_sks_nodepool_disk_size{"id", "name", "version"}||
|**SOS**||
|exoscale_sos_bucket{"name", "zone"}||
|exoscale_sos_bucket_count||
|**Block Storage**||
|exoscale_volume_size{"id", "name"}||
|exoscale_volume_count||