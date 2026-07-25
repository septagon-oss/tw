// Implements: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

// Package emission renders CSS for the classes tw compiles, against the
// PlatformKit design system's custom properties.
//
// The stack composes in three layers, each replaceable:
//
//	:root { --pk-* }        pk-design theme tokens — rendered by the consumer
//	                        (pk-design tokens.CSSVars), not by this package
//	:root { --pk-role-* }   RoleVars(): tw's semantic color roles mapped onto
//	                        the token variables (or color-mix derivations)
//	utility rules           Base() for the full enumerable universe, or
//	                        For()/Rules() for exactly the classes an app uses
//
// This package deliberately does not import pk-design at runtime: role
// variables reference --pk-color-* custom properties by name, and the naming
// contract is pinned by a test-only dependency. Consumers on a different
// token pipeline can define those properties however they like.
//
// Escape hatches stay escape hatches: Raw classes, the *Raw methods, and
// PlatformKitClass handles are not resolvable here — the caller owns that
// CSS, and Rules fails closed on them rather than guessing.
package emission

import (
	"sort"
	"strings"

	"github.com/septagon-oss/styleengine"
	"github.com/septagon-oss/tw"
)

// RoleVars emits the :root block mapping every tw color role to the design
// system's token variables, plus the keyframes the animate-* utilities
// reference. Serve it after the theme's token block and before the utilities.
func RoleVars() *styleengine.Sheet {
	s := styleengine.NewSheet()
	roles := roleValues()
	names := make([]string, 0, len(roles))
	for c := range roles {
		names = append(names, string(c))
	}
	sort.Strings(names)
	for _, n := range names {
		s.Var("pk-role-"+n, roles[tw.Color(n)])
	}
	s.Keyframes("pk-spin", func(kb *styleengine.KeyframesBuilder) {
		kb.At("from", styleengine.Decl("transform", styleengine.Literal("rotate(0deg)")))
		kb.At("to", styleengine.Decl("transform", styleengine.Literal("rotate(360deg)")))
	})
	s.Keyframes("pk-pulse", func(kb *styleengine.KeyframesBuilder) {
		kb.At("0%, 100%", styleengine.Decl("opacity", styleengine.Literal("1")))
		kb.At("50%", styleengine.Decl("opacity", styleengine.Literal("0.5")))
	})
	return s
}

// Base emits one rule for every class in tw's enumerable universe — the
// closure of the typed enumerations and no-argument toggles over their
// builder methods, plus the bounded integer families. Prefixed variants
// (hover:, md:) are not pre-generated; use For or Rules for the exact
// variants an application composes.
func Base() (*styleengine.Sheet, error) {
	return Rules(baseClasses()...)
}

// For emits the rules for exactly the classes the given lists compile to,
// prefixed variants included. This is the deterministic, Go-native
// alternative to source scanning: components declare their ClassLists, and
// the stylesheet is derived from those declarations.
func For(lists ...tw.ClassList) (*styleengine.Sheet, error) {
	seen := map[string]bool{}
	var classes []string
	for _, l := range lists {
		for _, class := range strings.Fields(l.Compile()) {
			if !seen[class] {
				seen[class] = true
				classes = append(classes, class)
			}
		}
	}
	sort.Strings(classes)
	return Rules(classes...)
}

// baseClasses enumerates the class universe by driving tw's own enumerators
// and builder methods, so coverage tracks tw by construction.
func baseClasses() []string {
	var out []string
	add := func(cl tw.ClassList) {
		if c := cl.Compile(); c != "" { // a zero-valued constant may compile to nothing
			out = append(out, c)
		}
	}

	colorUniverse := append(tw.AllColors(), tw.ColorWhite, tw.ColorBlack)
	for _, c := range colorUniverse {
		add(tw.New().Bg(c))
		add(tw.New().TextColor(c))
		add(tw.New().BorderColor(c))
		add(tw.New().RingColor(c))
		add(tw.New().Accent(c))
	}
	for _, s := range tw.AllSpacings() {
		add(tw.New().Padding(s))
		add(tw.New().PaddingX(s))
		add(tw.New().PaddingY(s))
		add(tw.New().PaddingLeft(s))
		add(tw.New().PaddingRight(s))
		add(tw.New().PaddingTop(s))
		add(tw.New().PaddingBottom(s))
		add(tw.New().Margin(s))
		add(tw.New().MarginX(s))
		add(tw.New().MarginY(s))
		add(tw.New().MarginLeft(s))
		add(tw.New().MarginRight(s))
		add(tw.New().MarginTop(s))
		add(tw.New().MarginBottom(s))
		add(tw.New().NegTop(s))
		add(tw.New().NegBottom(s))
		add(tw.New().NegLeft(s))
		add(tw.New().NegRight(s))
		add(tw.New().Gap(s))
		add(tw.New().GapX(s))
		add(tw.New().GapY(s))
		add(tw.New().Top(s))
		add(tw.New().Bottom(s))
		add(tw.New().Left(s))
		add(tw.New().Right(s))
		add(tw.New().Inset(s))
		add(tw.New().InsetX(s))
		add(tw.New().InsetY(s))
		add(tw.New().Width(s))
		add(tw.New().Height(s))
		add(tw.New().MinWidth(s))
		add(tw.New().MinHeight(s))
		add(tw.New().DivideX(s))
		add(tw.New().DivideY(s))
		add(tw.New().SpaceX(s))
		add(tw.New().SpaceY(s))
		add(tw.New().UnderlineOffset(s))
		add(tw.New().TranslateXStep(s))
		add(tw.New().TranslateYStep(s))
	}
	for _, v := range tw.AllTranslates() {
		add(tw.New().TranslateX(v))
		add(tw.New().TranslateY(v))
	}
	for _, v := range tw.AllDisplays() {
		add(tw.New().Display(v))
	}
	for _, v := range tw.AllItems() {
		add(tw.New().Items(v))
	}
	for _, v := range tw.AllJustify() {
		add(tw.New().Justify(v))
	}
	for _, v := range tw.AllFlexDirs() {
		add(tw.New().FlexDir(v))
	}
	for _, v := range tw.AllPositions() {
		add(tw.New().Position(v))
	}
	for _, v := range tw.AllOverflows() {
		add(tw.New().Overflow(v))
		add(tw.New().OverflowX(v))
		add(tw.New().OverflowY(v))
	}
	for _, v := range tw.AllFontSizes() {
		add(tw.New().FontSize(v))
	}
	for _, v := range tw.AllFontWeights() {
		add(tw.New().FontWeight(v))
	}
	for _, v := range tw.AllFontFamilies() {
		add(tw.New().FontFamily(v))
	}
	for _, v := range tw.AllTrackings() {
		add(tw.New().Tracking(v))
	}
	for _, v := range tw.AllLeadings() {
		add(tw.New().Leading(v))
	}
	for _, v := range tw.AllTextAligns() {
		add(tw.New().TextAlign(v))
	}
	for _, v := range tw.AllRadii() {
		add(tw.New().Rounded(v))
		add(tw.New().RoundedTop(v))
		add(tw.New().RoundedBottom(v))
		add(tw.New().RoundedLeft(v))
		add(tw.New().RoundedRight(v))
	}
	for _, v := range tw.AllShadows() {
		add(tw.New().Shadow(v))
	}
	for _, v := range tw.AllBorderWidths() {
		add(tw.New().Border(v))
		add(tw.New().BorderTop(v))
		add(tw.New().BorderBottom(v))
		add(tw.New().BorderLeft(v))
		add(tw.New().BorderRight(v))
	}
	for _, v := range tw.AllBorderStyles() {
		add(tw.New().BorderStyle(v))
	}
	for _, v := range tw.AllRingWidths() {
		add(tw.New().Ring(v))
	}
	for _, v := range tw.AllRingOffsets() {
		add(tw.New().RingOffset(v))
	}
	for _, v := range tw.AllOpacities() {
		add(tw.New().Opacity(v))
	}
	for _, v := range tw.AllCursors() {
		add(tw.New().Cursor(v))
	}
	for _, v := range tw.AllPointerEvents() {
		add(tw.New().PointerEvents(v))
	}
	for _, v := range tw.AllOutlines() {
		add(tw.New().Outline(v))
	}
	for _, v := range tw.AllTransitions() {
		add(tw.New().Transition(v))
	}
	for _, v := range tw.AllDurations() {
		add(tw.New().Duration(v))
	}
	for _, v := range tw.AllEasings() {
		add(tw.New().Easing(v))
	}
	for _, v := range tw.AllSelects() {
		add(tw.New().UserSelect(v))
	}
	// MaxWidth has no All* enumerator in tw; enumerate the constants here and
	// let the coverage meta-test flag any future addition via AllSizes parity.
	for _, v := range []tw.MaxWidth{
		tw.MaxWXS, tw.MaxWSM, tw.MaxWMD, tw.MaxWLG, tw.MaxWXL, tw.MaxW2XL,
		tw.MaxW3XL, tw.MaxW4XL, tw.MaxW5XL, tw.MaxW6XL, tw.MaxW7XL,
		tw.MaxWFull, tw.MaxWNone, tw.MaxWScreen, tw.MaxWProse,
	} {
		add(tw.New().MaxWScaled(v))
	}
	for _, v := range tw.AllZLayers() {
		add(tw.New().ZLayer(v))
	}
	for n := 1; n <= 12; n++ {
		add(tw.New().GridCols(n))
		add(tw.New().ColSpan(n))
	}
	add(tw.New().ColSpanFull())
	for n := 1; n <= 10; n++ {
		add(tw.New().LineClamp(n))
	}

	// No-argument toggles.
	add(tw.New().Truncate())
	add(tw.New().SrOnly())
	add(tw.New().Italic())
	add(tw.New().NotItalic())
	add(tw.New().Underline())
	add(tw.New().NoUnderline())
	add(tw.New().LineThrough())
	add(tw.New().NoLineThrough())
	add(tw.New().Uppercase())
	add(tw.New().Lowercase())
	add(tw.New().Capitalize())
	add(tw.New().NormalCase())
	add(tw.New().TabularNums())
	add(tw.New().WhitespaceNowrap())
	add(tw.New().WhitespacePreWrap())
	add(tw.New().BreakAll())
	add(tw.New().BreakWords())
	add(tw.New().Flex1())
	add(tw.New().FlexGrow())
	add(tw.New().FlexGrow0())
	add(tw.New().FlexShrink0())
	add(tw.New().FlexNone())
	add(tw.New().FlexWrap())
	add(tw.New().FlexNoWrap())
	add(tw.New().Group())
	add(tw.New().Peer())
	add(tw.New().Relative())
	add(tw.New().Transform())
	add(tw.New().AnimateSpin())
	add(tw.New().AnimatePulse())
	add(tw.New().AppearanceNone())
	add(tw.New().AspectSquare())
	add(tw.New().AspectVideo())
	add(tw.New().ObjectContain())
	add(tw.New().ObjectCover())
	add(tw.New().OverscrollContain())
	add(tw.New().ResizeNone())
	add(tw.New().ResizeY())
	add(tw.New().RingInset())

	sort.Strings(out)
	return dedupe(out)
}

func dedupe(in []string) []string {
	out := in[:0]
	var prev string
	for i, s := range in {
		if i == 0 || s != prev {
			out = append(out, s)
		}
		prev = s
	}
	return out
}
