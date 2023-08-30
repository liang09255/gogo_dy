package ggIP

import (
	"common/ggLog"
	"testing"
)

func TestIP(t *testing.T) {
	ggLog.Infof("IP: %v", GetIP())
}

func TestAllIP(t *testing.T) {
	ggLog.Infof("IP: %v", getAllIP())
}
