package frontmatter_test

import (
	"testing"

	"github.com/lorenzogood/blog/internal/frontmatter"
	"github.com/stretchr/testify/assert"
)

type SampleFrontmatter struct {
	Title       string `toml:"title"`
	Description string `toml:"description"`
	Author      string `toml:"author"`
	Date        string `toml:"date"`
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        string
		expectedData SampleFrontmatter
		expectedRest string
		expectedErr  error
	}{
		{
			name: "valid frontmatter",
			input: `+++
title = "Test Document"
description = "This is a test document."
author = "John Doe"
date = "2024-11-17"
+++
This is the body of the document.`,
			expectedData: SampleFrontmatter{
				Title:       "Test Document",
				Description: "This is a test document.",
				Author:      "John Doe",
				Date:        "2024-11-17",
			},
			expectedRest: "\nThis is the body of the document.",
			expectedErr:  nil,
		},
		{
			name:         "missing delimiters",
			input:        "\nThis is just a body without frontmatter.",
			expectedData: SampleFrontmatter{},
			expectedRest: "",
			expectedErr:  frontmatter.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsedData SampleFrontmatter
			rest, err := frontmatter.Parse([]byte(test.input), &parsedData)

			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedData, parsedData)
			assert.Equal(t, test.expectedRest, string(rest))
		})
	}
}
