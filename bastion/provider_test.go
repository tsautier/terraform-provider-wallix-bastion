package bastion_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/wallix/terraform-provider-wallix-bastion/bastion"
)

var (
	testAccProviders = map[string]*schema.Provider{ //nolint: gochecknoglobals
		"wallix-bastion": testAccProvider,
	}
	testAccProvider = bastion.Provider() //nolint: gochecknoglobals
)

func TestProvider(t *testing.T) {
	if err := bastion.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	_ = bastion.Provider()
}

func testAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("WALLIX_BASTION_HOST") == "" {
		t.Fatal("WALLIX_BASTION_HOST must be set for acceptance tests")
	}
	if os.Getenv("WALLIX_BASTION_TOKEN") == "" {
		t.Fatal("WALLIX_BASTION_TOKEN must be set for acceptance tests")
	}
	if os.Getenv("WALLIX_BASTION_USER") == "" {
		t.Fatal("WALLIX_BASTION_USER must be set for acceptance tests")
	}

	if err := testAccProvider.Configure(t.Context(), terraform.NewResourceConfigRaw(nil)); err != nil {
		t.Fatal(err)
	}
}
