package main

import (
	"reflect"
	"testing"
)

func TestParseModuleName(t *testing.T) {
	cases := map[string]string{
		"module userv4\n\ngo 1.25\n":     "userv4",
		"// c\nmodule   \"tx\"\n":         "tx",
		"go 1.25\nrequire x v1.0.0\n": "",
	}
	for in, want := range cases {
		if got := parseModuleName(in); got != want {
			t.Fatalf("parseModuleName(%q)=%q want %q", in, got, want)
		}
	}
}

func TestNormalizeScanExclude(t *testing.T) {
	got := normalizeScanExclude([]string{
		" internal/modules/catalog/ ",
		"internal/modules/catalog",
		"internal/modules/subscription/shadow",
		"",
	})
	want := []string{
		"internal/modules/catalog",
		"internal/modules/subscription/shadow",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeScanExclude() = %#v, want %#v", got, want)
	}
}

func TestIsScanExcludedMatchesSubtreeOnly(t *testing.T) {
	excluded := []string{"internal/modules/catalog"}
	for _, path := range []string{
		"internal/modules/catalog",
		"internal/modules/catalog/shadow",
		"internal/modules/catalog/migration/file.go",
	} {
		if !isScanExcluded(path, excluded) {
			t.Fatalf("expected %q to be excluded", path)
		}
	}
	for _, path := range []string{
		"internal/domain/plan",
		"internal/modules/catalogue",
		"internal/modules/subscription",
	} {
		if isScanExcluded(path, excluded) {
			t.Fatalf("did not expect %q to be excluded", path)
		}
	}
}

func TestIsSharedRPCClientConstructor(t *testing.T) {
	for _, name := range []string{
		"NewAuthRPCClient",
		"NewUserRPCClient",
		"NewTQDRPCClient",
	} {
		if !isSharedRPCClientConstructor(name) {
			t.Fatalf("expected %q to be treated as shared RPC constructor", name)
		}
	}
	for _, name := range []string{
		"NewAuthClient",
		"NewAdminUserProfileClient",
		"NewHttpClient",
		"NewPoiGrpcHandler",
	} {
		if isSharedRPCClientConstructor(name) {
			t.Fatalf("did not expect %q to be skipped as shared RPC constructor", name)
		}
	}
}

func TestInjectionProvidesCleanup(t *testing.T) {
	if !injectionProvidesCleanup([]InjectionPackage{{
		Import:       "common/injection",
		Constructors: []string{"RPCClientSet", "NewHttpClient"},
	}}) {
		t.Fatal("RPCClientSet should enable cleanup")
	}
	if injectionProvidesCleanup([]InjectionPackage{{
		Import:       "common/injection",
		Constructors: []string{"NewHttpClient"},
	}}) {
		t.Fatal("NewHttpClient alone should not enable cleanup")
	}
}
