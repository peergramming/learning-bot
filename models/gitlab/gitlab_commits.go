package gitlab

type GitlabRepoCommits struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}
