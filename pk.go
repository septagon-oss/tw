package tw

// Implements: REQ-004.
// Per: ADR-0004.
// Discipline: C-14.

// PlatformKitClass is a typed handle for application-specific or
// custom CSS classes that live outside Tailwind's utility namespace.
// Typical uses: your own design-system primitives, component handles
// (data-controller targets, htmx attributes surfaced as classes),
// progress/keyframe animations, or admin-shell chrome defined in a
// separate stylesheet. The values are emitted verbatim by the PK
// method; no "tw-" rewriting is applied.
type PlatformKitClass string

// Example non-Tailwind classes. Consumers are expected to define
// their own values (or alias these) that match classes present in
// their CSS bundles. The names and values here are retained for
// compatibility with existing call sites; new consumers should
// introduce their own typed constants of this type.
const (
	// Motion primitives (shared across multiple components).
	PKTransitionStandard PlatformKitClass = "pk-transition-standard"
	PKTransitionColors   PlatformKitClass = "pk-transition-colors"
	PKTransitionOpacity  PlatformKitClass = "pk-transition-opacity"
	PKTransitionProgress PlatformKitClass = "pk-transition-progress"
	// PKTransitionTransform is the transform-tier transition treatment
	// used by disclosure chevrons and other geometry-only transitions.
	PKTransitionTransform PlatformKitClass = "pk-transition-transform"
	// PKTransitionEmphasis is the emphasis-tier transition treatment
	// used by modal / editor surfaces that require a more pronounced
	// motion curve.
	PKTransitionEmphasis PlatformKitClass = "pk-transition-emphasis"
	// PKMotionModerate tunes the motion intensity for components that
	// want a "moderate" animation setting without re-declaring the
	// transition property group.
	PKMotionModerate PlatformKitClass = "pk-motion-moderate"
	// PKMotionSlow pairs with PKMotionModerate — used by page-level
	// progress surfaces (wizards, long-running flows) where the
	// animation curve should feel gentler than the default.
	PKMotionSlow PlatformKitClass = "pk-motion-slow"

	// Progress-bar animations.
	PKProgressFill              PlatformKitClass = "pk-progress-fill"
	PKProgressFillIndeterminate PlatformKitClass = "pk-progress-fill-indeterminate"
	AnimatePKProgressBar        PlatformKitClass = "animate-pk-progress-bar"

	// Component handles (data-attribute selectors used by HTMX /
	// Stimulus controllers and by custom stylesheets).
	PKHandleChatFab PlatformKitClass = "chat-fab"

	// Example admin / shell handles. These demonstrate the pattern
	// for routing non-Tailwind class names through the typed builder.
	PKAdminSidebarShell        PlatformKitClass = "admin-sidebar-shell"
	PKAdminSidebarBrand        PlatformKitClass = "admin-sidebar-brand"
	PKAdminSidebarBrandMark    PlatformKitClass = "admin-sidebar-brand-mark"
	PKAdminSidebarBrandCopy    PlatformKitClass = "admin-sidebar-brand-copy"
	PKAdminSidebarBrandEyebrow PlatformKitClass = "admin-sidebar-brand-eyebrow"
	PKAdminSidebarBrandTitle   PlatformKitClass = "admin-sidebar-brand-title"
	PKAdminSidebarBrandText    PlatformKitClass = "admin-sidebar-brand-text"
	PKAdminSidebarChevron      PlatformKitClass = "admin-sidebar-chevron"
	PKAdminSidebarSectionTitle PlatformKitClass = "admin-sidebar-section-title"
	PKAdminSidebarLink         PlatformKitClass = "admin-sidebar-link"
	PKAdminSidebarDock         PlatformKitClass = "admin-sidebar-dock"
	PKAdminSidebarTenantSwitch PlatformKitClass = "admin-sidebar-tenant-switcher"
	PKAdminSidebarTenantLabel  PlatformKitClass = "admin-sidebar-tenant-label"
	PKAdminSidebarUserMenu     PlatformKitClass = "admin-sidebar-user-menu"
	PKAdminTopbar              PlatformKitClass = "admin-topbar"
	PKAdminTopbarShell         PlatformKitClass = "admin-topbar-shell"
	PKAdminTopbarContext       PlatformKitClass = "admin-topbar-context"
	PKAdminTopbarKicker        PlatformKitClass = "admin-topbar-kicker"
	PKAdminTopbarHeading       PlatformKitClass = "admin-topbar-heading"
	PKAdminTopbarTitle         PlatformKitClass = "admin-topbar-title"
	PKAdminTopbarBadge         PlatformKitClass = "admin-topbar-badge"
	PKAdminTopbarSubtitle      PlatformKitClass = "admin-topbar-subtitle"
	PKAdminTopbarActions       PlatformKitClass = "admin-topbar-actions"
	PKAdminTopbarTenantSwitch  PlatformKitClass = "admin-topbar-tenant-switcher"
	PKAdminTopbarTenantLabel   PlatformKitClass = "admin-topbar-tenant-label"
	PKAdminTopbarNotifications PlatformKitClass = "admin-topbar-notifications"
	PKAdminTopbarUserMenu      PlatformKitClass = "admin-topbar-user-menu"
	PKAdminTopbarThemeToggle   PlatformKitClass = "admin-topbar-theme-toggle"
	PKAdminToolbarSearch       PlatformKitClass = "admin-toolbar-search"
	PKAdminToolbarShortcut     PlatformKitClass = "admin-toolbar-shortcut"

	// State-attribute classes (is-active, is-idle, is-parent, is-open).
	// They are not Tailwind utilities; include them when your CSS
	// responds to these state markers.
	PKStateActive PlatformKitClass = "is-active"
	PKStateIdle   PlatformKitClass = "is-idle"
	PKStateParent PlatformKitClass = "is-parent"
	PKStateOpen   PlatformKitClass = "is-open"
)

// AllPlatformKitClasses returns every PlatformKitClass const in stable
// order. Useful for exhaustive coverage tests of the escape-hatch path.
func AllPlatformKitClasses() []PlatformKitClass {
	return []PlatformKitClass{
		PKTransitionStandard, PKTransitionColors, PKTransitionOpacity, PKTransitionProgress,
		PKTransitionTransform, PKTransitionEmphasis, PKMotionModerate, PKMotionSlow,
		PKProgressFill, PKProgressFillIndeterminate, AnimatePKProgressBar,
		PKHandleChatFab,
		PKAdminSidebarShell, PKAdminSidebarBrand, PKAdminSidebarBrandMark, PKAdminSidebarBrandCopy,
		PKAdminSidebarBrandEyebrow, PKAdminSidebarBrandTitle, PKAdminSidebarBrandText,
		PKAdminSidebarChevron, PKAdminSidebarSectionTitle, PKAdminSidebarLink,
		PKAdminSidebarDock, PKAdminSidebarTenantSwitch, PKAdminSidebarTenantLabel, PKAdminSidebarUserMenu,
		PKAdminTopbar, PKAdminTopbarShell, PKAdminTopbarContext, PKAdminTopbarKicker,
		PKAdminTopbarHeading, PKAdminTopbarTitle, PKAdminTopbarBadge, PKAdminTopbarSubtitle,
		PKAdminTopbarActions, PKAdminTopbarTenantSwitch, PKAdminTopbarTenantLabel,
		PKAdminTopbarNotifications, PKAdminTopbarUserMenu, PKAdminTopbarThemeToggle,
		PKAdminToolbarSearch, PKAdminToolbarShortcut,
		PKStateActive, PKStateIdle, PKStateParent, PKStateOpen,
	}
}

// PK appends a custom (non-Tailwind) class to the ClassList. The
// value is emitted verbatim. Use this for anything that should not
// be interpreted as a Tailwind utility (e.g. your "pk-*" animation
// classes, admin chrome handles, or arbitrary component markers).
func (cl ClassList) PK(c PlatformKitClass) ClassList {
	if c == "" {
		return cl
	}
	return cl.append(string(c))
}
