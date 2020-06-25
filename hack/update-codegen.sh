#!/bin/bash

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/container-object-storage-interface/api/client \
  github.com/container-object-storage-interface/api/apis \
  cosi.sigs.k8s.io:v1alpha1

client-gen --input cosi.sigs.k8s.io/v1alpha1 \
	   --input-base github.com/container-object-storage-interface/api/apis/ \
	   --output-package github.com/container-object-storage-interface/api/client/ \
	   --output-base $GOPATH --clientset-name "clientset"
