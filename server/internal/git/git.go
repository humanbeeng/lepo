package git

import (
	"log"

	"github.com/go-git/go-git/v5"
)

type GitCloner struct{}

type GitCloneRequest struct {
	URL        string
	TargetPath string
}

func NewGitCloner() *GitCloner {
	return &GitCloner{}
}

func (gc *GitCloner) Clone(req GitCloneRequest) error {
	// Validate URL
	// Check if public repo
	// Send clone request
	log.Println("Clone requested for", req.URL)
	_, err := git.PlainClone(req.TargetPath, false, &git.CloneOptions{
		URL:               req.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}
	log.Println("Clone completed for", req.URL, "location", req.TargetPath)
	return nil
}
