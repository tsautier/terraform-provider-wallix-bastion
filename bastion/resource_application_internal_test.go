package bastion

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFillApplicationHandlesNilLocalDomains(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApplication().Schema, map[string]interface{}{})

	fillApplication(d, jsonApplication{
		ApplicationName:  "test-app",
		ConnectionPolicy: skProtoRDP,
		Category:         skStandard,
	})

	localDomains := d.Get("local_domains")
	if localDomains == nil {
		t.Fatalf("expected local_domains to be initialized")
	}
}

func TestPrepareApplicationJSONRejectsParametersForWebApplication(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApplication().Schema, map[string]interface{}{
		"application_name":  "web-app",
		"connection_policy": "WebApp",
		"category":          "web_application",
		"application_url":   "https://example.com",
		"parameters":        "some-value",
	})

	_, err := prepareApplicationJSON(d, true, VersionWallixAPI312)
	if err == nil {
		t.Fatalf("expected an error when parameters is set with category = web_application")
	}
	if want := "parameters cannot be configured when category = web_application"; err.Error() != want {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

func TestPrepareApplicationJSONOmitsParametersForWebApplication(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApplication().Schema, map[string]interface{}{
		"application_name":  "web-app",
		"connection_policy": "WebApp",
		"category":          "web_application",
		"application_url":   "https://example.com",
	})

	jsonData, err := prepareApplicationJSON(d, true, VersionWallixAPI312)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jsonData.Parameters != nil {
		t.Fatalf("expected Parameters to be nil, got %q", *jsonData.Parameters)
	}

	body, err := json.Marshal(jsonData)
	if err != nil {
		t.Fatalf("unexpected error marshaling: %v", err)
	}
	if strings.Contains(string(body), "parameters") {
		t.Fatalf("expected marshaled JSON to omit the parameters key, got: %s", body)
	}
}

func TestPrepareApplicationJSONKeepsParametersForStandard(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceApplication().Schema, map[string]interface{}{
		"application_name":  "std-app",
		"connection_policy": "RDP",
		"category":          "standard",
		"target":            "cluster",
		"parameters":        "some-value",
		"paths": []interface{}{
			map[string]interface{}{
				"target":      "Interactive@device:SSH",
				"program":     "application_path",
				"working_dir": "directory",
			},
		},
	})

	jsonData, err := prepareApplicationJSON(d, true, VersionWallixAPI312)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jsonData.Parameters == nil || *jsonData.Parameters != "some-value" {
		t.Fatalf("expected Parameters to be \"some-value\", got %v", jsonData.Parameters)
	}
}
