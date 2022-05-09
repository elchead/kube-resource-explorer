package monitoring

import "fmt"

type Controller struct {
	Client Clienter
}

func (c Controller) GetMigrations() bool {
	free, _ := c.Client.GetFreeMemoryNode("zone2")
	fmt.Println(free)
	if free < 20 {
		return true
	} else {
		return false
	}
}
