# tw

Typed builder for Tailwind utility class strings.

Package tw exports a small immutable fluent builder (ClassList) and a complete set of typed constants covering Tailwind's layout, spacing, sizing, color roles, typography, border, shadow, motion, state, breakpoint, z-index, and related utilities.

```go
import "github.com/septagon-oss/tw"

button := tw.New().
    Display(tw.DisplayInlineFlex).
    Items(tw.ItemsCenter).
    Justify(tw.JustifyCenter).
    Gap(tw.S2).
    PaddingX(tw.S4).
    PaddingY(tw.S2).
    Rounded(tw.RadiusXL).
    FontWeight(tw.FontSemibold).
    Bg(tw.SurfaceBrand).
    TextColor(tw.FgOnBrand).
    On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
        return c.Bg(tw.SurfaceBrandHover)
    }).
    Transition(tw.TransitionColors, tw.Duration200).
    Compile()
```

`Compile()` walks the accumulated segments once and returns a deterministic space-separated string suitable for gomponents, templ, or any HTML attribute writer. The builder allocates only on mutation; zero-value and empty compile to "".

All semantic colors (Surface*, Fg*, Border*, Ring*) are first-class; role methods (Bg, TextColor, BorderColor, RingColor) map them to the correct "bg-"/"text-"/... utility. Specials such as ColorTransparent and the plain ColorWhite/ColorBlack are supported.

Spacing, Radius, Shadow, Font*, Tracking, etc. follow the same pattern: pass the typed step to the role method.

Prefixing:

- `On(state, func(ClassList) ClassList)` wraps the inner result with "hover:", "focus-visible:", "group-hover:", etc. Nesting supported.
- `Breakpoint(bp, func...)` does the same for "sm:", "lg:", ...

Composition:

- `Merge(other)` appends another builder's segments.
- `Raw(s)` passes a pre-validated utility string through (use for runtime values or migration only).

Non-Tailwind classes (custom CSS, component handles, pk-* animations, admin chrome) are routed through the PlatformKitClass type and the `PK(c)` method so they bypass any Tailwind-only linters.

Enumerators (AllColors, AllStates, AllRadii, AllZLayers, ...) exist for exhaustive testing and coverage tooling.

See godoc for the complete method and constant set. The package has comprehensive tests and executable Examples.

The package is used as the single source of Tailwind class construction for component libraries and is intended to be generally useful anywhere a typed, zero-static-string Tailwind DSL is desired in Go.
