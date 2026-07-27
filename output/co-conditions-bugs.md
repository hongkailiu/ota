## OCPBugs on [jira/dashboards/22315](https://redhat.atlassian.net/jira/dashboards/22315) as exceptions in CI
Total: 13 issues | Pixaa: 5 (not closed: 5) | New: 1 | ASSIGNED: 2 | POST: 2 | Verified: 7 | Closed: 1

PIXAA: open/total: 5/5

PIXAA components: Cloud Compute, Cluster Autoscaler, Management Console, OLM

Descoped (not in table): 12 issues - [OCPBUGS-20056](https://redhat.atlassian.net/browse/OCPBUGS-20056), [OCPBUGS-23746](https://redhat.atlassian.net/browse/OCPBUGS-23746), [OCPBUGS-38661](https://redhat.atlassian.net/browse/OCPBUGS-38661), [OCPBUGS-38662](https://redhat.atlassian.net/browse/OCPBUGS-38662), [OCPBUGS-38663](https://redhat.atlassian.net/browse/OCPBUGS-38663), [OCPBUGS-62629](https://redhat.atlassian.net/browse/OCPBUGS-62629), [OCPBUGS-62633](https://redhat.atlassian.net/browse/OCPBUGS-62633), [OCPBUGS-63116](https://redhat.atlassian.net/browse/OCPBUGS-63116), [OCPBUGS-65896](https://redhat.atlassian.net/browse/OCPBUGS-65896), [OCPBUGS-65984](https://redhat.atlassian.net/browse/OCPBUGS-65984), [OCPBUGS-66213](https://redhat.atlassian.net/browse/OCPBUGS-66213), [OCPBUGS-86308](https://redhat.atlassian.net/browse/OCPBUGS-86308)

Descoped (others): 2 issues - [OCPBUGS-66027](https://redhat.atlassian.net/browse/OCPBUGS-66027), [OCPBUGS-38678](https://redhat.atlassian.net/browse/OCPBUGS-38678)

Won't Do: 4 issues - [OCPBUGS-22382](https://redhat.atlassian.net/browse/OCPBUGS-22382), [OCPBUGS-23744](https://redhat.atlassian.net/browse/OCPBUGS-23744), [OCPBUGS-65583](https://redhat.atlassian.net/browse/OCPBUGS-65583), [OCPBUGS-65647](https://redhat.atlassian.net/browse/OCPBUGS-65647)

| # | Key | Summary | Status | Resolution | Target Version | Release Blocker | Component | Pixaa | Assignee | Parent | Notes |
|---|-----|---------|--------|------------|----------------|-----------------|-----------|-------|----------|--------|-------|
| 0 | [OCPBUGS-25739](https://redhat.atlassian.net/browse/OCPBUGS-25739) |  ingress operator appears to be reporting unavailable in error during upgrade in bare metal cluster | Closed | Done | 5.0.0 | Rejected | Networking / router |  | Davide Salerno |  | origin#31345, dup of OCPBUGS-92835 |
| 1 | [OCPBUGS-38676](https://redhat.atlassian.net/browse/OCPBUGS-38676) | clusteroperator/console blips Degraded=True during CI job run | Verified |  | 5.0.0 | Rejected | Management Console | ✓ | Steve Goodwin | CONSOLE-5185 | origin#31280 |
| 2 | [OCPBUGS-42837](https://redhat.atlassian.net/browse/OCPBUGS-42837) | clusteroperator/cloud-controller-manager blips Degraded=True during upgrade test | Verified |  | 5.0.0 | Rejected | Cloud Compute / Cloud Controller Manager | ✓ | Nolan Brubaker | OCPCLOUD-3420 | To be removed in 5.1 |
| 3 | [OCPBUGS-62627](https://redhat.atlassian.net/browse/OCPBUGS-62627) | cluster operator ingress reported Progressing=True with reason=Reconciling for a node reboot | New |  | 4.22.0 | Rejected | Networking / router |  | Davide Salerno | NE-2267 |  |
| 4 | [OCPBUGS-64688](https://redhat.atlassian.net/browse/OCPBUGS-64688) | cluster operator console reported Progressing=True with reason= SyncLoopRefresh_InProgress for a node reboot | Verified |  | 5.0.0 |  | Management Console | ✓ | Steve Goodwin | CONSOLE-5185 | Under evaluation, origin#31280, dup of OCPBUGS-93982 |
| 5 | [OCPBUGS-64852](https://redhat.atlassian.net/browse/OCPBUGS-64852) | Cluster Operator cloud-controller-manager did not report Progressing=True during a cluster update | Verified |  | 5.0.0 |  | Cloud Compute / Cloud Controller Manager | ✓ | Theo Barber-Bany | OCPCLOUD-3420 | Under evaluation |
| 6 | [OCPBUGS-66101](https://redhat.atlassian.net/browse/OCPBUGS-66101) | Cluster Operator baremetal did not report Progressing=True during a cluster update | Verified |  | 5.0.0 | Approved | Bare Metal Hardware Provisioning |  | Honza Pokorny |  | Under evaluation |
| 7 | [OCPBUGS-66225](https://redhat.atlassian.net/browse/OCPBUGS-66225) | clusteroperator/image-registry blips Degraded=True during upgrade test | Verified |  | 5.0.0 | Approved | Image Registry |  | Ilias Rinis |  | Under evaluation |
| 8 | [OCPBUGS-67134](https://redhat.atlassian.net/browse/OCPBUGS-67134) | console ClusterOperator should not blip Available=False with the reason=Deployment_InsufficientReplicas | Verified |  | 5.0.0 |  | Management Console | ✓ | Steve Goodwin | CONSOLE-5185 | To be removed in 5.1 |
| 9 | [OCPBUGS-82160](https://redhat.atlassian.net/browse/OCPBUGS-82160) | Test failure in upgrade jobs- [bz-Image Registry] clusteroperator/image-registry should not change condition/Available (multi-p-p) | ASSIGNED |  | 5.0.0 | Rejected | Image Registry |  | Flavian Missi |  |  |
| 10 | [OCPBUGS-85677](https://redhat.atlassian.net/browse/OCPBUGS-85677) | cluster operator network reported Progressing=True for a node reboot | POST |  | 5.0.0 | Rejected | Networking / cluster-network-operator |  | Jamo Luhrsen |  | To be removed in 5.1 |
| 11 | [OCPBUGS-86009](https://redhat.atlassian.net/browse/OCPBUGS-86009) | cluster operator dns reported Progressing=True on cluster scaling up | POST |  | 5.0.0 | Rejected | Networking / DNS |  | Brett Tofel |  |  |
| 12 | [OCPBUGS-86017](https://redhat.atlassian.net/browse/OCPBUGS-86017) | Test failure in upgrade jobs- [bz-Image Registry] clusteroperator/image-registry should not change condition/Available (multi-z-z) | ASSIGNED |  | 5.0.0 | Rejected | Image Registry |  | Rohit Patil |  |  |

## Build info

* Build: v20260727-d240331

* `git_commit_origin`: v4.1.0-11730-g0c1e165773
