package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (c *mockClient) GetFreeMemoryOfNodes() (NodeMemMap, error) {
	args := c.Called()
	return args.Get(0).(NodeMemMap), args.Error(1)
}

func (c *mockClient) GetFreeMemoryNode(nodeName string) (float64, error) {
	args := c.Called(nodeName)
	return args.Get(0).(float64), args.Error(1)
}

func (c *mockClient) GetPodMemories(nodeName string) (PodMemMap, error) {
	args := c.Called(nodeName)
	return args.Get(0).(PodMemMap), args.Error(1)
}

type mockPolicy struct {
	mock.Mock
}

func (c *mockPolicy) GetCriticalNodes(clt Clienter) []string {
	args := c.Called()
	return args.Get(0).([]string)
}

func TestGetCriticalNodes(t *testing.T) {
	mockClient := &mockClient{}
	sut := ThresholdPolicy{20.}
	t.Run("do not migrate if 80% free", func(t *testing.T) {
		nodes := NodeMemMap{"z1": 80., "z2": 90.5}
		mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil).Once()
		assert.Equal(t, 0, len(sut.GetCriticalNodes(mockClient)))
	})
	t.Run("migrate if 10% free", func(t *testing.T) {
		nodes := NodeMemMap{"z1": 80., "z2": 10.5}
		mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil).Once()
		assert.Equal(t, 1, len(sut.GetCriticalNodes(mockClient)))
	})
}

func TestMigration(t *testing.T) {
	sut := setupControllerWithMocks()
	t.Run("migrate correct pod on critical node", func(t *testing.T) {
		migs, err := sut.GetMigrations()
		assert.NoError(t, err)
		assert.Equal(t, "z2_q", migs[0].Pod)
	})
}

func setupControllerWithMocks() Controller {
	mockClient := &mockClient{}
	mockPolicy := &mockPolicy{}
	sut := Controller{mockClient, mockPolicy}
	podsZ2 := PodMemMap{"z2_w": 50, "z2_q": 10000000}
	mockPolicy.On("GetCriticalNodes", mock.Anything).Return([]string{"z2"}, nil).Once()
	mockClient.On("GetPodMemories", "z2").Return(podsZ2, nil).Once()
	return sut
}

func TestGetMaxPod(t *testing.T) {
	assert.Equal(t, "z1_q", GetMaxPod(PodMemMap{"z1_w": 1000, "z1_q": 5000000}))
}
