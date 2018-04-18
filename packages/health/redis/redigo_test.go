package redis

import (
	"errors"
	"testing"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

type mockConnectionProvider struct {
	conn redigo.Conn
}

func (c mockConnectionProvider) getConnection(network, addr string) (redigo.Conn, error) {
	return c.conn, nil
}

type nilMockConnectionProvider struct{}

func (n nilMockConnectionProvider) getConnection(network, addr string) (redigo.Conn, error) {
	return nil, errors.New("Nil connection")
}

func Test_RedigoRedis_GetVersion_nil(t *testing.T) {
	r := NewRedigo("tcp", ":6379")
	r.provider = nilMockConnectionProvider{}

	_, err := r.GetVersion()

	if err == nil {
		t.Error("err == nil, wants !nil")
	}
}

func Test_RedigoRedis_GetVersion_up(t *testing.T) {
	conn := redigomock.NewConn()

	conn.Command("INFO").Expect("# Server\r\nredis_version:3.0.5\r\nredis_git_sha1:00000000\r\nredis_git_dirty:0\r\n\r\n# Clients\r\nconnected_clients:1\r\nclient_longest_output_list:0\r\n")

	r := NewRedigo("tcp", ":6379")
	r.provider = mockConnectionProvider{conn: conn}
	version, _ := r.GetVersion()

	expectedVersion := "3.0.5"

	if version != expectedVersion {
		t.Errorf("version == %s, wants %s", version, expectedVersion)
	}
}

func Test_RedigoRedis_GetVersion_down(t *testing.T) {
	conn := redigomock.NewConn()

	conn.Command("INFO").ExpectError(errors.New("Error"))

	r := NewRedigo("tcp", ":6379")
	r.provider = mockConnectionProvider{conn: conn}
	_, err := r.GetVersion()

	if err == nil {
		t.Error("err == nil, wants !nil")
	}
}

func Test_RedigoRedis_GetVersion_down_version_not_present(t *testing.T) {
	conn := redigomock.NewConn()

	conn.Command("INFO").Expect("# Server\r\nedis_git_sha1:00000000\r\nredis_git_dirty:0\r\n\r\n# Clients\r\nconnected_clients:1\r\nclient_longest_output_list:0\r\n")

	r := NewRedigo("tcp", ":6379")
	r.provider = mockConnectionProvider{conn: conn}
	_, err := r.GetVersion()

	if err == nil {
		t.Error("err == nil, wants !nil")
	}
}
