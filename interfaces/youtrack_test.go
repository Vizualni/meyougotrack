package interfaces

import (
	"testing"
)

func TestExtractIssueIdWithOnlyUrl(t *testing.T) {
	validateExtractIssueId(
		"http://youtrack.trikoder.net/issue/asa-123",
		"asa-123",
		false,
		t,
	)
}

func TestExtractIssueIdWithOnlyHttpsUrl(t *testing.T) {
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/asa-123",
		"asa-123",
		false,
		t,
	)
}

func TestExtractIssueWithMultipleIssueIdTypes(t *testing.T) {
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/3K1AP-123",
		"3K1AP-123",
		false,
		t,
	)
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/3K1a-123",
		"3K1a-123",
		false,
		t,
	)
}

func TestExtractIssueWithInvalidUrl(t *testing.T) {
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/3K1AP-123whatisthis",
		"something wrong",
		true,
		t,
	)
}

func TestExtractIssueWithExtraTextAfter(t *testing.T) {
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/3K1AP-123 Oh look something",
		"3K1AP-123",
		false,
		t,
	)
}

func TestExtractIssueWithExtraTextBefore(t *testing.T) {
	validateExtractIssueId(
		"Something before https://youtrack.trikoder.net/issue/3K1AP-123",
		"3K1AP-123",
		false,
		t,
	)
}

func TestExtractIssueWithExtraTextBeforeAndAfter(t *testing.T) {
	validateExtractIssueId(
		"Something before https://youtrack.trikoder.net/issue/3K1AP-123 Something after",
		"3K1AP-123",
		false,
		t,
	)
}

func TestExtractIssueWithMultipleUrlsFindsOnlyFirstOne(t *testing.T) {
	validateExtractIssueId(
		"https://youtrack.trikoder.net/issue/3K1AP-123 hello https://youtrack.trikoder.net/issue/3K1AP-456",
		"3K1AP-123",
		false,
		t,
	)
}

func TestExtractWithoutUrl(t *testing.T) {
	validateExtractIssueId(
		"Hello",
		"not found",
		true,
		t,
	)
}

func validateExtractIssueId(input, output string, expectsError bool, t *testing.T) {

	s := SimpleRegexIssueIdExtractor{}

	id, err := s.Extract(input)

	if err != nil {
		if expectsError == false {
			t.Fatalf("Unexpected error: %s", err)
			return
		}
	}

	if id == output {
		if expectsError == true {
			// well...unexpected error
			t.Fatal("expected error here")
		}
		// nice
		return
	}

	if expectsError == false {
		// didnt expect that
		t.Fatal("found id and expected id not matching")
	}

}
