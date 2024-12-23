package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LitmusInstallSpec определяет параметры для установки Litmus
type LitmusInstallSpec struct {
	// Добавьте необходимые поля для спецификации установки
}

// LitmusInstallStatus определяет статус установки Litmus
type LitmusInstallStatus struct {
	// Добавьте статус установки
}

// LitmusInstall представляет объект для установки Litmus Chaos Operator
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type LitmusInstall struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LitmusInstallSpec   `json:"spec,omitempty"`
	Status LitmusInstallStatus `json:"status,omitempty"`
}

// LitmusInstallList содержит список объектов LitmusInstall
// +kubebuilder:object:root=true
type LitmusInstallList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LitmusInstall `json:"items"`
}
