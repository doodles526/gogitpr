package api

type CommitData struct {
	Label string   `json:"label"`
	Ref   string   `json:"ref"`
	Sha   string   `json:"sha"`
	User  UserData `json:"user"`
	Repo  RepoData `json:"repo"`
}
