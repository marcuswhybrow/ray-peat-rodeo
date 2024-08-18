package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

const TimecodePeriod = 60

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("First argument must be the path to a whisper json file.")
	}

	jsonFilePath := os.Args[1]

	timestampOffset := 0

	if len(os.Args) >= 3 {
		seconds, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Second argument must be the number of seconds by which to offset timestamps in seconds.")
		}
		timestampOffset = seconds
	}

	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var whisperData WhisperData
	json.Unmarshal(jsonContent, &whisperData)

	prevSegmentStart := float32(-TimecodePeriod * 2)
	for _, segment := range whisperData.Segments {
		if segment.Start > prevSegmentStart+TimecodePeriod {
			prevSegmentStart = segment.Start
			fmt.Printf(" [%v]", segment.timestamp(timestampOffset))
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

func (s *WhisperSegment) naturalStart(offset int) (int, int, int) {
	rawSeconds := s.Start + float32(offset)
	hours := int(rawSeconds / 3600)
	remainder := int(s.Start) % 3600
	minutes := int(remainder / 60)
	seconds := remainder % 60
	return hours, minutes, seconds
}

func (s *WhisperSegment) timestamp(offset int) string {
	hours, minutes, seconds := s.naturalStart(offset)
	if hours > 0 {
		return fmt.Sprintf("%v:%02v:%02v", hours, minutes, seconds)
	} else {
		return fmt.Sprintf("%v:%02v", minutes, seconds)
	}
}
