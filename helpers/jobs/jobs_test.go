package jobs

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"testing"
	"time"
)

func testFun() {
	fmt.Println(":hell ya")
}

func testJob(c *gocron.Scheduler) error {
	_, err := c.Every(1).Second().Tag("refreshToken").Do(testFun)
	if err != nil {
		return err
	}
	return nil
}

func TestNew(t *testing.T) {
	c := New().Get()
	c.StartAsync()

	c.Every(1).Second().Do(testFun)
	time.Sleep(60 * time.Second)

}
