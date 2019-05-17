package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExec(t *testing.T) {
	// Connect to the cluster and pull pod specific information
	clientset :=connect("")
	assert.Equal(t, execute(clientset), 0)
}