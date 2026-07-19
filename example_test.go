package tw_test

// Validates: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

import (
	"fmt"

	"github.com/septagon-oss/tw"
)

func ExampleClassList_basic() {
	cls := tw.New().
		Display(tw.DisplayFlex).
		Items(tw.ItemsCenter).
		Gap(tw.S2).
		PaddingX(tw.S4).
		Rounded(tw.RadiusXL).
		FontWeight(tw.FontSemibold).
		Compile()
	fmt.Println(cls)
	// Output: flex items-center gap-2 px-4 rounded-xl font-semibold
}

func ExampleClassList_colorsAndStates() {
	cls := tw.New().
		Bg(tw.SurfacePrimary).
		TextColor(tw.FgPrimary).
		On(tw.StateHover, func(c tw.ClassList) tw.ClassList {
			return c.Bg(tw.SurfaceHover).TextColor(tw.FgPrimary)
		}).
		Compile()
	fmt.Println(cls)
	// Output: bg-surface-primary text-fg-primary hover:bg-surface-hover hover:text-fg-primary
}

func ExampleClassList_breakpointAndMerge() {
	base := tw.New().Padding(tw.S4).Rounded(tw.RadiusMD)
	responsive := base.Merge(
		tw.New().Breakpoint(tw.BreakpointMD, func(c tw.ClassList) tw.ClassList {
			return c.Padding(tw.S6)
		}),
	)
	fmt.Println(responsive.Compile())
	// Output: p-4 rounded-md md:p-6
}

func ExampleClassList_zLayerAndCustom() {
	// ZLayer produces arbitrary-value z-* classes from typed numeric layers.
	modal := tw.New().ZLayer(tw.ZModal).Compile()
	fmt.Println(modal)

	// PlatformKitClass (or your own typed alias) lets non-Tailwind classes
	// flow through the same builder without triggering utility-only linters.
	custom := tw.New().PK(tw.PlatformKitClass("my-custom-handle")).Compile()
	fmt.Println(custom)
	// Output:
	// z-[1400]
	// my-custom-handle
}

func ExampleClassList_rawAndCompileEmpty() {
	// Raw is the escape hatch for runtime-computed or legacy values.
	// Prefer the typed API for all static utilities.
	hybrid := tw.New().Flex1().Raw("data-active").Compile()
	fmt.Printf("%q\n", hybrid)

	empty := tw.New().Compile()
	fmt.Printf("empty=%q\n", empty)
	// Output:
	// "flex-1 data-active"
	// empty=""
}
