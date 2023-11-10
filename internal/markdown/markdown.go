package markdown

type FrontMatter struct {
	Source struct {
		Series   string
		Title    string
		Url      string
		Kind     string
		Duration string
	}
	Speakers      map[string]string
	Transcription struct {
		Url    string
		Kind   string
		Date   string
		Author string
	}
}
