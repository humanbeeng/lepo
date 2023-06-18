package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

type GitCloner struct {
	URL        string
	TargetPath string
}

type GitClonerOpts struct {
	URL        string
	TargetPath string
}

func NewGitCloner(opts GitClonerOpts) *GitCloner {
	return &GitCloner{
		URL:        opts.URL,
		TargetPath: opts.TargetPath,
	}
}

func (gc *GitCloner) Clone() error {
	// Validate URL
	// Check if public repo
	// Send clone request
	r, err := git.PlainClone(gc.TargetPath, false, &git.CloneOptions{
		URL:               gc.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return nil
	}

	fmt.Println(r.Branches())
	return nil
}
