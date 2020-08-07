#!/bin/bash

deepcopy-gen --input-dirs github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1 \
	     --output-base $GOPATH/src \
	     --output-file-base zz_generated.deepcopy \
	     --output-package github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1

openapi-gen --input-dirs github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1 \
	    --output-base $GOPATH/src \
	    --output-package github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1

defaulter-gen --input-dirs github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1 \
	      --output-base $GOPATH/src \
	      --output-package github.com/container-object-storage-interface/api/defaulters

lister-gen --input-dirs github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1 \
	      --output-base $GOPATH/src \
	      --output-package github.com/container-object-storage-interface/api/listers

informer-gen --input-dirs github.com/container-object-storage-interface/api/apis/cosi.sigs.k8s.io/v1alpha1 \
	      --output-base $GOPATH/src \
	      --listers-package github.com/container-object-storage-interface/api/listers \
	      --versioned-clientset-package github.com/container-object-storage-interface/api/clientset \
	      --output-package github.com/container-object-storage-interface/api/informers


controller-gen crd:crdVersions=v1 paths=./apis/...

client-gen --input cosi.sigs.k8s.io/v1alpha1 \
	   --input-base github.com/container-object-storage-interface/api/apis/ \
	   --output-package github.com/container-object-storage-interface/api/ \
	   --output-base $GOPATH/src \
	   --clientset-name "clientset"
