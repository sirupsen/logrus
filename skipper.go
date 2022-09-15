package logrus

import (
	"regexp"
	"runtime"
)

var pkgRegex = regexp.MustCompile(`^(.+\/)?[^\. ]+\.?`)

// Skipper defines a behaviour to determine if a Frame should be skipped from
// the logs.
type Skipper interface {
	// ShouldSkip tells if the frame should be skipped.
	ShouldSkip(*runtime.Frame) bool
}

// CallerSkippers defines an internal type for storing skippers on a logger
// instance.
type CallerSkippers []Skipper

func (skippers CallerSkippers) with(s Skipper) CallerSkippers {
	return append(skippers, s)
}

func (skippers CallerSkippers) shouldSkip(f *runtime.Frame) bool {
	for _, skp := range skippers {
		if skp.ShouldSkip(f) {
			return true
		}
	}

	return false
}

// PackageSkipper defines a simple package skipper by their name.
type PackageSkipper struct {
	pkgName string
}

// NewPackageSkipper creates a skipper for the given package.
func NewPackageSkipper(funcName string) Skipper {
	return &PackageSkipper{
		pkgName: extractPackageName(funcName),
	}
}

func (ps *PackageSkipper) ShouldSkip(f *runtime.Frame) bool {
	pkgName := extractPackageName(f.Function)

	return pkgName == ps.pkgName
}

// extractPackageName reduces a fully qualified function name to the package
// name.
func extractPackageName(fn string) string {
	pkgName := pkgRegex.FindString(fn)
	if pkgName == "" {
		return fn
	}

	if pkgName[len(pkgName)-1] == '.' {
		return pkgName[:len(pkgName)-1]
	}

	return pkgName
}
