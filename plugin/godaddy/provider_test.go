package godaddy

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"godaddy": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

//nolint:unused // used by acceptance tests
func testAccPreCheck(t *testing.T) {
	t.Helper()
	for _, key := range []string{"GODADDY_API_KEY", "GODADDY_API_SECRET", "GODADDY_DOMAIN"} {
		if os.Getenv(key) == "" {
			t.Fatalf("%s must be set for acceptance tests.", key)
		}
	}
}
