package git

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// GithubRepository defines a github (enterprise) repo
type GithubRepository struct {
	Name              string `mapstructure:"github-repo-name"`
	ServerURL 		  string `mapstructure:"github-enterprise-server-url"`
	Organization      string `mapstructure:"github-organization"`
	GithubAccessToken string `mapstructure:"github-auth-token"`
	License           string `mapstructure:"license"`
	ProjectDir 		  string `mapstructure:"project-dir"`
	Private 		  bool   `mapstructure:"private"`
	CloneURL 		  string
}

// option is a generic functional configuration option for a GithubRepository
type option func(*GithubRepository)

// New creates a new GithubRepository configured with the provided options
func New(name string, options ...option) (*GithubRepository, error) {
	r := GithubRepository{
		Name:              name,
		License:           "mit",
		ProjectDir:        ".",
	}
	for _, opt := range options {
		opt(&r)
	}
	return &r, nil
}

// Configure allows us to modify the configuration of a GithubRepository
func (r *GithubRepository) Configure(opts ...option) {
	for _, opt := range opts {
		opt(r)
	}
}

func Name(n string) option {
	return func(r *GithubRepository) {
		r.Name = n
	}
}

func ServerURL(u string) option {
	return func(r *GithubRepository) {
		r.ServerURL = u
	}
}

func Organization(o string) option {
	return func(r *GithubRepository) {
		r.Organization = o
	}
}

func GithubAccessToken(t string) option {
	return func(r *GithubRepository) {
		r.GithubAccessToken = t
	}
}

func License(l string) option {
	return func(r *GithubRepository) {
		r.License = l
	}
}

func ProjectDir(d string) option {
	return func(r *GithubRepository) {
		r.ProjectDir = d
	}
}

func CloneURL(u string) option {
	return func(r *GithubRepository) {
		r.CloneURL = u
	}
}

// Create uses a GithubRepository to construct a Github (Enterprise) repo via API calls
func (r GithubRepository) Create() error {
	var private bool
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: r.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	var client *github.Client
	if r.ServerURL != "" {
		client, _ = github.NewEnterpriseClient(r.ServerURL, r.ServerURL, tc)
		private = true
	} else {
		client = github.NewClient(tc)
	}
	repo := &github.Repository{
		Name:    github.String(r.Name),
		Private: github.Bool(private),
		GitignoreTemplate: github.String("Go"),
		DefaultBranch: github.String("develop"),
	}
	if r.License != "" {
		repo.LicenseTemplate = github.String(strings.ToLower(r.License))
	}
	repository, _, err := client.Repositories.Create(ctx, r.Organization, repo);if err != nil {
		return errors.Wrap(err, "failed to create remote/origin")
	}
	r.CloneURL = *repository.CloneURL
	return nil
}

// Clone performs a plain `git clone` of "url" as the project dir.  If the output arg is supplied, progress is logged there.
func (r GithubRepository) Clone(output io.Writer) error {
	_, err := git.PlainClone(r.ProjectDir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "scaffold", // yes, this can be anything except an empty string
			Password: r.GithubAccessToken,
		},
		URL:      r.CloneURL,
		Progress: output,
	})
	return err
}

func (r GithubRepository) URL() (string, error) {
	if r.CloneURL == "" {
		return "", errors.Errorf("%s is empty", "CloneURL")
	}
	return r.CloneURL, nil
}