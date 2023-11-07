package meta

type FrontMatter struct {
	Source struct {
		Title    string
		Series   string
		Url      string
		Duration string
	}
	Speakers      map[string]string
	Transcription struct {
		Source string
		Date   string
		Author string
	}
}
