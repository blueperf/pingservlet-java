package controller

import (
	"ping-operator/pkg/controller/pingservlet"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pingservlet.Add)
}
