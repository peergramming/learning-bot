package gitlab

import "testing"

func TestIssues(t *testing.T) {
	t.Logf("%+v", GetRepoIssues("learning-bot/test-proj"))

}
