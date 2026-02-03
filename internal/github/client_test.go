package github

import (
	"testing"
)

func TestCheckGHInstalled(t *testing.T) {
	c := NewClient()
	err := c.CheckInstalled()
	if err != nil {
		t.Skipf("gh not installed: %v", err)
	}
}
