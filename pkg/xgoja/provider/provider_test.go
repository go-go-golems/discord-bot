package provider

import (
	"context"
	"testing"

	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerutil"
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

func TestCollectModuleSections(t *testing.T) {
	sectionCapability := fakeSectionCapability{slug: "http"}
	sections, err := providerutil.CollectConfigSections([]providerapi.ModuleDescriptor{{
		PackageID:           "test-http",
		ModuleID:            "express",
		PackageCapabilities: []providerapi.PackageCapability{sectionCapability},
	}}, providerapi.SectionContext{RuntimeProfile: "bot", CommandProviderID: "bots"}, map[string]string{schema.DefaultSlug: "bot command schema"})
	if err != nil {
		t.Fatalf("collect sections: %v", err)
	}
	if len(sections) != 1 || sections[0].GetSlug() != "http" {
		t.Fatalf("sections = %#v", sections)
	}
}

func TestInitSelectedModulesInvokesRuntimeInitializer(t *testing.T) {
	factory, err := engine.NewRuntimeFactoryBuilder().Build()
	if err != nil {
		t.Fatalf("build runtime factory: %v", err)
	}
	rt, err := factory.NewRuntime(engine.WithStartupContext(context.Background()), engine.WithLifetimeContext(context.Background()))
	if err != nil {
		t.Fatalf("runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	initializer := &fakeRuntimeInitializer{}
	vals := values.New()
	if err := initSelectedModules(context.Background(), vals, rt, []providerapi.ModuleDescriptor{{
		PackageID:           "test-http",
		ModuleID:            "express",
		PackageCapabilities: []providerapi.PackageCapability{initializer},
	}}); err != nil {
		t.Fatalf("init selected modules: %v", err)
	}
	if !initializer.called {
		t.Fatal("expected runtime initializer to be called")
	}
	if initializer.values != vals {
		t.Fatal("expected parsed values to be passed through")
	}
	if initializer.handle == nil || initializer.handle.Runtime() != rt.VM {
		t.Fatal("expected runtime handle for created runtime")
	}
}

type fakeSectionCapability struct{ slug string }

func (c fakeSectionCapability) CapabilityID() string { return "fake-section" }

func (c fakeSectionCapability) ConfigSections(providerapi.SectionContext) ([]schema.Section, error) {
	section, err := schema.NewSection(c.slug, "Fake section", schema.WithFields(fields.New("enabled", fields.TypeBool)))
	if err != nil {
		return nil, err
	}
	return []schema.Section{section}, nil
}

type fakeRuntimeInitializer struct {
	called bool
	values *values.Values
	handle providerapi.RuntimeHandle
}

func (i *fakeRuntimeInitializer) CapabilityID() string { return "fake-runtime-initializer" }

func (i *fakeRuntimeInitializer) InitRuntimeFromSections(_ context.Context, vals *values.Values, handle providerapi.RuntimeHandle) error {
	i.called = true
	i.values = vals
	i.handle = handle
	return nil
}
