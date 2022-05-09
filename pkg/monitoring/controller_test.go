package monitoring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (c *mockClient) GetFreeMemoryNode(nodeName string) (float64, error) {
	args := c.Called(nodeName)
	fmt.Println(args)
	return args.Get(0).(float64), nil
}

func (c *mockClient) GetPodMemories(nodeName string) (PodMemMap, error) {
	return PodMemMap{}, nil
}

func TestControl(t *testing.T) {
	mockClient := &mockClient{}
	sut := Controller{mockClient}
	t.Run("do not migrate if 80%% free", func(t *testing.T) {
		mockClient.On("GetFreeMemoryNode", mock.Anything).Return(80., nil).Once()
		assert.Equal(t, false, sut.GetMigrations())
	})
	t.Run("migrate if 10%% free", func(t *testing.T) {
		mockClient.On("GetFreeMemoryNode", mock.Anything).Return(10., nil).Once()
		assert.Equal(t, true, sut.GetMigrations())
	})
}
