package migration_test

import (
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/migration"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	pod := "o10n-worker-s-q697b-lvfw7"
	node := "zone2"
	sut := migration.New(pod, node)
	manifest := sut.GetManifest() // yaml... send to kubectl apply
	assert.Equal(t, "", manifest)
}
