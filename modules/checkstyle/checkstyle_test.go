package checkstyle

import (
	"testing"
)

func TestParseLineIssue(t *testing.T) {
	ok, issue := parseLineIssue(`[WARN] /tmp/path/to/file/ModificationsTest.java:23: Line is longer than 100 characters (found 122). [LineLength]`)
	if !ok {
		t.Errorf("Not OK!")
	}
	if issue.CheckName != "LineLength" {
		t.Errorf("CheckName = \"%s\"; want \"%s\"\n", issue.CheckName, "LineLength")
	}
	if issue.FilePath != "/tmp/path/to/file/ModificationsTest.java" {
		t.Errorf("FilePath = \"%s\"; want \"%s\"\n", issue.FilePath, "/tmp/path/to/file/ModificationsTest.java")
	}
	if issue.LineNumber != 23 {
		t.Errorf("LineNumber = \"%d\"; want \"%d\"\n", issue.LineNumber, 23)
	}
	if issue.ColumnNumber != 0 {
		t.Errorf("ColumnNumber = \"%d\"; want \"%d\"\n", issue.ColumnNumber, 0)
	}
	if issue.Description != "Line is longer than 100 characters (found 122)." {
		t.Errorf("Description = \"%s\"; want \"%s\"\n", issue.Description, "Line is longer than 100 characters (found 122).")
	}
}
