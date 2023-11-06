package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Building Ray Peat Rodeo...")

	if err := os.MkdirAll("build", os.ModePerm); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	f, err := os.Create("build/index.html")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}

	err = hello("Marcus").Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}

	fmt.Println("Done.")
}
