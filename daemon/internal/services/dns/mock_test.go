package dns

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderValidateAndPreview(t *testing.T) {
	service := NewMockService(filepath.Join(t.TempDir(), "generated", "dns"))
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
	if !strings.Contains(preview, "listen-address=127.0.0.1") {
		t.Fatalf("expected listen-address in preview, got: %s", preview)
	}
}
