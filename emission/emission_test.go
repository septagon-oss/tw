// Validates: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

package emission

// emission_test.go pins the three promises this package makes: every color
// role tw declares is mapped, every enumerable class resolves to at least one
// declaration, and the emitted sheet is valid CSS that styleengine can parse
// back. The golden file makes any change to the emitted CSS reviewable as a
// diff instead of a rendering surprise.

import (
	"os"
	"strings"
	"testing"

	"github.com/septagon-oss/pk-design/pkg/themes"
	"github.com/septagon-oss/pk-design/pkg/tokens"
	"github.com/septagon-oss/styleengine"
	"github.com/septagon-oss/tw"
)

func TestRoleMapCoversEveryColor(t *testing.T) {
	t.Parallel()
	roles := roleValues()
	for _, c := range tw.AllColors() {
		if _, ok := roles[c]; !ok {
			t.Errorf("tw.Color %q has no role mapping; add it to roleValues()", c)
		}
	}
	// ColorWhite and ColorBlack are compilable but deliberately outside
	// AllColors (tw documents them as hard-contrast escape values); they are
	// the only entries allowed beyond the enumerator.
	if want := len(tw.AllColors()) + 2; len(roles) != want {
		t.Errorf("role map has %d entries, want %d (AllColors + white + black)", len(roles), want)
	}
}

// TestRoleVarsReferenceRealThemeTokens is the naming contract with pk-design:
// every --pk-color-* variable a role references must exist in the canonical
// theme's rendered CSS. pk-design is a test-only dependency for exactly this.
func TestRoleVarsReferenceRealThemeTokens(t *testing.T) {
	t.Parallel()
	css, err := tokens.CSSVars(themes.Default().Tokens)
	if err != nil {
		t.Fatalf("render theme: %v", err)
	}
	for role, value := range roleValues() {
		for rest := value; ; {
			i := strings.Index(rest, "var(--pk-color-")
			if i < 0 {
				break
			}
			rest = rest[i+len("var("):]
			end := strings.IndexAny(rest, "),")
			name := rest[:end]
			if !strings.Contains(css, name+":") {
				t.Errorf("role %q references %s, which the canonical theme does not define", role, name)
			}
		}
	}
}

func TestBaseCoversEveryEnumerableClass(t *testing.T) {
	t.Parallel()
	classes := baseClasses()
	if len(classes) < 1000 {
		t.Fatalf("enumerated only %d classes; the enumeration itself regressed", len(classes))
	}
	for _, class := range classes {
		if _, err := resolveBase(class); err != nil {
			t.Errorf("enumerable class %q does not resolve: %v", class, err)
		}
	}
}

func TestEscapeHatchesFailClosed(t *testing.T) {
	t.Parallel()
	for _, class := range []string{
		"pk-transition-standard",      // PlatformKitClass handle
		"w-[37px]",                    // arbitrary value
		"bg-[#123456]",                // arbitrary color
		"rotate-45",                   // parametric without table entry
		"completely-made-up-nonsense", // typo
	} {
		if _, err := Rules(class); err == nil {
			t.Errorf("Rules(%q) succeeded; escape hatches must fail closed", class)
		}
	}
}

func TestPrefixedRules(t *testing.T) {
	t.Parallel()
	sheet, err := Rules(
		"hover:bg-surface-brand-hover",
		"md:flex",
		"lg:hover:bg-surface-hover",
		"focus-visible:ring-2",
		"group-hover:underline",
		"dark:bg-surface-inverse",
		"first:pt-0",
		"placeholder:text-fg-placeholder",
	)
	if err != nil {
		t.Fatalf("Rules: %v", err)
	}
	css, err := sheet.Render(styleengine.RenderOptions{Pretty: true})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	for _, want := range []string{
		`.hover\:bg-surface-brand-hover:hover`,
		"@media (min-width: 768px)",
		`.md\:flex`,
		`.lg\:hover\:bg-surface-hover:hover`,
		"@media (min-width: 1024px)",
		`.focus-visible\:ring-2:focus-visible`,
		`.group:hover .group-hover\:underline`,
		`.dark .dark\:bg-surface-inverse`,
		`.first\:pt-0:first-child`,
		`.placeholder\:text-fg-placeholder::placeholder`,
	} {
		if !strings.Contains(css, want) {
			t.Errorf("rendered CSS missing %q", want)
		}
	}
}

func TestPeerStateFailsClosed(t *testing.T) {
	t.Parallel()
	if _, err := Rules("peer:underline"); err == nil {
		t.Fatal("peer-prefixed classes have no CSS mapping yet and must fail closed")
	}
}

func TestForDeduplicatesAcrossLists(t *testing.T) {
	t.Parallel()
	button := tw.New().Display(tw.DisplayInlineFlex).Gap(tw.S2).Bg(tw.SurfaceBrand)
	badge := tw.New().Display(tw.DisplayInlineFlex).Bg(tw.SurfaceBrandSoft)
	sheet, err := For(button, badge)
	if err != nil {
		t.Fatalf("For: %v", err)
	}
	css, err := sheet.Render(styleengine.RenderOptions{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if got := strings.Count(css, ".inline-flex"); got != 1 {
		t.Errorf("inline-flex emitted %d times, want 1", got)
	}
}

// TestEmittedSheetsRoundTrip proves everything emitted is CSS styleengine
// itself accepts, and regenerates the golden when -update is set.
func TestEmittedSheetsRoundTrip(t *testing.T) {
	t.Parallel()
	base, err := Base()
	if err != nil {
		t.Fatalf("Base: %v", err)
	}
	full := RoleVars().Merge(base)
	css, err := full.Render(styleengine.RenderOptions{Pretty: true})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if _, err := styleengine.Parse(css); err != nil {
		t.Fatalf("emitted CSS does not parse back: %v", err)
	}

	golden := "testdata/base.golden.css"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, []byte(css), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run with UPDATE_GOLDEN=1 to create): %v", err)
	}
	if string(want) != css {
		t.Fatalf("emitted CSS differs from golden (%d vs %d bytes); rerun with UPDATE_GOLDEN=1 and review the diff", len(css), len(want))
	}
}
