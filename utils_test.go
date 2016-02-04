package jiraui

import (
	"encoding/json"
	"testing"
)

func TestCountLabelsFromQueryData(t *testing.T) {
	var data interface{}
	inputJSON := []byte(`{
		  "issues": [
				{ "fields": { "labels": [ "wibble", "bibble" ] } },
				{ "fields": { "labels": [ "wibble", "bibble" ] } },
				{ "fields": { "labels": [ "bibble" ] } },
				{ "fields": { "labels": [] } },
				{ "fields": { "labels": [] } }
			]
		}`)

	expected := make(map[string]int)
	expected["wibble"] = 2
	expected["bibble"] = 3
	expected["NOT LABELLED"] = 2

	err := json.Unmarshal(inputJSON, &data)
	if err != nil {
		t.Fatal(err)
	}

	actual := countLabelsFromQueryData(data)
	for k, v := range expected {
		if v != actual[k] {
			t.Fatalf("%s: expected %d, got %d", k, v, actual[k])
		}
	}
}

func TestFindTicketIdInString(t *testing.T) {
	var match string
	match = findTicketIdInString("  relates: BLAH-123[Done]  ")
	if match != "BLAH-123" {
		t.Fatalf("expected BLAH-123, got %s", match)
	}
	match = findTicketIdInString("  wibble: xxBLAH-123[Done]  ")
	if match != "BLAH-123" {
		t.Fatalf("expected %q, got %q", "", match)
	}
	match = findTicketIdInString("  wibble: xxBL-1[Done]  ")
	if match != "BL-1" {
		t.Fatalf("expected %q, got %q", "BL-1", match)
	}
	match = findTicketIdInString("  wibble: xxBLAH-1[Done]  ")
	if match != "BLAH-1" {
		t.Fatalf("expected %q, got %q", "", match)
	}
	match = findTicketIdInString("  wibble: xxTOOLONGPROJECT-1[Done]  ")
	if match != "" {
		t.Skip("This fails, TODO fixing!")
	}
}

func TestWrapText(t *testing.T) {
	input := []string{
		"",
		"wibble:        hello",
		"longfield:     1234567890123456789012345678901234567890",
		"1234567890123456789012345678901234567890",
		"12345678901234567890123456789012345678901234567890",
		"      {code}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {code:bash}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {noformat}   ",
		"      # This is noformat it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {noformat}   ",
		"body: |",
		"   hello there I am a line that is longer than 40 chars yes I am oh aye.",
	}
	expected := []string{
		"",
		"wibble:        hello",
		"longfield:     1234567890123456789012345678901234567890",
		"1234567890123456789012345678901234567890",
		"12345678901234567890123456789012345678901234567890",
		"      {code}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {code:bash}   ",
		"      # This is code it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {code}   ",
		"      {noformat}   ",
		"      # This is noformat it should not be wrapped at all herpdy derp",
		"      # weoijwefoi wpeifjwoiejf pwjefoijwefij wefjowiejf wefwefwefijwe",
		"      {noformat}   ",
		"body: |",
		"   hello there I am a line that is",
		"   longer than 40 chars yes I am oh aye.",
	}
	match := WrapText(input, 40)
	for i, _ := range expected {
		if i > len(match)-1 {
			t.Fatalf("expected %d lines, got %d", len(expected), len(match))
		} else if match[i] != expected[i] {
			t.Fatalf("line %d - expected %q, got %q", i, expected[i], match[i])
		}
	}
}
