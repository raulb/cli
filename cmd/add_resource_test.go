package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	mock "github.com/meroxa/cli/mock-cmd"
	utils "github.com/meroxa/cli/utils"
	"github.com/meroxa/meroxa-go"
	"github.com/spf13/cobra"
	"reflect"
	"strings"
	"testing"
)

func TestAddResourceArgs(t *testing.T) {
	tests := []struct {
		args []string
		err  error
		name string
	}{
		{nil, nil, ""},
		{[]string{"resName"}, nil, "resName"},
	}

	for _, tt := range tests {
		ar := &AddResource{}
		err := ar.setArgs(tt.args)

		if tt.err != err {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.err, err)
		}

		if tt.name != ar.name {
			t.Fatalf("expected \"%s\" got \"%s\"", tt.name, ar.name)
		}
	}
}

func TestAddResourceFlags(t *testing.T) {
	expectedFlags := []struct {
		name      string
		required  bool
		shorthand string
	}{
		{"type", true, ""},
		{"url", true, "u"},
		{"credentials", false, ""},
		{"metadata", false, "m"},
	}

	c := &cobra.Command{}
	ar := &AddResource{}
	ar.setFlags(c)

	for _, f := range expectedFlags {
		cf := c.Flags().Lookup(f.name)
		if cf == nil {
			t.Fatalf("expected flag \"%s\" to be present", f.name)
		}

		if f.shorthand != cf.Shorthand {
			t.Fatalf("expected shorthand \"%s\" got \"%s\" for flag \"%s\"", f.shorthand, cf.Shorthand, f.name)
		}

		if f.required && !utils.IsFlagRequired(cf) {
			t.Fatalf("expected flag \"%s\" to be required", f.name)
		}
	}
}

func TestAddResourceExecution(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockAddResourceClient(ctrl)

	r := meroxa.CreateResourceInput{
		Type:        "postgres",
		Name:        "",
		URL:         "https://foo.url",
		Credentials: nil,
		Metadata:    nil,
	}

	cr := utils.GenerateResource()

	client.
		EXPECT().
		CreateResource(
			ctx,
			&r,
		).
		Return(&cr, nil)

	output := utils.CaptureOutput(func() {
		ar := &AddResource{}
		got, err := ar.execute(ctx, client, r)

		if !reflect.DeepEqual(got, &cr) {
			t.Fatalf("expected \"%v\", got \"%v\"", &cr, got)
		}

		if err != nil {
			t.Fatalf("not expected error, got \"%s\"", err.Error())
		}
	})

	expected := fmt.Sprintf("Adding %s resource...", r.Type)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestAddResourceOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = false

	output := utils.CaptureOutput(func() {
		ar := &AddResource{}
		ar.output(&r)
	})

	expected := fmt.Sprintf("%s resource with name %s successfully added!", r.Type, r.Name)

	if !strings.Contains(output, expected) {
		t.Fatalf("expected output \"%s\" got \"%s\"", expected, output)
	}
}

func TestAddResourceJSONOutput(t *testing.T) {
	r := utils.GenerateResource()
	flagRootOutputJSON = true

	output := utils.CaptureOutput(func() {
		ar := &AddResource{}
		ar.output(&r)
	})

	var parsedOutput meroxa.Resource
	json.Unmarshal([]byte(output), &parsedOutput)

	if !reflect.DeepEqual(r, parsedOutput) {
		t.Fatalf("not expected output, got \"%s\"", output)
	}
}
