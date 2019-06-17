package gitlab

type GitLabProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	NameWithNamespace string `json:"name_with_namespace"`
	WebURL            string `json:"web_url"`
}
