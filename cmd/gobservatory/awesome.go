package main

import (
	"fmt"
	"sort"
	"strings"
)

func awesome(pc PonzuConnection, tagged bool) {
	// Get existing stars
	stars, err := GetFromPonzu(fmt.Sprintf("%s://%s:%s/api/contents?type=Star&count=-1", pc.Scheme, pc.Host, pc.Port))
	if err != nil {
		panic("Error getting stars from Ponzu")
	}
	sort.Sort(stars)
	var langs = make(map[string]interface{})
	var tg = make(map[string]interface{})
	// Build up the tags
	for i, star := range stars.Stars {
		if star.Language == "" {
			stars.Stars[i].Language = "Unknown"
			star.Language = stars.Stars[i].Language
		}
		if _, ok := langs[star.Language]; !ok {
			langs[star.Language] = nil
		}
		if len(star.Tags) == 0 {
			continue
		}
		for _, t := range star.Tags {
			if _, ok := tg[t]; !ok {
				tg[t] = nil // add tag into map
			}
		}
	}

	languages := []string{}
	for k := range langs {
		languages = append(languages, k)
	}
	sort.Strings(languages)

	tags := []string{}
	for k := range tg {
		tags = append(tags, k)
	}
	sort.Strings(tags)

	// Print contents
	fmt.Printf("\n## Contents\n")
	// languages
	for _, lang := range languages {
		fmt.Printf("- [%s](#%s)\n", lang, strings.Replace(strings.ToLower(lang), " ", "-", -1))
	}
	// tags
	if tagged {
		for _, tag := range tags {
			fmt.Printf("- [%s](#%s)\n", tag, strings.Replace(strings.ToLower(tag), " ", "-", -1))
		}
	}
	// Print by language
	for _, lang := range languages {
		fmt.Printf("\n## %s\n", lang)
		for _, star := range stars.Stars {
			if star.Language == lang {
				fmt.Printf("* [%s](%s) - %s\n", star.Name, star.HtmlUrl, star.Description)
			}
		}
	}
	if tagged {
		// Print by tag
		for _, tag := range tags {
			fmt.Printf("\n## %s\n", tag)
			for _, star := range stars.Stars {
				if len(star.Tags) == 0 {
					continue
				}
				if StarContainsTag(star, tag) {
					fmt.Printf("* [%s](%s) - %s\n", star.Name, star.HtmlUrl, star.Description)
				}
			}
		}
	}
}
