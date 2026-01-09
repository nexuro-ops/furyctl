# SIGHUP Distribution Compatibility Matrix

## Overview

This document provides a comprehensive compatibility matrix showing supported combinations of Kubernetes versions, module versions, and furyctl across different cluster types.

## Kubernetes Version Support

### Current Support Status

| Kubernetes Version | Release Date | EOL Date | Status | OnPremises | EKS | KFDDistribution |
|-------------------|--------------|----------|--------|------------|-----|-----------------|
| 1.25.x | Aug 2022 | Aug 2023 | Deprecated | ✓ | ✓ | ✓ |
| 1.26.x | May 2023 | May 2024 | Deprecated | ✓ | ✓ | ✓ |
| 1.27.x | Aug 2023 | Aug 2024 | Deprecated | ✓ | ✓ | ✓ |
| 1.28.x | Aug 2023 | Aug 2024 | Deprecated | ✓ | ✓ | ✓ |
| 1.29.x | Dec 2023 | Dec 2024 | Deprecated | ✓ | ✓ | ✓ |
| 1.30.x | Apr 2024 | Apr 2025 | Deprecated | ✓ | ✓ | ✓ |
| 1.31.x | Aug 2024 | Aug 2025 | Supported | ✓ | ✓ | ✓ |
| 1.32.x | Oct 2024 | Oct 2025 | Supported | ✓ | ✓ | ✓ |
| 1.33.x | Nov 2024 | Nov 2025 | Supported | ✓ | ✓ | ✓ |
| 1.34.x | Dec 2024 | Dec 2025 | Supported | ✓ | ✓ | ✓ |
| **1.35.x** | **Jan 2025** | **Jan 2026** | **NEW** | **✓** | **✓** | **✓** |

## Kubernetes 1.35.x Detailed Support

### Version Range

Supported patch versions: **v1.35.0 through v1.35.5**

| Patch Version | Release Date | Status | OnPremises | EKS | KFDDistribution |
|---------------|--------------|--------|------------|-----|-----------------|
| v1.35.0 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |
| v1.35.1 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |
| v1.35.2 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |
| v1.35.3 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |
| v1.35.4 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |
| v1.35.5 | Jan 2025 | ✓ Supported | ✓ | ✓ | ✓ |

## Module Compatibility by Kubernetes Version

### Kubernetes 1.35.x Module Requirements

The following table shows minimum required module versions for each Kubernetes version:

| Module | 1.29 | 1.30 | 1.31 | 1.32 | 1.33 | 1.34 | **1.35** | Notes |
|--------|------|------|------|------|------|------|----------|-------|
| **auth** | v0.5.0+ | v0.5.0+ | v0.5.0+ | v0.5.0+ | v0.5.0+ | v0.5.0+ | **v0.6.0+** | WebSocket RBAC validation |
| **aws** | v5.0.0+ | v5.0.0+ | v5.0.0+ | v5.0.0+ | v5.0.0+ | v5.1.0+ | **v5.1.0+** | AWS API changes |
| **dr** | v3.0.0+ | v3.0.0+ | v3.0.0+ | v3.0.0+ | v3.2.0+ | v3.2.0+ | **v3.2.0+** | Velero compatibility |
| **ingress** | v4.0.0+ | v4.0.0+ | v4.0.0+ | v4.1.0+ | v4.1.0+ | v4.1.0+ | **v4.1.1+** | Nginx Ingress updates |
| **logging** | v5.0.0+ | v5.0.0+ | v5.0.0+ | v5.1.0+ | v5.2.0+ | v5.2.0+ | **v5.2.0+** | New audit log formats |
| **monitoring** | v3.0.0+ | v3.0.0+ | v3.0.0+ | v4.0.0+ | v4.0.0+ | v4.0.0+ | **v4.0.1+** | Prometheus/Grafana updates |
| **networking** | v2.0.0+ | v2.0.0+ | v2.0.0+ | v3.0.0+ | v3.0.0+ | v3.0.0+ | **v3.0.0+** | Calico CNI improvements |
| **opa** | v1.10.0+ | v1.10.0+ | v1.10.0+ | v1.14.0+ | v1.15.0+ | v1.15.0+ | **v1.15.0+** | Gatekeeper enhancements |
| **tracing** | v1.0.0+ | v1.0.0+ | v1.0.0+ | v1.1.0+ | v1.3.0+ | v1.3.0+ | **v1.3.0+** | Jaeger compatibility |

### Breaking Changes by Module

#### auth (v0.5.0 → v0.6.0)
- **Breaking Change:** WebSocket RBAC validation
- **Impact:** Authentication/Authorization for WebSocket connections now enforced
- **Migration:** Update RBAC policies to explicitly allow WebSocket endpoints
- **Action Required:** Review and update network policies

#### logging (v5.1.0 → v5.2.0)
- **Breaking Change:** New audit log format
- **Impact:** Audit logs now include additional fields for K8s 1.35 events
- **Migration:** Update log parsers and analysis tools
- **Action Required:** Update Elasticsearch/ELK configurations

#### monitoring (v4.0.0 → v4.0.1)
- **Breaking Change:** Prometheus metrics renamed
- **Impact:** Some metric names changed to match K8s 1.35 standards
- **Migration:** Update Grafana dashboards and alerting rules
- **Action Required:** Update Prometheus scrape configs and alerts

#### networking (v3.0.0)
- **Breaking Change:** CNI enhancements
- **Impact:** Network policy enforcement improvements
- **Migration:** No manual action typically needed
- **Action Required:** Verify network policies still working as expected

## Tool Version Compatibility

### Kubernetes 1.35.x Tool Versions

| Tool | Version | Mandatory | Notes |
|------|---------|-----------|-------|
| kubectl | 1.35.x | Yes | Must match cluster version |
| kustomize | 5.6.0+ | Yes | For manifest management |
| terraform | 1.4.6+ | Yes (EKS) | For infrastructure as code |
| helm | 3.12.3+ | Yes | For chart deployments |
| helmfile | 0.156.0+ | No | For helm orchestration |
| yq | 4.34.1+ | No | For YAML processing |
| kapp | 0.64.2+ | No | For application management |
| furyagent | 0.4.0+ | Yes | For cluster management |

### containerd Requirement (OnPremises only)

| Requirement | Version | Impact | Action |
|-------------|---------|--------|--------|
| **containerd** | **2.0+** | **MANDATORY** | Upgrade before cluster creation |
| cgroup | **v2** | **MANDATORY** | Requires modern OS (Ubuntu 22.04+, RHEL 9+) |
| cgroup-v1 | N/A | **DEPRECATED** | No longer supported |

## Cluster Type Specific Support

### OnPremises Clusters

| Feature | 1.35 Support | Notes |
|---------|-------------|-------|
| Master/Worker nodes | ✓ Full | HA configurations supported |
| Custom networking | ✓ Full | Calico v3.0.0+ required |
| Local storage | ✓ Full | No changes to storage |
| Load balancing | ✓ Full | No changes required |
| Monitoring | ✓ Full | Monitoring v4.0.1+ required |
| High Availability | ✓ Full | etcd HA fully supported |
| Multi-network | ✓ Full | CNI improvements apply |

**Key Requirements:**
- cgroup v2 filesystem
- containerd 2.0+
- Supported OS versions (see below)

**Supported Operating Systems:**
- Ubuntu 22.04 LTS or newer
- RHEL/CentOS 9 or newer
- Debian 12 or newer
- AlmaLinux 9 or newer

### EKS Clusters

| Feature | 1.35 Support | Notes |
|---------|-------------|-------|
| Auto-scaling | ✓ Full | AWS Autoscaling groups |
| Load balancing | ✓ Full | AWS NLB/ALB fully supported |
| IAM integration | ✓ Full | IRSA fully compatible |
| VPC networking | ✓ Full | All network modes supported |
| Managed nodes | ✓ Full | Both managed and self-managed nodes |
| Fargate | ✓ Full | ECS Fargate available |

**Key Requirements:**
- AWS CLI 2.8.12+
- Proper IAM permissions
- VPC with subnets (public/private)
- No local system requirements (AWS managed)

**Supported Regions:**
- All standard AWS regions
- AWS GovCloud (available)
- AWS China regions (available)

### KFDDistribution (Existing Clusters)

| Feature | 1.35 Support | Notes |
|---------|-------------|-------|
| Module installation | ✓ Full | All modules v1.35 compatible |
| Namespace creation | ✓ Full | Automatic namespace setup |
| RBAC setup | ✓ Full | Full RBAC support |
| Network policies | ✓ Full | Enhanced with v3.0.0 networking |
| Storage classes | ✓ Full | No breaking changes |
| Custom values | ✓ Full | Helm values fully compatible |

**Key Requirements:**
- Existing Kubernetes 1.35 cluster
- kubeconfig with admin access
- Sufficient cluster resources
- Module versions matching requirements

## Upgrade Path Matrix

### Supported Upgrade Paths

```
1.33 → 1.34 → 1.35 (Recommended path)
1.34 → 1.35       (Direct upgrade)
```

### Not Supported

```
1.32 → 1.35       (Must upgrade through 1.33/1.34)
```

### Downgrade Support

| From | To | Status | Notes |
|------|-----|--------|-------|
| 1.35 | 1.34 | ⚠️ Limited | Possible but not recommended |
| 1.35 | 1.33 | ✗ Not Supported | Requires cluster recreation |

**Note:** Always backup etcd before any upgrade/downgrade operations.

## Known Incompatibilities

### Products Not Yet Supporting 1.35

| Product | Issue | Workaround | Timeline |
|---------|-------|-----------|----------|
| Legacy custom CRDs | API deprecations | Update CRD definitions | Release 1.35.3+ |
| Old RBAC configs | WebSocket validation | Update policies | Immediate |
| IPVS kube-proxy | IPVS deprecated | Migrate to nftables | 1.36+ will remove |

## Feature Gates

### Kubernetes 1.35 Feature Gates (GA)

These features are enabled by default in 1.35:

| Feature | Status | Impact |
|---------|--------|--------|
| PodDisruptionBudgetV2 | GA | Pod disruption improvements |
| CPUManagerPolicyOptions | GA | CPU management enhancements |
| MemoryManagerPolicyOptions | GA | Memory management improvements |
| WindowsHostProcessContainers | GA (Windows only) | Windows container support |

## Performance Characteristics

### Expected Performance Changes from 1.34 to 1.35

| Metric | Change | Impact | Notes |
|--------|--------|--------|-------|
| API latency | +2-5% | Low | RBAC validation overhead |
| Memory usage | +3-8% | Low | Additional WebSocket handling |
| CPU usage | +1-3% | Low | Enhanced security checks |
| Disk I/O | No change | None | Same as 1.34 |

**Recommendation:** Allow 10-15% headroom for resource planning.

## Testing Recommendations

### Before Upgrading

- [ ] Test with test/e2e/testdata/k8s-1.35.x-*.yaml configurations
- [ ] Run ./scripts/validate-k8s-1.35.sh
- [ ] Verify all modules meet minimum versions
- [ ] Test custom applications on 1.35
- [ ] Backup etcd and persistent volumes
- [ ] Review module changelog for 1.35 versions

### During Upgrade

- [ ] Monitor cluster health
- [ ] Watch for pod restart issues
- [ ] Check RBAC access logs
- [ ] Verify networking connectivity
- [ ] Monitor application logs

### After Upgrade

- [ ] Run validation script
- [ ] Verify all pods running
- [ ] Test application functionality
- [ ] Check monitoring and logging
- [ ] Update documentation

## Support Timeline

```
1.35.x Support Timeline:
├─ Jan 2025: Release (1.35.0)
├─ Apr 2025: 1.35 becomes standard
├─ Jan 2026: 1.35 EOL (12 months support)
└─ After EOL: Security patches only
```

## Deprecation Warnings

### In Kubernetes 1.35

⚠️ **IPVS kube-proxy mode** - Will be removed in 1.36
- **Migration path:** Switch to nftables
- **Timeline:** Update recommended for 1.35, required for 1.36

⚠️ **Some API versions** - Older API versions deprecated
- **Check:** kubectl api-resources for version status
- **Action:** Use newer API versions in manifests

## Migration Resources

- **Upgrade Guide:** See [K8S-1.35-UPGRADE-GUIDE.md](./K8S-1.35-UPGRADE-GUIDE.md)
- **Breaking Changes:** See [K8S-1.35-BREAKING-CHANGES.md](./K8S-1.35-BREAKING-CHANGES.md)
- **Migration Guide:** See [K8S-1.35-MIGRATION-GUIDE.md](./K8S-1.35-MIGRATION-GUIDE.md)
- **Testing Guide:** See [K8S-1.35-TESTING.md](./K8S-1.35-TESTING.md)

## Support Matrix Versions

- **Last Updated:** January 2025
- **Valid For:** furyctl v0.34.0+
- **Distribution:** fury-distribution v1.35.0+

## Quick Reference

### Minimum Versions for 1.35

```
Kubernetes:    1.35.0+
containerd:    2.0.0+ (OnPremises)
cgroup:        v2 (OnPremises)
auth module:   v0.6.0+
logging:       v5.2.0+
monitoring:    v4.0.1+
networking:    v3.0.0+
kubectl:       1.35.x
```

### Recommended Versions for 1.35

```
Kubernetes:    1.35.5+
containerd:    2.0.2+ or latest stable
auth module:   v0.6.0+
logging:       v5.2.0+
monitoring:    v4.0.1+
networking:    v3.0.0+
kubectl:       1.35.5
furyctl:       v0.34.0+
```

## Related Documentation

- [Kubernetes 1.35 Release Notes](https://kubernetes.io/blog/2025/01/kubernetes-1-35-release/)
- [SIGHUP Distribution Documentation](https://docs.fury.sighup.io)
- [furyctl Configuration Guide](./CONFIG-GUIDE.md)
- [SIGHUP Distribution Changelog](https://github.com/sighupio/fury-distribution/releases)
