package bump

import (
	"fmt"
	"regexp"

	"github.com/coreos/go-semver/semver"
)

const (
	semverMatcher = `(\d+\.){1}(\d+\.){1}(\*|\d+)`
)

type (
	Options struct {
		Replace    string
		Part       string
		Index      int
		PreRelease string
		Metadata   string
	}

	Bump struct {
		Old, New string
		Loc      []int
		NewBytes []byte
	}
)

func NewBump(v string, part string) *Bump {
	s := semver.New(v)
	switch part {
	case "major":
		s.BumpMajor()
	case "minor":
		s.BumpMinor()
	default:
		s.BumpPatch()
	}

	return &Bump{
		Old: v,
		New: s.String(),
	}
}

func (b *Bump) String() string {
	return b.New
}

func (b *Bump) Replace(v string) {
	b.Old = v
	b.New = v
}

func (b *Bump) Bump(part string) {
	s := semver.New(b.Old)
	switch part {
	case "major":
		s.BumpMajor()
	case "minor":
		s.BumpMinor()
	default:
		s.BumpPatch()
	}
	b.New = s.String()
}

//func (b *Bump) ReplaceInContent(vbytes []byte) (newcontents []byte, err error) {
//	return replace(vbytes, b.New, "", 0)
//}
//
//func (b *Bump) BumpInContent(vbytes []byte) (newcontents []byte, err error) {
//	return replace(vbytes, "", "", 0)
//}

func (b *Bump) Tag() error {
	return nil
}

// BumpInContent takes finds the first semver string in the content, bumps it, then returns the same content with the new version
func BumpInContent(vbytes []byte, part string, index int) (old, new string, loc []int, newcontents []byte, err error) {
	return replace(vbytes, "", part, index)
}

// ReplaceInContent takes finds the first semver string in the content and replaces it with replaceWith
func ReplaceInContent(vbytes []byte, replaceWith string, index int) (old, new string, loc []int, newcontents []byte, err error) {
	return replace(vbytes, replaceWith, "", index)
}

func ReplaceInContent2(vbytes []byte, options *Options) (old, new string, loc []int, newcontents []byte, err error) {
	return replace2(vbytes, options)
}

func BumpInContent2(vbytes []byte, options *Options) (old, new string, loc []int, newcontents []byte, err error) {
	return replace2(vbytes, options)
}

func BumpString(input string, options *Options) (string, error) {
	oldC, newC, loc, n2, err := replace2([]byte(input), options)
	if err != nil {
		return "", err
	}
	fmt.Println(oldC, newC, loc, n2)
	return newC, nil
}

// if index is set, it will find all matches and choose the one at the given index, -1 means last
func replace(vbytes []byte, replace, part string, index int) (old, new string, loc []int, newcontents []byte, err error) {
	options := &Options{
		Replace: replace,
		Part:    part,
		Index:   index,
	}
	return replace2(vbytes, options)
}

func replace2(vbytes []byte, options *Options) (old, new string, loc []int, newcontents []byte, err error) {
	re := regexp.MustCompile(semverMatcher)
	if options.Index == 0 {
		loc = re.FindIndex(vbytes)
	} else {
		locs := re.FindAllIndex(vbytes, -1)
		if locs == nil {
			return "", "", nil, nil, fmt.Errorf("did not find semantic version")
		}

		locsLen := len(locs)
		if options.Index >= locsLen {
			return "", "", nil, nil, fmt.Errorf("semver index to replace out of range. Found %v, want %v", locsLen, options.Index)
		}

		if options.Index < 0 {
			loc = locs[locsLen+options.Index]
		} else {
			loc = locs[options.Index]
		}
	}
	// fmt.Println(loc)
	if loc == nil {
		return "", "", nil, nil, fmt.Errorf("Did not find semantic version")
	}

	vs := string(vbytes[loc[0]:loc[1]])

	if options.Replace == "" {
		// fmt.Println("bumping", vs, "part", options.Part)
		v := semver.New(vs)
		switch options.Part {
		case "major":
			v.BumpMajor()
		case "minor":
			v.BumpMinor()
		default:
			v.BumpPatch()
		}

		if options.PreRelease != "" {
			v.PreRelease = semver.PreRelease(options.PreRelease)
		}

		if options.Metadata != "" {
			v.Metadata = options.Metadata
		}

		options.Replace = v.String()
	}

	len1 := loc[1] - loc[0]
	additionalBytes := len(options.Replace) - len1
	// Create and fill an extended buffer
	b := make([]byte, len(vbytes)+additionalBytes)
	copy(b[:loc[0]], vbytes[:loc[0]])
	copy(b[loc[0]:loc[1]+additionalBytes], options.Replace)
	copy(b[loc[1]+additionalBytes:], vbytes[loc[1]:])
	// fmt.Printf("writing: '%v'", string(b))

	return vs, options.Replace, loc, b, nil
}
