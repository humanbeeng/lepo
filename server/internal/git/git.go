package git

import (
	"github.com/go-git/go-git/v5"
	"go.uber.org/zap"
)

type GitCloner struct {
	logger *zap.Logger
}

type GitCloneRequest struct {
	URL        string
	TargetPath string
}

func NewGitCloner(logger *zap.Logger) *GitCloner {
	return &GitCloner{
		logger: logger,
	}
}

func (gc *GitCloner) Clone(req GitCloneRequest) error {
	// Validate URL
	// Check if public repo
	// Send clone request
	gc.logger.Info("Clone requested", zap.String("url", req.URL))
	_, err := git.PlainClone(req.TargetPath, false, &git.CloneOptions{
		URL:               req.URL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return err
	}
	gc.logger.Info(
		"Clone completed",
		zap.String("url", req.URL),
		zap.String("location", req.TargetPath),
	)
	return nil
}
