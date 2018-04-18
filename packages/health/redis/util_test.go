package redis

import "testing"

func Test_parseInfo(t *testing.T) {
	expected := map[string]string{
		"redis_version":              "3.0.5",
		"redis_git_sha1":             "00000000",
		"redis_git_dirty":            "0",
		"connected_clients":          "1",
		"client_longest_output_list": "0",
	}

	info := parseInfo("# Server\r\nredis_version:3.0.5\r\nredis_git_sha1:00000000\r\nredis_git_dirty:0\r\n\r\n# Clients\r\nconnected_clients:1\r\nclient_longest_output_list:0\r\n")

	for key, expectedValue := range expected {
		v, ok := info[key]

		if !ok {
			t.Errorf("%s is not present in the info map", key)
			t.FailNow()
		}

		if v != expectedValue {
			t.Errorf("v = %s, but expected v = %s", v, expectedValue)
			t.FailNow()
		}
	}
}
