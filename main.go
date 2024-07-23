package main

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Date  string   `yaml:"date"`
	Title string   `yaml:"title"`
	Slug  string   `yaml:"slug"`
	Tags  []string `yaml:"tags"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("not enough args %v", os.Args[0])
	}
	path := os.Args[1]

	fileBytes, err := os.ReadFile(path)
	checkErr(err)

	frontMatter, md, err := extractFrontMatter(fileBytes)
	checkErr(err)
	title := frontMatter.Title
	slug := frontMatter.Slug
	insert(title, slug, md)
}

func insert(title, slug, md string) {
	db, err := sql.Open("sqlite3", "../zalgorithm/db/notes.db")
	checkErr((err))
	defer db.Close()

	insertNoteQuery := `INSERT INTO notes (title, slug, markdown) VALUES (?, ?, ?)`
	stmt, err := db.Prepare(insertNoteQuery)
	checkErr(err)

	defer stmt.Close()

	_, err = stmt.Exec(title, slug, md)
	checkErr(err)
}

func extractFrontMatter(mdBytes []byte) (FrontMatter, string, error) {
	var frontMatter FrontMatter
	parts := bytes.SplitN(mdBytes, []byte("---\n"), 3)
	if len(parts) == 3 {
		err := yaml.Unmarshal(parts[1], &frontMatter)
		checkErr(err)
		return frontMatter, string(parts[2]), nil
	}

	// Return empty front matter if not present
	return frontMatter, string(mdBytes), errors.New("missing or invalid front matter")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
}
