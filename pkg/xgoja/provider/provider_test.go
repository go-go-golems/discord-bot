package provider

import (
	"testing"

	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

func TestRegister(t *testing.T) {
	registry := providerapi.NewRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register: %v", err)
	}
	for _, name := range []string{"discord", "ui"} {
		if _, ok := registry.ResolveModule(PackageID, name); !ok {
			t.Fatalf("expected module %s", name)
		}
	}
	if _, ok := registry.ResolveCommandSetProvider(PackageID, "bots"); !ok {
		t.Fatal("expected bots command set provider")
	}
}
