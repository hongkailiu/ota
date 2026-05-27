| # | Key | Summary | Status | Resolution | Component | Assignee | Notes |
|---|-----|---------|--------|------------|-----------|----------|-------|
| 0 | [OCPBUGS-20056](https://redhat.atlassian.net/browse/OCPBUGS-20056) | Single short-lived operand blip shouldn't cause authentication operator Available=False | POST |  | apiserver-auth | Ondra Kupka | Possibly dup of  OCPBUGS-66027 |
| 1 | ~~[OCPBUGS-22382](https://redhat.atlassian.net/browse/OCPBUGS-22382)~~ | Image registry experiencing disruption during vSphere serial jobs | Closed | Won't Do | Image Registry | Flavian Missi | Won't Do confirmed |
| 2 | ~~[OCPBUGS-23744](https://redhat.atlassian.net/browse/OCPBUGS-23744)~~ | operator-lifecycle-manager-packageserver ClusterOperator should not blip Available=False on 4.14 to 4.15 updates | Closed | Done-Errata | OLM | Kevin Rizza | Won't Do confirmed: OLMv0 in maintenance mode |
| 3 | [OCPBUGS-23746](https://redhat.atlassian.net/browse/OCPBUGS-23746) | openshift-apiserver ClusterOperator should not blip Available=False on brief missing HTTP content-type | POST |  | openshift-apiserver | Ondra Kupka |  |
| 4 | [OCPBUGS-25739](https://redhat.atlassian.net/browse/OCPBUGS-25739) |  ingress operator appears to be reporting unavailable in error during upgrade in bare metal cluster | ASSIGNED |  | Networking / router | Davide Salerno |  |
| 5 | [OCPBUGS-38661](https://redhat.atlassian.net/browse/OCPBUGS-38661) | clusteroperator/kube-apiserver blips Degraded=True during upgrade test | POST |  | kube-apiserver |  |  |
| 6 | [OCPBUGS-38662](https://redhat.atlassian.net/browse/OCPBUGS-38662) | clusteroperator/kube-controller-manager blips Degraded=True during upgrade test | POST |  | kube-controller-manager | Ondra Kupka |  |
| 7 | [OCPBUGS-38663](https://redhat.atlassian.net/browse/OCPBUGS-38663) | clusteroperator/kube-scheduler blips Degraded=True during upgrade test | New |  | kube-scheduler | Ondra Kupka |  |
| 8 | [OCPBUGS-38676](https://redhat.atlassian.net/browse/OCPBUGS-38676) | clusteroperator/console blips Degraded=True during CI job run | New |  | Management Console | Steve Goodwin |  |
| 9 | [OCPBUGS-42837](https://redhat.atlassian.net/browse/OCPBUGS-42837) | clusteroperator/cloud-controller-manager blips Degraded=True during upgrade test | POST |  | Cloud Compute / Cloud Controller Manager | Nolan Brubaker |  |
| 10 | [OCPBUGS-45921](https://redhat.atlassian.net/browse/OCPBUGS-45921) | clusteroperator/ingress blips Degraded=True during hypershift conformance test | ASSIGNED |  | Networking / router | Andrey Lebedev |  |
| 11 | [OCPBUGS-62627](https://redhat.atlassian.net/browse/OCPBUGS-62627) | cluster operator ingress reported Progressing=True with reason=Reconciling for a node reboot | ASSIGNED |  | Networking / router | Davide Salerno |  |
| 12 | [OCPBUGS-62629](https://redhat.atlassian.net/browse/OCPBUGS-62629) | cluster operator kube-storage-version-migrator reported Progressing=True with reason=KubeStorageVersionMigrator_Deploying for a node reboot | New |  | kube-storage-version-migrator |  |  |
| 13 | [OCPBUGS-62633](https://redhat.atlassian.net/browse/OCPBUGS-62633) | cluster operator service-ca reported Progressing=True with reason=_ManagedDeploymentsAvailable for a node reboot | POST |  | service-ca | Ondra Kupka |  |
| 14 | [OCPBUGS-63116](https://redhat.atlassian.net/browse/OCPBUGS-63116) | cluster operator openshift-controller-manager reported Progressing=True with reason= _DesiredStateNotYetAchieved or RouteControllerManager_DesiredStateNotYetAchieved for a node reboot | New |  | openshift-controller-manager / controller-manager | Prabhakar Palepu |  |
| 15 | [OCPBUGS-64688](https://redhat.atlassian.net/browse/OCPBUGS-64688) | cluster operator console reported Progressing=True with reason= SyncLoopRefresh_InProgress for a node reboot | ASSIGNED |  | Management Console | Steve Goodwin |  |
| 16 | [OCPBUGS-64852](https://redhat.atlassian.net/browse/OCPBUGS-64852) | Cluster Operator cloud-controller-manager did not report Progressing=True during a cluster update | ASSIGNED |  | Cloud Compute / Cloud Controller Manager | Nolan Brubaker |  |
| 17 | ~~[OCPBUGS-65583](https://redhat.atlassian.net/browse/OCPBUGS-65583)~~ | Cluster Operator operator-lifecycle-manager did not report Progressing=True during a cluster update | Closed | Won't Do | OLM | Jordan Keister | Won't Do confirmed: OLMv0 in maintenance mode |
| 18 | [OCPBUGS-65647](https://redhat.atlassian.net/browse/OCPBUGS-65647) | Cluster Operator openshift-samples did not report Progressing=True during a cluster update | New |  | Samples Operator | Shannon Poole |  |
| 19 | [OCPBUGS-65896](https://redhat.atlassian.net/browse/OCPBUGS-65896) | cluster operator authentication reported Progressing=True on cluster scaling up | POST |  | apiserver-auth | Ondra Kupka |  |
| 20 | [OCPBUGS-65984](https://redhat.atlassian.net/browse/OCPBUGS-65984) | kube-storage-version-migrator goes Available=False with reason=KubeStorageVersionMigrator_Deploying during updates | Verified |  | kube-storage-version-migrator | Luis Sanchez | Under evaluation |
| 21 | [OCPBUGS-66101](https://redhat.atlassian.net/browse/OCPBUGS-66101) | Cluster Operator baremetal did not report Progressing=True during a cluster update | New |  | Bare Metal Hardware Provisioning | Iury Gregory Melo Ferreira |  |
| 22 | [OCPBUGS-66213](https://redhat.atlassian.net/browse/OCPBUGS-66213) | image-registry operator changed condition/Available to false during non-upgrade job | New |  | Image Registry | Flavian Missi |  |
| 23 | [OCPBUGS-66225](https://redhat.atlassian.net/browse/OCPBUGS-66225) | clusteroperator/image-registry blips Degraded=True during upgrade test | ASSIGNED |  | Image Registry | Thomas Jungblut |  |
| 24 | [OCPBUGS-67134](https://redhat.atlassian.net/browse/OCPBUGS-67134) | console ClusterOperator should not blip Available=False with the reason=Deployment_InsufficientReplicas | New |  | Management Console | Jakub Hadvig |  |
| 25 | [OCPBUGS-82160](https://redhat.atlassian.net/browse/OCPBUGS-82160) | Test failure in upgrade jobs- [bz-Image Registry] clusteroperator/image-registry should not change condition/Available | ASSIGNED |  | Image Registry | Flavian Missi |  |
| 26 | [OCPBUGS-85677](https://redhat.atlassian.net/browse/OCPBUGS-85677) | cluster operator network reported Progressing=True for a node reboot | New |  | Networking / cluster-network-operator | Jamo Luhrsen |  |
| 27 | [OCPBUGS-86009](https://redhat.atlassian.net/browse/OCPBUGS-86009) | cluster operator dns reported Progressing=True on cluster scaling up | New |  | Networking / DNS | Brett Tofel |  |
| 28 | [OCPBUGS-86308](https://redhat.atlassian.net/browse/OCPBUGS-86308) | cluster operator image-registry reported Progressing=True on cluster scaling up | New |  | Image Registry | Flavian Missi |  |

## Build info

* Build: v20260527-b99a0a5

* `git_commit_origin`: v4.1.0-11397-g8e592a009c
