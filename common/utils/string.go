package utils

import (
	"regexp"
	"strings"

	common "github.com/parinithshekar/gitsink/common"
)

// FilterRepos filters the repositories according to regex and string patters
func FilterRepos(repos []common.Repository, include []string, exclude []string) []common.Repository {
	var finalRepos []common.Repository

	// RegEx patterns
	var includePatterns []string
	var excludePatterns []string

	// Plain repository names
	var includeNames []string
	var excludeNames []string

	for _, pattern := range include {
		isRE, _ := regexp.MatchString("^/.*/$", pattern)
		if isRE { // pattern is a regular expression -> /dock.*-pl.*/
			barePattern := strings.TrimSuffix(strings.TrimPrefix(pattern, "/"), "/")
			includePatterns = append(includePatterns, barePattern)
		} else { // pattern is a plain string -> repo-name-1
			includeNames = append(includeNames, pattern)
		}
	}

	for _, pattern := range exclude {
		isRE, _ := regexp.MatchString("^/.*/$", pattern)
		if isRE {
			barePattern := strings.TrimSuffix(strings.TrimPrefix(pattern, "/"), "/")
			excludePatterns = append(excludePatterns, barePattern)
		} else {
			excludeNames = append(excludeNames, pattern)
		}
	}

	for _, repo := range repos {
		includePass := false
		excludePass := true

		// Check if repository should be included
		for _, pattern := range includePatterns {
			match, _ := regexp.MatchString(pattern, repo.Slug)
			if match {
				includePass = true
				break
			}
		}
		if !includePass {
			for _, name := range includeNames {
				if name == repo.Slug {
					includePass = true
					break
				}
			}
		}

		// Check if repository should be excluded
		for _, pattern := range excludePatterns {
			match, _ := regexp.MatchString(pattern, repo.Slug)
			if match {
				excludePass = false
				break
			}
		}
		if excludePass {
			for _, name := range excludeNames {
				if name == repo.Slug {
					excludePass = false
					break
				}
			}
		}

		// Append to final repos if all filters pass
		if includePass && excludePass {
			finalRepos = append(finalRepos, repo)
		}
	}

	return finalRepos
}
