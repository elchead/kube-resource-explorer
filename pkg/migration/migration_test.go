package migration_test

import (
	"testing"

	"github.com/elchead/kube-resource-explorer/pkg/migration"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	pod := "o10n-worker-s-dqknc-9xbl2"
	namespace := "playground"
	sut := migration.New(pod, namespace)
	err := sut.Migrate()
	assert.Error(t, err)
}
