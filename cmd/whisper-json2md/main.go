package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const TimecodePeriod = 60

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: First argument must be the path to a whisper json file.")
		os.Exit(1)
	}

	jsonFilePath := os.Args[1]

	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var whisperData WhisperData
	json.Unmarshal(jsonContent, &whisperData)

	prevSegmentStart := float32(-1)
	for _, segment := range whisperData.Segments {
		if segment.Start > prevSegmentStart+TimecodePeriod {
			prevSegmentStart = segment.Start
			fmt.Printf(" [%v]", segment.timecode())
		}
		fmt.Printf(segment.Text)
	}
}

type WhisperData struct {
	Text     string           `json:"text"`
	Segments []WhisperSegment `json:"segments"`
}

type WhisperSegment struct {
	Id               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float32 `json:"start"`
	End              float32 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float32 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

func (s *WhisperSegment) naturalStart() (int, int, int) {
	hours := int(s.Start / 3600)
	remainder := int(s.Start) % 3600
	minutes := int(remainder / 60)
	seconds := remainder % 60
	return hours, minutes, seconds
}

func (s *WhisperSegment) timecode() string {
	hours, minutes, seconds := s.naturalStart()
	if hours > 0 {
		return fmt.Sprintf("%v:%02v:%02v", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("%v:%02v", minutes, seconds)
	}
}
