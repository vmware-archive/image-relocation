package pathmapping

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"

	"github.com/pivotal/image-relocation/pkg/image"
)

type PathMapping func(repoPrefix string, originalImage image.Name) image.Name

func FlattenRepoPath(repoPrefix string, originalImage image.Name) image.Name {
	hasher := md5.New()
	hasher.Write([]byte(originalImage.Name()))
	hash := hex.EncodeToString(hasher.Sum(nil))
	available := reference.NameTotalLengthMax - len(mappedPath(repoPrefix, "", hash))
	fp := flatPath(originalImage.Path(), available)
	var mp string
	if fp == "" {
		mp = fmt.Sprintf("%s/%s", repoPrefix, hash)
	} else {
		mp = mappedPath(repoPrefix, fp, hash)
	}
	mn, err := image.NewName(mp)
	if err != nil {
		panic(err) // handle more gracefully
	}
	return mn
}

func mappedPath(repoPrefix string, repoPath string, hash string) string {
	return fmt.Sprintf("%s/%s-%s", repoPrefix, repoPath, hash)
}

func flatPath(repoPath string, size int) string {
	return strings.Join(crunch(strings.Split(repoPath, "/"), size), "-")
}

func crunch(components []string, size int) []string {
	for n := len(components); n > 0; n-- {
		comp := reduce(components, n)
		if len(strings.Join(comp, "-")) <= size {
			return comp
		}

	}
	return []string{}
}

func reduce(components []string, n int) []string {
	if len(components) < 2 || len(components) <= n {
		return components
	}

	tmp := make([]string, len(components))
	copy(tmp, components)

	last := components[len(tmp)-1]
	if n < 2 {
		return []string{last}
	}

	front := tmp[0 : n-1]
	return append(front, "-", last)
}
