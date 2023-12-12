package atcom

import (
	"reflect"
	"testing"
)

func TestGetMeaningfulPart(t *testing.T) {

	t.Run("error for 'not OK' response", func(t *testing.T) {
		response := []string{"command: command", "not OK"}
		_, _, err := GetMeaningfulPart(response, "command", "")
		if err.Error() != "no ok response" {
			t.Errorf("Expected error 'no ok response', got: %v", err)
		}
	})

	t.Run("error for empty response", func(t *testing.T) {
		response := []string{}
		_, _, err := GetMeaningfulPart(response, "command", "")
		if err.Error() != "no ok response" {
			t.Errorf("Expected error 'no ok response', got: %v", err)
		}
	})

	t.Run("Valid response without prefix", func(t *testing.T) {
		response := []string{"command: someCommand", "data1", "data2", "OK"}
		data, len, err := GetMeaningfulPart(response, "command", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expectedData := []string{"data1", "data2"}
		expectedLen := 2
		if !reflect.DeepEqual(data, expectedData) || len != expectedLen {
			t.Errorf("Unexpected result. Got data: %v, len: %v", data, len)
		}
	})

	t.Run("Valid response with prefix", func(t *testing.T) {
		response4 := []string{"command: someCommand", "prefixData1", "prefixData2", "OK"}
		data, len, err := GetMeaningfulPart(response4, "command", "prefix")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		expectedData4 := []string{"Data1", "Data2"}
		expectedLen4 := 2
		if !reflect.DeepEqual(data, expectedData4) || len != expectedLen4 {
			t.Errorf("Unexpected result. Got data: %v, len: %v", data, len)
		}
	})

}
