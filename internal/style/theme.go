package style

import (
	"flag"
	"regexp"
	"strings"
)

var (
	largeToolbarForce        = flag.Bool("force-large-toolbar-icons", false, "Force use of large icons in the toolbar.")
	largeToolbarThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
		"win32",
	}, "|"))

	nonSymbolicIconForce        = flag.Bool("force-non-symbolic-icons", false, "Force use of non-symbolic icons in the toolbar.")
	nonSymbolicIconThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
	}, "|"))

	unlinkedNavButtonsForce        = flag.Bool("force-unlinked-nav-buttons", false, "Force unlinked styling for navigation buttons.")
	unlinkedNavButtonsThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
	}, "|"))

	compactMenuForce        = flag.Bool("force-compact-menu", false, "Force the window menu to use compact styling.")
	compactMenuThemesRegexp = regexp.MustCompile(strings.Join([]string{
		"elementary(-x)?",
		"io\\.elementary\\.stylesheet.*",
	}, "|"))
)

// IsLargeToolbarTheme returns true if we should use large toolbar buttons with
// the given theme.
func IsLargeToolbarTheme(theme string) bool {
	return *largeToolbarForce || largeToolbarThemesRegexp.MatchString(theme)
}

// IsSymbolicIconTheme returns true if we should use symbolic icons with the
// given theme.
func IsSymbolicIconTheme(theme string, darkMode bool) bool {
	return !*nonSymbolicIconForce && (darkMode || !nonSymbolicIconThemesRegexp.MatchString(theme))
}

// IsLinkedNavButtonsTheme returns true if we should visually "link" the buttons
// in the navigation button box for the given theme.
func IsLinkedNavButtonsTheme(theme string) bool {
	return !*unlinkedNavButtonsForce && !unlinkedNavButtonsThemesRegexp.MatchString(theme)
}

// IsCompactMenuTheme returns true if we should reduce the left and right
// margins of popover-style menus.
func IsCompactMenuTheme(theme string) bool {
	return *compactMenuForce || compactMenuThemesRegexp.MatchString(theme)
}
