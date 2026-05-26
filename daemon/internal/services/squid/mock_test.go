package squid

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderValidateAndPreview(t *testing.T) {
	service := NewMockService(filepath.Join(t.TempDir(), "generated", "squid"))
	if err := service.RenderConfig(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := service.ValidateConfig(context.Background()); err != nil {
		t.Fatal(err)
	}

	preview, err := service.ConfigPreview(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(preview, "cache_dir ufs") {
		t.Fatalf("expected cache_dir line in preview, got: %s", preview)
	}
}
