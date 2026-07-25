// Implements: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

package emission

// roles.go owns the single mapping from tw's semantic color roles to the
// design system's theme tokens. Utility rules never reference theme tokens
// directly: they reference --pk-role-* variables, and this file emits those
// variables from --pk-* token variables (or derives them with color-mix when
// the theme has no dedicated token). Retheming therefore never touches the
// utility rules — a different theme changes the values behind the same roles.

import (
	"github.com/septagon-oss/tw"
)

// tokenVar renders a var() reference to a pk-design token custom property.
func tokenVar(path string) string { return "var(--pk-color-" + path + ")" }

// mix derives a tint: pct% of colorA over colorB, in sRGB. Used for the soft,
// hover, and disabled roles the theme intentionally does not enumerate.
func mix(a string, pct string, b string) string {
	return "color-mix(in srgb, " + a + " " + pct + "%, " + b + ")"
}

// roleValues maps every tw.Color to its CSS value in terms of theme token
// variables. TestRoleMapCoversEveryColor pins this to tw.AllColors(), so a new
// role in tw fails the build of this package's tests until it is mapped here.
func roleValues() map[tw.Color]string {
	surfacePrimary := tokenVar("surface-primary")
	textPrimary := tokenVar("text-primary")
	textMuted := tokenVar("text-muted")
	accent := tokenVar("accent-default")
	focus := tokenVar("focus")

	return map[tw.Color]string{
		// Surfaces.
		tw.SurfacePrimary:     surfacePrimary,
		tw.SurfaceSecondary:   tokenVar("surface-canvas"),
		tw.SurfaceTertiary:    tokenVar("surface-muted"),
		tw.SurfaceBrand:       accent,
		tw.SurfaceBrandHover:  tokenVar("accent-hover"),
		tw.SurfaceBrandSoft:   mix(accent, "12", surfacePrimary),
		tw.SurfaceSuccess:     tokenVar("status-ok"),
		tw.SurfaceSuccessSoft: tokenVar("status-okbg"),
		tw.SurfaceWarning:     tokenVar("status-warning"),
		tw.SurfaceWarningSoft: tokenVar("status-warningbg"),
		tw.SurfaceDanger:      tokenVar("status-danger"),
		tw.SurfaceDangerSoft:  tokenVar("status-dangerbg"),
		tw.SurfaceInfo:        focus,
		tw.SurfaceInfoSoft:    mix(focus, "12", surfacePrimary),
		tw.SurfaceDisabled:    tokenVar("surface-muted"),
		tw.SurfaceHover:       mix(textPrimary, "4", surfacePrimary),
		tw.SurfaceActive:      mix(textPrimary, "8", surfacePrimary),
		tw.SurfaceOverlay:     mix(tokenVar("sidebar-bg"), "55", "transparent"),
		tw.SurfaceInverse:     tokenVar("sidebar-bg"),

		// Foreground.
		tw.FgPrimary:     textPrimary,
		tw.FgSecondary:   mix(textPrimary, "78", surfacePrimary),
		tw.FgTertiary:    mix(textPrimary, "60", surfacePrimary),
		tw.FgMuted:       textMuted,
		tw.FgPlaceholder: mix(textMuted, "70", surfacePrimary),
		tw.FgBrand:       accent,
		tw.FgOnBrand:     tokenVar("accent-on"),
		tw.FgSuccess:     tokenVar("status-ok"),
		tw.FgWarning:     tokenVar("status-warning"),
		tw.FgDanger:      tokenVar("status-danger"),
		tw.FgInfo:        focus,
		tw.FgDisabled:    mix(textMuted, "55", surfacePrimary),
		tw.FgOnSurface:   textPrimary,
		tw.FgOnInverse:   tokenVar("sidebar-text"),
		tw.FgLink:        accent,
		tw.FgLinkHover:   tokenVar("accent-hover"),

		// Borders.
		tw.BorderPrimary:   tokenVar("border-default"),
		tw.BorderSecondary: mix(tokenVar("border-default"), "60", surfacePrimary),
		tw.BorderBrand:     accent,
		tw.BorderSuccess:   tokenVar("status-ok"),
		tw.BorderWarning:   tokenVar("status-warning"),
		tw.BorderDanger:    tokenVar("status-danger"),
		tw.BorderInfo:      focus,

		// Rings.
		tw.RingBrand:  accent,
		tw.RingFocus:  focus,
		tw.RingDanger: tokenVar("status-danger"),

		// Neutrals.
		tw.ColorTransparent: "transparent",
		tw.ColorWhite:       "#ffffff",
		tw.ColorBlack:       "#000000",
	}
}

// roleVar renders the var() reference utility rules use for a color role.
func roleVar(c tw.Color) string { return "var(--pk-role-" + string(c) + ")" }
