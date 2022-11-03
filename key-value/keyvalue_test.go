package main

import "testing"

func TestPut(t *testing.T) {
	tests := []struct {
		description string
		key         string
		value       string
	}{
		{
			description: "Key and value filled no errors",
			key:         "x",
			value:       "abc",
		},
		{
			description: "Empty key no errors",
			key:         "",
			value:       "cba",
		},
	}
	for _, c := range tests {
		t.Run(c.description, func(t *testing.T) {
			if got := Put(c.key, c.value); got != nil {
				t.Errorf("Put(%s, %s) = %s", c.key, c.value, got.Error())
			}
		})
	}
}

func TestGet(t *testing.T) {
	Put("x", "abc")

	tests := []struct {
		description string
		key         string
		value       string
	}{
		{
			description: "Get valid key with value",
			key:         "x",
			value:       "abc",
		},
		{
			description: "Get no existing key",
			key:         "y",
			value:       "",
		},
	}
	for _, c := range tests {
		t.Run(c.description, func(t *testing.T) {
			if got, err := Get(c.key); got != c.value {
				t.Errorf("Get(%s) = %s, want %s - error: %s", c.key, got, c.value, err.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	existingKey := "x"
	Put(existingKey, "abc")

	tests := []struct {
		description string
		key         string
		err         error
	}{
		{
			description: "Delete valid key with value",
			key:         existingKey,
			err:         nil,
		},
		{
			description: "Delete no existing key",
			key:         "y",
			err:         nil,
		},
	}
	for _, c := range tests {
		t.Run(c.description, func(t *testing.T) {
			if err := Delete(c.key); err != c.err {
				t.Errorf("Delete(%s), error: %s", c.key, err)
			}
		})
	}
}
