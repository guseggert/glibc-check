package glibccheck

import (
	"debug/elf"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PaesslerAG/gval"
)

var r = regexp.MustCompile(`^GLIBC_([0-9.]+)$`)

type GLIBCVersion struct {
	Full  string
	Major int
	Minor int
	Patch int
}

func (v GLIBCVersion) String() string { return v.Full }

type GLIBCVersions []GLIBCVersion

func (v GLIBCVersions) Len() int { return len(v) }
func (v GLIBCVersions) Less(i, j int) bool {
	if v[i].Major == v[j].Major {
		if v[i].Minor == v[j].Minor {
			return v[i].Patch < v[j].Patch
		}
		return v[i].Minor < v[j].Minor
	}
	return v[i].Major < v[j].Major
}
func (v GLIBCVersions) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v GLIBCVersions) String() string {
	var ss []string
	for _, version := range v {
		ss = append(ss, version.String())
	}
	return strings.Join(ss, " ")
}
func (v GLIBCVersions) FindViolations(expr string) (GLIBCVersions, error) {
	lang := gval.NewLanguage(gval.Base(), gval.Arithmetic(), gval.PropositionalLogic())
	var violations GLIBCVersions
	for _, version := range v {
		res, err := lang.Evaluate(expr, map[string]interface{}{
			"major": version.Major,
			"minor": version.Minor,
			"patch": version.Patch,
		})
		if err != nil {
			return GLIBCVersions{}, fmt.Errorf("evaluating '%s' on '%s': %w", expr, version, err)
		}
		b, ok := res.(bool)
		if !ok {
			return GLIBCVersions{}, fmt.Errorf("'%s' did not evaluate to a boolean", expr)
		}
		if !b {
			violations = append(violations, version)
		}
	}
	return violations, nil
}

func ParseGLIBCVersion(s string) (GLIBCVersion, error) {
	split := strings.Split(s, ".")
	// glibc is either <major>.<minor> or <major>.<minor>.<patch>
	if len(split) != 2 && len(split) != 3 {
		return GLIBCVersion{}, fmt.Errorf("unable to parse glibc version '%s'", s)
	}
	major, err := strconv.Atoi(split[0])
	if err != nil {
		return GLIBCVersion{}, err
	}
	minor, err := strconv.Atoi(split[1])
	if err != nil {
		return GLIBCVersion{}, err
	}

	v := GLIBCVersion{
		Full:  s,
		Major: major,
		Minor: minor,
	}

	if len(split) == 3 {
		patch, err := strconv.Atoi(split[2])
		if err != nil {
			return GLIBCVersion{}, err
		}
		v.Patch = patch
	}

	return v, nil
}

func ParseFile(f string) (GLIBCVersions, error) {
	elfFile, err := elf.Open(f)
	if err != nil {
		return GLIBCVersions{}, err
	}

	syms, err := elfFile.ImportedSymbols()
	if err != nil {
		return GLIBCVersions{}, err
	}
	versionStrs := map[string]bool{}
	for _, sym := range syms {
		versionStrs[sym.Version] = true
	}
	var versions GLIBCVersions
	for versionStr := range versionStrs {
		matches := r.FindStringSubmatch(versionStr)
		if len(matches) == 2 {
			version, err := ParseGLIBCVersion(matches[1])
			if err != nil {
				return GLIBCVersions{}, fmt.Errorf("parsing '%s' from '%s': %w", matches[1], f, err)
			}
			versions = append(versions, version)
		}
	}
	sort.Sort(versions)
	return versions, nil
}
