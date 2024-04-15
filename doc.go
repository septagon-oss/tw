// Package tw is a typed, allocation-efficient DSL for constructing
// Tailwind utility class strings from Go.
//
// ClassList is an immutable builder. Every modifier returns a new value.
// Compile walks the internal segment list once and emits a deterministic
// space-separated string. The intended use is to build component base
// classes at package init time or inside pure view functions; runtime
// cost after the first Compile is negligible.
//
// The package ships exhaustive typed enumerations for:
//   - layout (Display, Items, Justify, Position, FlexDir, ...)
//   - spacing and sizing (Spacing with S0..S64 + SPX, Width/Height/Min/Max variants, Gap...)
//   - semantic color roles (Surface*, Fg*, Border*, Ring* plus transparent/white/black)
//   - typography (FontSize, FontWeight, Tracking, Leading, FontFamily, TextAlign...)
//   - border, ring, shadow, radius, outline, opacity, cursor, ...
//   - motion (Transition, Duration, Easing, Translate...)
//   - state and responsive prefixes (State, Breakpoint) with func nesting
//   - z-index layers via ZLayer (arbitrary-value z-[N] emitted by .Class())
//   - plus many convenience methods (Truncate, SrOnly, LineClamp, Aspect*, GridCols, ...)
//
// Prefixing via On(state, fn) and Breakpoint(bp, fn) supports arbitrary
// nesting and stacking. Merge and Raw provide composition and escape hatches.
//
// A typed escape for non-Tailwind classes (your own CSS, design-system
// primitives, component data-handles, progress animations, etc.) is
// provided by the PlatformKitClass type and the PK method. Values are
// emitted verbatim.
//
// # Example
//
//	import "github.com/septagon-oss/tw"
//
//	base := tw.New().
//	    Display(tw.DisplayInlineFlex).
//	    Items(tw.ItemsCenter).
//	    Gap(tw.S2).
//	    Rounded(tw.RadiusXL).
//	    FontWeight(tw.FontSemibold).
//	    Bg(tw.SurfaceBrand).
//	    TextColor(tw.FgOnBrand).
//	    On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
//	        return c.Bg(tw.SurfaceBrandHover)
//	    }).
//	    Compile()
//
// All*() functions (AllColors, AllStates, AllRadii, AllZLayers, ...) are
// supplied for exhaustive tests and linter coverage of the compile tables.
//
// See the godoc for the full method and constant list.
package tw

// Implements: REQ-004. // Per: ADR-0004. // Discipline: C-14.
