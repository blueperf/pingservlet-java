// +build !ignore_autogenerated

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"ping-operator/pkg/apis/benchmark/v1alpha1.PingServlet":       schema_pkg_apis_benchmark_v1alpha1_PingServlet(ref),
		"ping-operator/pkg/apis/benchmark/v1alpha1.PingServletSpec":   schema_pkg_apis_benchmark_v1alpha1_PingServletSpec(ref),
		"ping-operator/pkg/apis/benchmark/v1alpha1.PingServletStatus": schema_pkg_apis_benchmark_v1alpha1_PingServletStatus(ref),
	}
}

func schema_pkg_apis_benchmark_v1alpha1_PingServlet(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PingServlet is the Schema for the pingservlets API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("ping-operator/pkg/apis/benchmark/v1alpha1.PingServletSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("ping-operator/pkg/apis/benchmark/v1alpha1.PingServletStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta", "ping-operator/pkg/apis/benchmark/v1alpha1.PingServletSpec", "ping-operator/pkg/apis/benchmark/v1alpha1.PingServletStatus"},
	}
}

func schema_pkg_apis_benchmark_v1alpha1_PingServletSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PingServletSpec defines the desired state of PingServlet",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_benchmark_v1alpha1_PingServletStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "PingServletStatus defines the observed state of PingServlet",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
