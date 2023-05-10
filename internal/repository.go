package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Release struct {
	ID          int    `json:"id"`
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
}

type Response struct {
	URL             string  `json:"url"`
	AssetsURL       string  `json:"assets_url"`
	UploadURL       string  `json:"upload_url"`
	HTMLURL         string  `json:"html_url"`
	ID              int     `json:"id"`
	Author          Author  `json:"author"`
	NodeID          string  `json:"node_id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Name            string  `json:"name"`
	Draft           bool    `json:"draft"`
	Prerelease      bool    `json:"prerelease"`
	CreatedAt       string  `json:"created_at"`
	PublishedAt     string  `json:"published_at"`
	Assets          []Asset `json:"assets"`
	TarballURL      string  `json:"tarball_url"`
	ZipballURL      string  `json:"zipball_url"`
	Body            string  `json:"body"`
}

type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Asset struct {
	URL                string      `json:"url"`
	ID                 int         `json:"id"`
	NodeID             string      `json:"node_id"`
	Name               string      `json:"name"`
	Label              interface{} `json:"label"`
	Uploader           Author      `json:"uploader"`
	ContentType        string      `json:"content_type"`
	State              string      `json:"state"`
	Size               int         `json:"size"`
	DownloadCount      int         `json:"download_count"`
	CreatedAt          string      `json:"created_at"`
	UpdatedAt          string      `json:"updated_at"`
	BrowserDownloadURL string      `json:"browser_download_url"`
}

type Repository struct {
	Name        string
	Owner       string
	ReleaseName string
	LocalDir    string
	Token       string
	Version     string
}

func (r *Repository) GetName() string {
	return r.Name
}

func (r *Repository) GetOwner() string {
	return r.Owner
}

func (r *Repository) GetContent() string {
	return fmt.Sprintf("%s.rb", r.GetName())
}

func (r *Repository) GetVersion() string {
	return r.Version
}

func (r *Repository) GetLocalDir() string {
	return r.LocalDir + "/" + r.GetName() + "/" + r.GetVersion() + "/bin/"
}

func (r *Repository) GetReleaseName() string {
	return fmt.Sprintf("homebrew-%s", r.GetName())
}

func (r *Repository) GetGitCloneURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s.rb",
		r.GetOwner(),
		r.GetReleaseName(),
		r.GetName())
}

func (r *Repository) GetGitReleaseLatestURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		r.GetOwner(),
		r.GetName())
}

func (r *Repository) GetDir() string {
	return r.GetLocalDir() + "/" + r.GetOwner() + "/" + r.GetReleaseName() + "/" + r.GetContent()
}

func (r *Repository) GetToken() string {
	return fmt.Sprintf("Bearer %s", r.Token)
}

func createFolder(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return nil
}

func CloneReleaseRepo(r *Repository) error {
	client := resty.New()

	if _, err := os.Stat(r.GetLocalDir()); errors.Is(err, os.ErrNotExist) {
		return err
	}
	createFolder(r.GetLocalDir() + "/" + r.GetOwner())
	createFolder(r.GetLocalDir() + "/" + r.GetOwner() + "/" + r.GetReleaseName())

	resp, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3.raw").
		SetHeader("Authorization", r.GetToken()).
		Get(r.GetGitCloneURL())
	if err != nil {
		return err
	}

	err = os.WriteFile(r.GetDir(), resp.Body(), os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}

func GetReleaseLatest(r *Repository, os *OS) (map[string]Asset, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", r.GetToken()).
		Get(r.GetGitReleaseLatestURL())
	if err != nil {
		return nil, err
	}

	result := &Response{}
	json.Unmarshal(resp.Body(), &result)

	// var table map[string]Asset
	table := make(map[string]Asset)

	for _, v := range result.Assets {
		if len(strings.Split(v.Name, "_")) > 3 {
			if os.GetName() == strings.Split(v.Name, "_")[2] && os.GetArch() == strings.Split(strings.Split(v.Name, ".")[2], "_")[2] {
				table[extractVersionFromURL(result.HTMLURL)] = v
			}
		}
	}

	return table, nil
}

func DownloadRelease(asset map[string]Asset, r *Repository) (string, error) {
	client := resty.New()

	path := r.GetLocalDir() + r.GetName() + ".tar.gz"
	_, err := client.R().
		SetHeader("Accept", "application/octet-stream").
		SetHeader("Authorization", r.GetToken()).
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetOutput(path).
		Get(asset[r.GetVersion()].URL)
	if err != nil {
		return "", err
	}

	return path, nil
}

func extractVersionFromURL(url string) string {
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
