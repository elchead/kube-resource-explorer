package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEnoughSpaceAvailable(t *testing.T) {
	sut := Cluster{NbrNodes: 2, NodeGb: 100, TargetFreePercent: 20.}

	t.Run("fail if other node full before already", func(t *testing.T) {
		nodes := NodeFreeMemMap{"z1": 19., "z2": 15.5}
		assert.Equal(t, "", sut.enoughSpaceAvailableOn("z1", 50., nodes))
	})
	t.Run("fail if other node would be full after migration", func(t *testing.T) {
		nodes := NodeFreeMemMap{"z1": 19., "z2": 60.}
		assert.Equal(t, "", sut.enoughSpaceAvailableOn("z1", 50., nodes))
	})
	t.Run("succeed if enough space", func(t *testing.T) {
		nodes := NodeFreeMemMap{"z1": 19., "z2": 60.}
		assert.Equal(t, "z2", sut.enoughSpaceAvailableOn("z1", 35., nodes))
	})
}
