package content

import (
	"fmt"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
	"net/http"
)

type Star struct {
	item.Item

	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	GithubId        int      `json:"github_id"`
	HtmlUrl         string   `json:"html_url"`
	Description     string   `json:"description"`
	Private         bool     `json:"private"`
	Fork            bool     `json:"fork"`
	Language        string   `json:"language"`
	OwnerLogin      string   `json:"owner_login"`
	OwnerAvatarUrl  string   `json:"owner_avatar_url"`
	OwnerUrl        string   `json:"owner_url"`
	OwnerId         int      `json:"owner_id"`
	OwnerType       string   `json:"owner_type"`
	Homepage        string   `json:"homepage"`
	Forks           int      `json:"forks"`
	Size            int      `json:"size"`
	StargazersCount int      `json:"stargazers_count"`
	DefaultBranch   string   `json:"default_branch"`
	StarredAt       string   `json:"starred_at"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	PushedAt        string   `json:"pushed_at"`
	Comments        string   `json:"comments"`
	Tags            []string `json:"tags"`
}

func (s *Star) Accept(res http.ResponseWriter, req *http.Request) error {
	return nil
}
func (s *Star) Approve(res http.ResponseWriter, req *http.Request) error {
	return nil
}
func (s *Star) AutoApprove(http.ResponseWriter, *http.Request) error {
	return nil
}

// MarshalEditor writes a buffer of html to edit a Star within the CMS
// and implements editor.Editable
func (s *Star) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(s,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Star field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Name", s, map[string]string{
				"label":       "Name",
				"type":        "text",
				"placeholder": "Enter the Name here",
			}),
		},
		editor.Field{
			View: editor.Input("FullName", s, map[string]string{
				"label":       "FullName",
				"type":        "text",
				"placeholder": "Enter the FullName here",
			}),
		},
		editor.Field{
			View: editor.Input("GithubId", s, map[string]string{
				"label":       "GithubId",
				"type":        "text",
				"placeholder": "Enter the GithubId here",
			}),
		},
		editor.Field{
			View: editor.Input("HtmlUrl", s, map[string]string{
				"label":       "HtmlUrl",
				"type":        "text",
				"placeholder": "Enter the HtmlUrl here",
			}),
		},
		editor.Field{
			View: editor.Tags("Tags", s, map[string]string{
				"label":       "Tags",
				"type":        "text",
				"placeholder": "Enter the Tags here",
			}),
		},
		editor.Field{
			View: editor.Input("Description", s, map[string]string{
				"label":       "Description",
				"type":        "text",
				"placeholder": "Enter the Description here",
			}),
		},
		editor.Field{
			View: editor.Richtext("Comments", s, map[string]string{
				"label":       "Comments",
				"type":        "text",
				"placeholder": "Comments",
			}),
		},
		editor.Field{
			View: editor.Input("Private", s, map[string]string{
				"label":       "Private",
				"type":        "text",
				"placeholder": "Enter the Private here",
			}),
		},
		editor.Field{
			View: editor.Input("Fork", s, map[string]string{
				"label":       "Fork",
				"type":        "text",
				"placeholder": "Enter the Fork here",
			}),
		},
		editor.Field{
			View: editor.Input("Language", s, map[string]string{
				"label":       "Language",
				"type":        "text",
				"placeholder": "Enter the Language here",
			}),
		},
		editor.Field{
			View: editor.Input("OwnerLogin", s, map[string]string{
				"label":       "OwnerLogin",
				"type":        "text",
				"placeholder": "Enter the OwnerLogin here",
			}),
		},
		editor.Field{
			View: editor.Input("OwnerAvatarUrl", s, map[string]string{
				"label":       "OwnerAvatarUrl",
				"type":        "text",
				"placeholder": "Enter the OwnerAvatarUrl here",
			}),
		},
		editor.Field{
			View: editor.Input("OwnerUrl", s, map[string]string{
				"label":       "OwnerUrl",
				"type":        "text",
				"placeholder": "Enter the OwnerUrl here",
			}),
		},
		editor.Field{
			View: editor.Input("OwnerId", s, map[string]string{
				"label":       "OwnerId",
				"type":        "text",
				"placeholder": "Enter the OwnerId here",
			}),
		},
		editor.Field{
			View: editor.Input("OwnerType", s, map[string]string{
				"label":       "OwnerType",
				"type":        "text",
				"placeholder": "Enter the OwnerType here",
			}),
		},
		editor.Field{
			View: editor.Input("Homepage", s, map[string]string{
				"label":       "Homepage",
				"type":        "text",
				"placeholder": "Enter the Homepage here",
			}),
		},
		editor.Field{
			View: editor.Input("Forks", s, map[string]string{
				"label":       "Forks",
				"type":        "text",
				"placeholder": "Enter the Forks here",
			}),
		},
		editor.Field{
			View: editor.Input("Size", s, map[string]string{
				"label":       "Size",
				"type":        "text",
				"placeholder": "Enter the Size here",
			}),
		},
		editor.Field{
			View: editor.Input("StargazersCount", s, map[string]string{
				"label":       "StargazersCount",
				"type":        "text",
				"placeholder": "Enter the StargazersCount here",
			}),
		},
		editor.Field{
			View: editor.Input("DefaultBranch", s, map[string]string{
				"label":       "DefaultBranch",
				"type":        "text",
				"placeholder": "Enter the DefaultBranch here",
			}),
		},
		editor.Field{
			View: editor.Input("StarredAt", s, map[string]string{
				"label":       "StarredAt",
				"type":        "text",
				"placeholder": "Enter the StarredAt here",
			}),
		},
		editor.Field{
			View: editor.Input("CreatedAt", s, map[string]string{
				"label":       "CreatedAt",
				"type":        "text",
				"placeholder": "Enter the CreatedAt here",
			}),
		},
		editor.Field{
			View: editor.Input("UpdatedAt", s, map[string]string{
				"label":       "UpdatedAt",
				"type":        "text",
				"placeholder": "Enter the UpdatedAt here",
			}),
		},
		editor.Field{
			View: editor.Input("PushedAt", s, map[string]string{
				"label":       "PushedAt",
				"type":        "text",
				"placeholder": "Enter the PushedAt here",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Star editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Star"] = func() interface{} { return new(Star) }
}
