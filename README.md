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

