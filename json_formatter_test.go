package logrus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields"].(map[string]interface{})["time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry["timestamp"] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry["time"])
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("msg", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields"].(map[string]interface{})["msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields"].(map[string]interface{})["level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestFieldClashWithRemappedFields(t *testing.T) {
	formatter := &ApexUpJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime:  "@timestamp",
			FieldKeyLevel: "@level",
			FieldKeyMsg:   "@message",
		},
	}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"@timestamp": "@timestamp",
		"@level":     "@level",
		"@message":   "@message",
		"timestamp":  "timestamp",
		"level":      "level",
		"msg":        "msg",
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	for _, field := range []string{"timestamp", "level", "msg"} {
		if entry["fields"].(map[string]interface{})[field] != field {
			t.Errorf("Expected field %v to be untouched; got %v", field, entry[field])
		}

		remappedKey := field
		if remapped, ok := entry[remappedKey]; ok {
			t.Errorf("Expected %s to be empty; got %v", remappedKey, remapped)
		}
	}

	for _, field := range []string{"@timestamp", "@level", "@message"} {
		if entry[field] == field {
			t.Errorf("Expected field %v to be mapped to an Entry value", field)
		}

		remappedKey := field
		if remapped, ok := entry["fields"].(map[string]interface{})[remappedKey]; ok {
			if remapped != field {
				t.Errorf("Expected field %v to be copied to %s; got %v", field, remappedKey, remapped)
			}
		} else {
			t.Errorf("Expected field %v to be copied to %s; was absent", field, remappedKey)
		}
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &ApexUpJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&logrus.Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "message") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &ApexUpJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &ApexUpJSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime: "timeywimey",
		},
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "timeywimey") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &ApexUpJSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, FieldKeyTime) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &ApexUpJSONFormatter{}

	b, err := formatter.Format(logrus.WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, FieldKeyTime) {
		t.Error("Timestamp not present", s)
	}
}
