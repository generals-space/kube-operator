package controller

import (
	"generals-space/kube-operator/pkg/controller/podcluster"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, podcluster.Add)
}
