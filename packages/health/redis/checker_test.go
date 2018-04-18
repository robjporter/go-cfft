package redis

import (
	"errors"
	"testing"
)

type upMockRedis struct {
	version string
}

func (r upMockRedis) GetVersion() (string, error) {
	return r.version, nil
}

func TestNewChecker(t *testing.T) {
	c := NewChecker("tcp", ":6379")

	if c.Redis == nil {
		t.Error("c.Redis == nil, wants !nil")
	}
}

func TestChecker_Check_up(t *testing.T) {
	dummyVersion := "3.0.5"
	c := NewCheckerWithRedis(upMockRedis{version: dummyVersion})

	health := c.Check()

	if health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), false)
	}

	version := health.GetInfo("version")

	if version != dummyVersion {
		t.Errorf("version == %s, wants %s", version, dummyVersion)
	}
}

type downMockRedis struct{}

func (r downMockRedis) GetVersion() (string, error) {
	return "", errors.New("Could not connect")
}

func TestChecker_Check_down(t *testing.T) {
	expectedError := "Could not connect"

	c := NewCheckerWithRedis(downMockRedis{})

	health := c.Check()

	if health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), false)
	}

	message := health.GetInfo("error")

	if message != expectedError {
		t.Errorf("message == %s, wants %s", message, expectedError)
	}
}
