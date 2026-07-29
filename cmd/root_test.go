package cmd

import "testing"

func TestResolveModelInExplainModeUsesDefaultExplainModel(t *testing.T) {
	got := resolveModel(defaultModel, true, false)
	if got != defaultExplainModel {
		t.Fatalf("expected explain model %q, got %q", defaultExplainModel, got)
	}
}

func TestResolveModelPrefersExplicitFlagInExplainMode(t *testing.T) {
	got := resolveModel("custom-model", true, true)
	if got != "custom-model" {
		t.Fatalf("expected explicit model %q, got %q", "custom-model", got)
	}
}

func TestResolveModelUsesDefaultModelWhenExplainIsOff(t *testing.T) {
	got := resolveModel(defaultModel, false, false)
	if got != defaultModel {
		t.Fatalf("expected default model %q, got %q", defaultModel, got)
	}
}
