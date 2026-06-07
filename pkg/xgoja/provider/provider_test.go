package provider

import (
	"context"
	"testing"

	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerutil"
)

func TestRegister(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
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
	sections, err := providerutil.CollectGlazedConfigSections([]providerapi.ModuleDescriptor{{
		PackageID:           "test-http",
		ModuleID:            "express",
		PackageCapabilities: []providerapi.PackageCapability{sectionCapability},
	}}, providerapi.SectionRequest{CommandProviderID: "bots"}, nil)
	if err != nil {
		t.Fatalf("collect sections: %v", err)
	}
	if len(sections) != 1 || sections[0].GetSlug() != "http" {
		t.Fatalf("sections = %#v", sections)
	}
}

func TestNewRuntimeUsesSectionAwareFactory(t *testing.T) {
	underlying := &fakeRuntimeFactory{}
	factory := &xgojaBotRuntimeFactory{factory: underlying}
	vals := values.New()

	rt, err := factory.newRuntime(contextWithCommandValues(context.Background(), vals))
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	if underlying.newRuntimeCalled {
		t.Fatal("expected section-aware NewRuntimeFromSections, got plain NewRuntime")
	}
	if !underlying.newRuntimeFromSectionsCalled {
		t.Fatal("expected NewRuntimeFromSections to be called")
	}
	if underlying.values != vals {
		t.Fatal("expected parsed values to be passed through")
	}
}

type fakeSectionCapability struct{ slug string }

func (c fakeSectionCapability) CapabilityID() string { return "fake-section" }

func (c fakeSectionCapability) GlazedConfigSections(providerapi.SectionRequest) ([]schema.Section, error) {
	section, err := schema.NewSection(c.slug, "Fake section", schema.WithFields(fields.New("enabled", fields.TypeBool)))
	if err != nil {
		return nil, err
	}
	return []schema.Section{section}, nil
}

type fakeRuntimeFactory struct {
	newRuntimeCalled             bool
	newRuntimeFromSectionsCalled bool
	values                       *values.Values
}

func (f *fakeRuntimeFactory) NewRuntime(ctx context.Context, opts ...require.Option) (*engine.Runtime, error) {
	f.newRuntimeCalled = true
	return newTestRuntime(ctx, opts...)
}

func (f *fakeRuntimeFactory) NewRuntimeFromSections(ctx context.Context, vals *values.Values, opts ...require.Option) (*engine.Runtime, error) {
	f.newRuntimeFromSectionsCalled = true
	f.values = vals
	return newTestRuntime(ctx, opts...)
}

func newTestRuntime(ctx context.Context, _ ...require.Option) (*engine.Runtime, error) {
	factory, err := engine.NewRuntimeFactoryBuilder().Build()
	if err != nil {
		return nil, err
	}
	return factory.NewRuntime(engine.WithStartupContext(ctx), engine.WithLifetimeContext(ctx))
}
