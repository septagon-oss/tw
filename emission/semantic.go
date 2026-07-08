package emission

// Implements: REQ-004. // Per: ADR-0004. // Discipline: C-14.

import (
	"sort"

	"github.com/septagon-dev/platformkit-design-system/moduletokens"
	"github.com/septagon-dev/platformkit-design-system/tokens"
	"github.com/septagon-oss/styleengine"
)

// EmitSemanticTokens emits the complete semantic token surface: color
// tokens (light in :root, dark in .dark), spacing, radius, shadow, motion,
// z-index, typography (font-size, line-height, letter-spacing), opacity,
// and the three density tiers (comfortable → :root, compact →
// .density-compact, spacious → .density-spacious).
//
// Returns an empty Sheet when theme is nil.
func EmitSemanticTokens(theme tokens.Theme) *styleengine.Sheet {
	s := styleengine.NewSheet()
	if theme == nil {
		return s
	}
	semantic := tokens.BuildSemanticTokens(theme)

	// Light color tokens land in :root.
	for _, tok := range semantic.AllTokens() {
		s.Var(tok.Name, tok.Light)
	}

	emitNamedEntries(s, ":root", "space-", semantic.Spacing.Entries())
	emitNamedEntries(s, ":root", "radius-", semantic.Radius.Entries())
	emitNamedEntries(s, ":root", "shadow-", semantic.Shadow.Entries())
	emitNamedEntries(s, ":root", "motion-", semantic.Motion.Entries())
	emitNamedEntries(s, ":root", "z-", semantic.ZIndex.Entries())
	emitNamedEntries(s, ":root", "text-", semantic.Typography.Entries())
	emitNamedEntries(s, ":root", "leading-", semantic.Typography.LineHeightEntries())
	emitNamedEntries(s, ":root", "tracking-", semantic.Typography.LetterSpacingEntries())
	emitNamedEntries(s, ":root", "opacity-", semantic.Opacity.Entries())
	emitNamedEntries(s, ":root", "density-", semantic.Density.Comfortable.Entries())

	// Dark mode override block.
	darkDecls := make([]styleengine.Declaration, 0, len(semantic.AllTokens()))
	for _, tok := range semantic.AllTokens() {
		darkDecls = append(darkDecls, styleengine.Decl("--"+tok.Name, styleengine.Literal(tok.Dark)))
	}
	if len(darkDecls) > 0 {
		s.AddRule(styleengine.Rule{Selector: styleengine.MustSelector(".dark"), Decls: darkDecls})
	}

	// Density variant blocks.
	emitDensityVariant(s, ".density-compact", semantic.Density.Compact.Entries())
	emitDensityVariant(s, ".density-spacious", semantic.Density.Spacious.Entries())

	return s
}

// EmitModuleColors emits per-module custom color vars under :root. Each
// module's CustomColors map contributes `--<moduleID>-<colorName>: <rgb>`
// declarations. Hex values that fail to parse are skipped (a lint rule
// should catch them upstream).
func EmitModuleColors(catalog *moduletokens.Catalog) *styleengine.Sheet {
	s := styleengine.NewSheet()
	if catalog == nil {
		return s
	}
	for _, moduleID := range catalog.ModuleIDs() {
		module, _ := catalog.Get(moduleID)
		names := make([]string, 0, len(module.CustomColors))
		for name := range module.CustomColors {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			rgb, err := tokens.HexToRGB(module.CustomColors[name])
			if err != nil {
				continue
			}
			s.Var(moduleID+"-"+name, rgb)
		}
	}
	return s
}

// emitNamedEntries adds vars of shape `--{prefix}{name}: {value}` under
// the given selector. Token registries upstream must supply ident-clean
// keys (see styleengine.validVarName) — this emitter does NOT normalize.
func emitNamedEntries(s *styleengine.Sheet, selector, prefix string, entries []struct{ Name, Value string }) {
	if selector != ":root" {
		decls := make([]styleengine.Declaration, 0, len(entries))
		for _, e := range entries {
			decls = append(decls, styleengine.Decl("--"+prefix+e.Name, styleengine.Literal(e.Value)))
		}
		if len(decls) > 0 {
			s.AddRule(styleengine.Rule{Selector: styleengine.MustSelector(selector), Decls: decls})
		}
		return
	}
	for _, e := range entries {
		s.Var(prefix+e.Name, e.Value)
	}
}

// emitDensityVariant emits one of the density override blocks (compact or
// spacious). The selector applies the same --density-* prop names as :root
// but with override values.
func emitDensityVariant(s *styleengine.Sheet, selector string, entries []struct{ Name, Value string }) {
	if len(entries) == 0 {
		return
	}
	decls := make([]styleengine.Declaration, 0, len(entries))
	for _, e := range entries {
		decls = append(decls, styleengine.Decl("--density-"+e.Name, styleengine.Literal(e.Value)))
	}
	s.AddRule(styleengine.Rule{Selector: styleengine.MustSelector(selector), Decls: decls})
}
