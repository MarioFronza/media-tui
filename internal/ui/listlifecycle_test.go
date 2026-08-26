package ui

import (
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
)

func TestListLifecycle_StartFetch_SetsLoading(t *testing.T) {
	lc := newListLifecycle()
	lc, _ = lc.startFetch()

	if !lc.loading {
		t.Fatalf("expected loading to be true, got false")
	}
}

func TestListLifecycle_StartFetch_ReturnsSpinnerTickCmd(t *testing.T) {
	lc := newListLifecycle()
	lc, cmd := lc.startFetch()
	msg := cmd()

	tick, ok := msg.(spinner.TickMsg)
	if !ok {
		t.Fatalf("expected spinner.TickMsg, got %T", msg)
	}
	if tick.ID != lc.spinner.ID() {
		t.Errorf("expected tick ID %d, got %d", lc.spinner.ID(), tick.ID)
	}
}
