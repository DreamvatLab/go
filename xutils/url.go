package xutils

import (
	"net/url"
	"path"
	"strings"

	"github.com/DreamvatLab/go/xerr"
)

// JointURL joins the basePath and additional paths, avoiding double slashes.
func JointURL(basePath string, paths ...string) (*url.URL, error) {
	// Parse the base URL
	r, err := url.Parse(basePath)
	if err != nil {
		return nil, xerr.WithStack(err)
	}

	// Join the paths using path.Join (it cleans double slashes)
	allPaths := append([]string{r.Path}, paths...)
	r.Path = path.Join(allPaths...)

	// Ensure path starts with '/' if basePath did
	if strings.HasPrefix(basePath, "/") && !strings.HasPrefix(r.Path, "/") {
		r.Path = "/" + r.Path
	}

	return r, nil
}

func JointURLString(basePath string, paths ...string) (string, error) {
	r, err := JointURL(basePath, paths...)
	if xerr.LogError(err) {
		return "", err
	}

	return r.String(), nil
}
