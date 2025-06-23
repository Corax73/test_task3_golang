package app_test

import (
	"checklist/customRedis"
	"context"
	"reflect"
	"slices"
	"testing"
)

func TestGetClient(t *testing.T) {
	client := customRedis.GetClient(context.Background(), "env.test")
	if reflect.TypeOf(client) == reflect.TypeOf(&customRedis.RedisClient{}) {
		t.Log("Done checking type")
	} else {
		t.Errorf("Checking type was incorrect")
	}
	if client.EnvPort == ":6379" {
		t.Log("Done checking default port")
	} else {
		t.Errorf("Checking default port was incorrect")
	}
	if client.EnvPassword == "" {
		t.Log("Done checking default password")
	} else {
		t.Errorf("Checking default password was incorrect")
	}
}
func TestGetAllKeys(t *testing.T) {
	client := customRedis.GetClient(context.Background(), "env.test")
	err := client.RedisClient.Set(
		client.Ctx,
		"1",
		"data1",
		0,
	).Err()
	if err != nil {
		t.Errorf("%s", err.Error())
	} else {
		err := client.RedisClient.Set(
			client.Ctx,
			"2",
			"data2",
			0,
		).Err()
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		keys := client.GetAllKeys()
		if len(keys) == 2 {
			t.Log("Done checking quantity keys")
		} else {
			t.Errorf("Checking quantity keys was incorrect")
		}
		if slices.Contains(keys, "1") && slices.Contains(keys, "2") {
			t.Log("Done checking keys")
		} else {
			t.Errorf("Checking keys was incorrect")
		}
	}
}

func TestRemoveModelKeys(t *testing.T) {
	client := customRedis.GetClient(context.Background(), "env.test")
	err := client.RedisClient.Set(
		client.Ctx,
		"1",
		"data1",
		0,
	).Err()
	if err != nil {
		t.Errorf("%s", err.Error())
	} else {
		err := client.RedisClient.Set(
			client.Ctx,
			"2",
			"data2",
			0,
		).Err()
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		keys := client.GetAllKeys()
		if len(keys) == 2 && slices.Contains(keys, "1") && slices.Contains(keys, "2") {
			client.RemoveModelKeys("1")
			keys = client.GetAllKeys()
			if len(keys) == 1 && !slices.Contains(keys, "1") && slices.Contains(keys, "2") {
				t.Log("Done checking remove keys")
			} else {
				t.Errorf("Checking remove keys was incorrect")
			}
		} else {
			t.Errorf("Checking keys was incorrect")
		}
	}
}
