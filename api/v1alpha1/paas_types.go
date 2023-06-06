/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	resourcev1 "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Dit zu eigenlijk in een Configmap voor de operator erbij moeten...
var (
	DEFAULT_QUOTAS = map[string]map[string]string{
		"argocd": {
			"limits.cpu":       "1",
			"limits.memory":    "4Gi",
			"requests.cpu":     "800m",
			"requests.memory":  "3Gi",
			"requests.storage": "10Gi",
			"thin.storageclass.storage.k8s.io/persistentvolumeclaims": "0",
		},
		"ci": {
			"limits.cpu":       "1",
			"limits.memory":    "4Gi",
			"requests.cpu":     "800m",
			"requests.memory":  "3Gi",
			"requests.storage": "10Gi",
			"thin.storageclass.storage.k8s.io/persistentvolumeclaims": "0",
		},
		"sso": {
			"limits.cpu":       "1",
			"limits.memory":    "4Gi",
			"requests.cpu":     "800m",
			"requests.memory":  "3Gi",
			"requests.storage": "10Gi",
			"thin.storageclass.storage.k8s.io/persistentvolumeclaims": "0",
		},
		"grafana": {
			"limits.cpu":       "1",
			"limits.memory":    "4Gi",
			"requests.cpu":     "800m",
			"requests.memory":  "3Gi",
			"requests.storage": "10Gi",
			"thin.storageclass.storage.k8s.io/persistentvolumeclaims": "0",
		},
	}
)

type PaasQuotas map[corev1.ResourceName]resourcev1.Quantity

func (pq PaasQuotas) QuotaWithDefaults(defaults string) (q PaasQuotas) {
	q = make(PaasQuotas)
	for key, value := range DEFAULT_QUOTAS[defaults] {
		q[corev1.ResourceName(key)] = resourcev1.MustParse(value)
	}
	for key, value := range pq {
		q[key] = value
	}
	return pq
}

// PaasSpec defines the desired state of Paas
type PaasSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Cabailities is a subset of capabilities that will be available in this Pass Project
	Capabilities PaasCapabilities `json:"capabilities,omitempty"`

	//Oplosgroep is an informational field which decides on the oplosgroep that is responsible
	Oplosgroep string `json:"oplosGroep"`

	//LdapGroups is a list of strings defining the ldap groups that should be brought synced from ldap and that should get permissions on all other projects
	LdapGroups []string `json:"ldapGroups,omitempty"`

	Groups map[string][]string `json:"groups,omitempty"`

	// Quota defines the quotas which should be set on the cluster resource quota as used by this PaaS project
	Quota PaasQuotas `json:"quota"`
}

// see config/samples/_v1alpha1_paas.yaml for example of CR

type PaasCapabilities struct {
	// ArgoCD defines the ArgoCD deployment that should be available.
	ArgoCD PaasArgoCD `json:"argocd,omitempty"`
	// CI defines the settings for a CI namespace (tekton) for this PAAS
	CI PaasCI `json:"ci,omitempty"`
	// SSO defines the settings for a SSO (KeyCloak) namwespace for this PAAS
	SSO PaasSSO `json:"sso,omitempty"`
	// Grafana defines the settings for a Grafana monitoring namespace for this PAAS
	Grafana PaasGrafana `json:"grafana,omitempty"`
}

type PaasArgoCD struct {
	// Do we want an ArgoCD namespace, default false
	Enabled bool `json:"enabled,omitempty"`
	// The URL that contains the Applications / Application Sets to be used by this ArgoCD
	GitUrl string `json:"gitUrl,omitempty"`
	// This project has it's own ClusterResourceQuota seetings
	Quota PaasQuotas `json:"quota,omitempty"`
}

func (pa PaasArgoCD) QuotaWithDefaults() (pq PaasQuotas) {
	return pa.Quota.QuotaWithDefaults("argocd")
}

type PaasCI struct {
	// Do we want a CI (Tekton) namespace, default false
	Enabled bool `json:"enabled,omitempty"`
	// This project has it's own ClusterResourceQuota seetings
	Quota PaasQuotas `json:"quota,omitempty"`
}

func (pc PaasCI) QuotaWithDefaults() (pq PaasQuotas) {
	return pc.Quota.QuotaWithDefaults("ci")
}

type PaasSSO struct {
	// Do we want an SSO namespace, default false
	Enabled bool `json:"enabled,omitempty"`
	// This project has it's own ClusterResourceQuota seetings
	Quota PaasQuotas `json:"quota,omitempty"`
}

func (ps PaasSSO) QuotaWithDefaults() (pq PaasQuotas) {
	return ps.Quota.QuotaWithDefaults("sso")
}

type PaasGrafana struct {
	// Do we want a Grafana namespace, default false
	Enabled bool `json:"enabled,omitempty"`
	// This project has it's own ClusterResourceQuota seetings
	Quota PaasQuotas `json:"quota,omitempty"`
}

func (pg PaasGrafana) QuotaWithDefaults() (pq PaasQuotas) {
	return pg.Quota.QuotaWithDefaults("grafana")
}

// PaasStatus defines the observed state of Paas
type PaasStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	ArgoCDUrl  string `json:"argocdUrl"`
	GrafanaUrl string `json:"grafanaUrl"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:resource:path=paas,scope=Cluster

// Paas is the Schema for the paas API
type Paas struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PaasSpec   `json:"spec,omitempty"`
	Status PaasStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PaasList contains a list of Paas
type PaasList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Paas `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Paas{}, &PaasList{})
}