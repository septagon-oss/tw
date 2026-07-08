package emission

// Implements: REQ-004. // Per: ADR-0004. // Discipline: C-14.

import (
	"regexp"

	"github.com/septagon-dev/platformkit-design-system/experiences"
	"github.com/septagon-oss/styleengine"
)

// surfaceClassRE constrains the shape of class-name surface tokens
// (RootBackground, RootGrain). The template layer applies these as body
// classes; values that fail this regex are dropped at emission time.
// Matching the CSS class-ident shape: lowercase letter start, then
// [a-z0-9-]. Underscores are excluded so the convention matches Tailwind
// and the existing pke-* class space.
var surfaceClassRE = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// ValidSurfaceClass reports whether a surface class-name token is shape-valid.
func ValidSurfaceClass(s string) bool { return surfaceClassRE.MatchString(s) }

// ValidateSurfaceClassTokens inspects a pack's class-name surface tokens
// (RootBackground, RootGrain) and returns the slice of token names that
// failed the [ValidSurfaceClass] shape check. Empty strings are ignored
// (an unset token is legitimate). Callers that build a pack catalog use
// this to fail-fast on pack registration; cssbundle.Compile does not
// re-run the check because Compile trusts in-tree pack authors.
func ValidateSurfaceClassTokens(surface *experiences.SurfaceTokens) []string {
	if surface == nil {
		return nil
	}
	var bad []string
	for _, pair := range []struct{ name, value string }{
		{"RootBackground", surface.RootBackground},
		{"RootGrain", surface.RootGrain},
	} {
		if pair.value != "" && !ValidSurfaceClass(pair.value) {
			bad = append(bad, pair.name+"="+pair.value)
		}
	}
	return bad
}

// EmitSurfaceTokens emits CSS for the body-applied portion of a pack's
// [experiences.SurfaceTokens]: cursor (default + hover) and scrollbar
// treatment. Class-name tokens (RootBackground, RootGrain) are NOT emitted
// here — they are applied by the template layer as body classes, not CSS.
//
// When surface is nil or carries only class-name tokens, the returned Sheet
// is empty. Tenant SurfaceOverrideV1 values must be merged INTO the
// SurfaceTokens via [overlays.ResolveSurface] before calling this.
func EmitSurfaceTokens(surface *experiences.SurfaceTokens) *styleengine.Sheet {
	s := styleengine.NewSheet()
	if surface == nil {
		return s
	}
	bodyDecls := make([]styleengine.Declaration, 0, 2)
	if surface.RootCursor != "" {
		bodyDecls = append(bodyDecls, styleengine.Decl("cursor", styleengine.Literal(surface.RootCursor)))
	}
	if v := scrollbarWidthCSS(surface.RootScrollbar); v != "" {
		bodyDecls = append(bodyDecls, styleengine.Decl("scrollbar-width", styleengine.Literal(v)))
	}
	if len(bodyDecls) > 0 {
		s.AddRule(styleengine.Rule{
			Selector: styleengine.MustSelector("body"),
			Decls:    bodyDecls,
		})
	}
	if surface.RootCursorHover != "" {
		s.AddRule(styleengine.Rule{
			Selector: styleengine.MustSelector("body a:hover, body button:hover, body [role=\"button\"]:hover"),
			Decls:    []styleengine.Declaration{styleengine.Decl("cursor", styleengine.Literal(surface.RootCursorHover))},
		})
	}
	return s
}

// scrollbarWidthCSS normalizes a SurfaceTokens.RootScrollbar token into a
// valid CSS scrollbar-width value. The CSS spec only accepts `auto`,
// `thin`, and `none`; PK's `default` token signals "keep UA default",
// which means we emit nothing. Unknown tokens are dropped (returning "")
// so a typo cannot leak invalid CSS to the browser — a lint rule should
// catch them upstream.
func scrollbarWidthCSS(token string) string {
	switch token {
	case "auto", "thin", "none":
		return token
	default:
		return ""
	}
}
