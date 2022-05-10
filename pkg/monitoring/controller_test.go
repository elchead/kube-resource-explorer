package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCriticalNodes(t *testing.T) {
	mockClient := &mockClient{}
	sut := ThresholdPolicy{20.}
	t.Run("do not migrate if 80% free", func(t *testing.T) {
		nodes := NodeFreeMemMap{"z1": 80., "z2": 90.5}
		mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil).Once()
		assert.Equal(t, 0, len(sut.GetCriticalNodes(mockClient)))
	})
	t.Run("migrate if 10% free", func(t *testing.T) {
		nodes := NodeFreeMemMap{"z1": 80., "z2": 10.5}
		mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil).Once()
		assert.Equal(t, 1, len(sut.GetCriticalNodes(mockClient)))
	})
}

func TestMigration(t *testing.T) {
	t.Run("migrate correct pod on critical node", func(t *testing.T) {
		sut := setupControllerWithMocks([]string{"z2"})
		migs, err := sut.GetMigrations()
		assert.NoError(t, err)
		assert.Equal(t, "z2_q", migs[0].Pod)
	})
	t.Run("do not migrate if other node is full", func(t *testing.T) {
		mockClient := &mockClient{}
		mockPolicy := &mockPolicy{}
		sut := NewControllerWithPolicy(mockClient, mockPolicy)
		mockPolicy.On("GetCriticalNodes", mock.Anything).Return([]string{"z1", "z2"}, nil)
		migs, err := sut.GetMigrations()
		assert.Error(t, err)
		assert.Empty(t, migs)
	})
	t.Run("do not migrate if other node is full after migration", func(t *testing.T) {
		mockClient := &mockClient{}
		sut := NewController(mockClient)
		podsZ1 := PodMemMap{"z1_w": 75, "z2_q": 200, "z3_t": 40}
		podsZ2 := PodMemMap{"z2_w": 100, "z2_q": 200}
		nodes := NodeFreeMemMap{"z1": 30., "z2": 15.5}
		mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil)
		mockClient.On("GetPodMemories", "z2").Return(podsZ2, nil)
		mockClient.On("GetPodMemories", "z1").Return(podsZ1, nil)
		migs, err := sut.GetMigrations()
		assert.Error(t, err)
		assert.Empty(t, migs)
	})
}

func setupControllerWithMocks(criticalNodes []string) *Controller {
	mockClient := &mockClient{}
	mockPolicy := &mockPolicy{}
	sut := NewControllerWithPolicy(mockClient, mockPolicy)
	podsZ2 := PodMemMap{"z2_w": 50, "z2_q": 10000000}
	mockPolicy.On("GetCriticalNodes", mock.Anything).Return(criticalNodes, nil)
	nodes := NodeFreeMemMap{"z1": 90., "z2": 90.}
	mockClient.On("GetFreeMemoryOfNodes").Return(nodes, nil)
	mockClient.On("GetPodMemories", "z2").Return(podsZ2, nil).Once()
	return sut
}

func TestGetMaxPod(t *testing.T) {
	assert.Equal(t, "z1_q", GetMaxPod(PodMemMap{"z1_w": 1000, "z1_q": 5000000}))
}

type mockClient struct {
	mock.Mock
}

func (c *mockClient) GetFreeMemoryOfNodes() (NodeFreeMemMap, error) {
	args := c.Called()
	return args.Get(0).(NodeFreeMemMap), args.Error(1)
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
