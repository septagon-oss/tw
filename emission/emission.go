// Package emission translates typed design-system tokens into styleengine
// rules. Each Emit* function returns a *styleengine.Sheet that can be merged
// into a larger bundle by the cssbundle orchestrator.
//
// The functions here are deliberately stateless: given the same input slice
// they produce identical output, enabling deterministic golden testing in
// the design-system and downstream consumers.
package emission

// Implements: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

import (
	"strings"

	"github.com/septagon-dev/platformkit-design-system/tokens/motion"
	"github.com/septagon-dev/platformkit-design-system/tokens/typography"
	"github.com/septagon-oss/styleengine"
)

// stripVarPrefix removes the leading "--" from a CSS variable name so it can
// be passed to styleengine.Sheet.Var (which adds the prefix back).
func stripVarPrefix(name string) string {
	return strings.TrimPrefix(name, "--")
}

// EmitFontStacks emits :root --font-* vars for every registered font stack.
func EmitFontStacks(stacks []typography.FontStack) *styleengine.Sheet {
	s := styleengine.NewSheet()
	for _, fs := range stacks {
		s.Var(stripVarPrefix(typography.CSSVarName(fs)), typography.StackCSSValue(fs))
	}
	return s
}

// EmitFontScales emits :root --scale-*-size, --scale-*-line, --scale-*-track
// vars for every registered Scale.
func EmitFontScales(scales []typography.Scale) *styleengine.Sheet {
	s := styleengine.NewSheet()
	for _, sc := range scales {
		s.Var(stripVarPrefix(typography.SizeVar(sc)), sc.Size)
		s.Var(stripVarPrefix(typography.LineVar(sc)), sc.Line)
		s.Var(stripVarPrefix(typography.TrackVar(sc)), sc.Track)
	}
	return s
}

// EmitMotionDurations emits :root --motion-duration-* vars for every
// registered Duration. The CSS value is formatted as "{Ms}ms".
func EmitMotionDurations(durations []motion.Duration) *styleengine.Sheet {
	s := styleengine.NewSheet()
	for _, d := range durations {
		s.Var(stripVarPrefix(motion.DurationCSSVarName(d)), formatMs(d.Ms))
	}
	return s
}

// EmitMotionEasings emits :root --motion-easing-* vars for every Easing.
func EmitMotionEasings(easings []motion.Easing) *styleengine.Sheet {
	s := styleengine.NewSheet()
	for _, e := range easings {
		s.Var(stripVarPrefix(motion.EasingCSSVarName(e)), e.CSS)
	}
	return s
}

// EmitMotionSprings emits :root --motion-spring-* vars for every Spring.
// Springs are linear() timing functions (spring physics without JS) and are
// consumed as transition/animation timing-functions, exactly like easings.
func EmitMotionSprings(springs []motion.Spring) *styleengine.Sheet {
	s := styleengine.NewSheet()
	for _, sp := range springs {
		s.Var(stripVarPrefix(motion.SpringCSSVarName(sp)), sp.CSS)
	}
	return s
}

// EmitReducedMotion emits @media (prefers-reduced-motion: reduce) { :root {
// --motion-duration-*: 0ms; } } for every registered Duration. Tenants who
// opt out of reduced-motion can override individual durations downstream.
func EmitReducedMotion(durations []motion.Duration) *styleengine.Sheet {
	s := styleengine.NewSheet()
	s.Media("(prefers-reduced-motion: reduce)", func(inner *styleengine.Sheet) {
		decls := make([]styleengine.Declaration, 0, len(durations))
		for _, d := range durations {
			decls = append(decls, styleengine.Decl(
				motion.DurationCSSVarName(d),
				styleengine.Literal("0ms"),
			))
		}
		inner.AddRule(styleengine.Rule{
			Selector: styleengine.MustSelector(":root"),
			Decls:    decls,
		})
	})
	return s
}

// formatMs renders an integer millisecond count as "Nms" without relying on
// fmt (single allocation, hot-path friendly).
func formatMs(ms int) string {
	return itoa(ms) + "ms"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
