package tw_test

// Validates: REQ-011.
// Per: ADR-0004.
// Discipline: C-14.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/septagon-oss/tw"
)

func TestNewIsEmpty(t *testing.T) {
	if !tw.New().IsEmpty() {
		t.Fatal("New() should be empty")
	}
	if tw.New().Compile() != "" {
		t.Fatal("New().Compile() should be empty string")
	}
}

func TestBasicChaining(t *testing.T) {
	got := tw.New().
		Display(tw.DisplayInlineFlex).
		Items(tw.ItemsCenter).
		Gap(tw.S2).
		Rounded(tw.RadiusXL).
		FontWeight(tw.FontSemibold).
		Compile()

	want := "inline-flex items-center gap-2 rounded-xl font-semibold"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestColorRoles(t *testing.T) {
	cases := []struct {
		name string
		list tw.ClassList
		want string
	}{
		{"bg", tw.New().Bg(tw.SurfacePrimary), "bg-surface-primary"},
		{"text", tw.New().TextColor(tw.FgPrimary), "text-fg-primary"},
		{"border-color", tw.New().BorderColor(tw.BorderPrimary), "border-border-primary"},
		{"ring-color", tw.New().RingColor(tw.RingBrand), "ring-ring-brand"},
		{"bg-transparent", tw.New().Bg(tw.ColorTransparent), "bg-transparent"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.list.Compile(); got != c.want {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestSpacingUtilities(t *testing.T) {
	cases := []struct {
		name string
		list tw.ClassList
		want string
	}{
		{"px-3.5", tw.New().PaddingX(tw.S3_5), "px-3.5"},
		{"py-2.5", tw.New().PaddingY(tw.S2_5), "py-2.5"},
		{"p-4", tw.New().Padding(tw.S4), "p-4"},
		{"mx-auto", tw.New().MarginX(tw.SAuto), "mx-auto"},
		{"gap-2", tw.New().Gap(tw.S2), "gap-2"},
		{"w-full", tw.New().Width(tw.SFull), "w-full"},
		{"min-h-11", tw.New().MinHeight(tw.S11), "min-h-11"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.list.Compile(); got != c.want {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestZLayer(t *testing.T) {
	cases := []struct {
		layer tw.ZLayer
		want  string
	}{
		{tw.ZBelow, "-z-[10]"},
		{tw.ZBase, "z-[0]"},
		{tw.ZModal, "z-[1400]"},
		{tw.ZPopover, "z-[1500]"},
		{tw.ZTooltip, "z-[1700]"},
	}
	for _, c := range cases {
		t.Run(c.layer.String(), func(t *testing.T) {
			if got := c.layer.Class(); got != c.want {
				t.Fatalf("Class() got %q, want %q", got, c.want)
			}
			if got := tw.New().ZLayer(c.layer).Compile(); got != c.want {
				t.Fatalf("compile got %q, want %q", got, c.want)
			}
		})
	}
}

func TestStatePrefix(t *testing.T) {
	got := tw.New().
		Bg(tw.SurfacePrimary).
		On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceHover).TranslateY(tw.TranslateNeg05)
		}).
		Compile()
	want := "bg-surface-primary hover:bg-surface-hover hover:-translate-y-0.5"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFocusVisibleRing(t *testing.T) {
	got := tw.New().
		On(tw.StateFocusVisible, func(c tw.ClassList) tw.ClassList {
			return c.Outline(tw.OutlineNone).
				Ring(tw.Ring2).
				RingColor(tw.RingBrand).
				RingOffset(tw.RingOffset2)
		}).
		Compile()
	want := "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring-brand focus-visible:ring-offset-2"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestNestedPrefixStack(t *testing.T) {
	// sm: + hover: should combine as "sm:hover:..."
	got := tw.New().
		Breakpoint(tw.BreakpointSM, func(c tw.ClassList) tw.ClassList {
			return c.On(tw.StateHover, func(c2 tw.ClassList) tw.ClassList {
				return c2.Bg(tw.SurfaceBrand)
			})
		}).
		Compile()
	want := "sm:hover:bg-surface-brand"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestMerge(t *testing.T) {
	base := tw.New().Display(tw.DisplayFlex).Items(tw.ItemsCenter)
	addition := tw.New().Padding(tw.S4).Rounded(tw.RadiusLG)
	got := base.Merge(addition).Compile()
	want := "flex items-center p-4 rounded-lg"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestImmutability(t *testing.T) {
	base := tw.New().Display(tw.DisplayFlex)
	branch1 := base.Items(tw.ItemsCenter)
	branch2 := base.Items(tw.ItemsStart)

	if got := branch1.Compile(); got != "flex items-center" {
		t.Fatalf("branch1: got %q", got)
	}
	if got := branch2.Compile(); got != "flex items-start" {
		t.Fatalf("branch2: got %q", got)
	}
	if got := base.Compile(); got != "flex" {
		t.Fatalf("base leaked: got %q", got)
	}
}

func TestRawEscapeHatch(t *testing.T) {
	got := tw.New().Raw("custom-class").Compile()
	if got != "custom-class" {
		t.Fatalf("got %q, want %q", got, "custom-class")
	}
	// Empty raw is a no-op.
	if tw.New().Raw("   ").Compile() != "" {
		t.Fatal("empty-ish Raw should be a no-op")
	}
}

func TestRawMultipleClassesRespectPrefix(t *testing.T) {
	got := tw.New().On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
		return c.Raw("opacity-80 scale-105")
	}).Compile()
	want := "hover:opacity-80 hover:scale-105"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEmptyValuesAreNoOps(t *testing.T) {
	got := tw.New().
		Bg(""). // zero Color
		PaddingX("").
		FontWeight("").
		Compile()
	if got != "" {
		t.Fatalf("empty-value chain should compile to empty, got %q", got)
	}
}

func TestCompileIsDeterministic(t *testing.T) {
	build := func() string {
		return tw.New().
			Bg(tw.SurfacePrimary).
			PaddingX(tw.S4).
			Rounded(tw.RadiusLG).
			On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
				return c.Bg(tw.SurfaceHover)
			}).
			Compile()
	}
	first := build()
	for i := range 50 {
		if got := build(); got != first {
			t.Fatalf("iteration %d: got %q, want %q", i, got, first)
		}
	}
}

func TestColorCoverageInCompileTables(t *testing.T) {
	// Every Color const that appears in a known role map must be
	// compilable via that role. We verify the Surface + Foreground +
	// Border + Ring roles each handle their canonical slice.
	surfaces := []tw.Color{
		tw.SurfacePrimary, tw.SurfaceSecondary, tw.SurfaceTertiary,
		tw.SurfaceBrand, tw.SurfaceBrandHover, tw.SurfaceBrandSoft,
		tw.SurfaceSuccess, tw.SurfaceSuccessSoft,
		tw.SurfaceWarning, tw.SurfaceWarningSoft,
		tw.SurfaceDanger, tw.SurfaceDangerSoft,
		tw.SurfaceInfo, tw.SurfaceInfoSoft,
		tw.SurfaceDisabled, tw.SurfaceHover, tw.SurfaceActive,
		tw.SurfaceOverlay, tw.SurfaceInverse,
	}
	for _, c := range surfaces {
		out := tw.New().Bg(c).Compile()
		want := fmt.Sprintf("bg-%s", c)
		if out != want {
			t.Fatalf("bg(%q): got %q, want %q", c, out, want)
		}
	}

	fg := []tw.Color{
		tw.FgPrimary, tw.FgSecondary, tw.FgTertiary, tw.FgMuted,
		tw.FgBrand, tw.FgOnBrand,
		tw.FgSuccess, tw.FgWarning, tw.FgDanger, tw.FgInfo,
		tw.FgDisabled, tw.FgOnSurface, tw.FgOnInverse,
		tw.FgLink, tw.FgLinkHover,
	}
	for _, c := range fg {
		out := tw.New().TextColor(c).Compile()
		want := fmt.Sprintf("text-%s", c)
		if out != want {
			t.Fatalf("text(%q): got %q, want %q", c, out, want)
		}
	}

	borders := []tw.Color{
		tw.BorderPrimary, tw.BorderSecondary, tw.BorderBrand,
		tw.BorderSuccess, tw.BorderWarning, tw.BorderDanger, tw.BorderInfo,
	}
	for _, c := range borders {
		out := tw.New().BorderColor(c).Compile()
		want := fmt.Sprintf("border-%s", c)
		if out != want {
			t.Fatalf("border(%q): got %q, want %q", c, out, want)
		}
	}

	rings := []tw.Color{tw.RingBrand, tw.RingFocus, tw.RingDanger}
	for _, c := range rings {
		out := tw.New().RingColor(c).Compile()
		want := fmt.Sprintf("ring-%s", c)
		if out != want {
			t.Fatalf("ring(%q): got %q, want %q", c, out, want)
		}
	}
}

func TestAllSpacingsHaveRoles(t *testing.T) {
	for _, s := range tw.AllSpacings() {
		// Numeric / fractional steps support all roles.
		// SFull/SAuto may not make sense for some roles but must not panic.
		_ = tw.New().PaddingX(s).Compile()
		_ = tw.New().Margin(s).Compile()
		_ = tw.New().Width(s).Compile()
		_ = tw.New().Gap(s).Compile()
	}
}

func TestNoRedundantSpaces(t *testing.T) {
	got := tw.New().
		Display(tw.DisplayFlex).
		Items("").
		Padding("").
		FontWeight(tw.FontSemibold).
		Compile()
	if strings.Contains(got, "  ") {
		t.Fatalf("output contains double space: %q", got)
	}
}
