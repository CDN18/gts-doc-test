/*
   GoToSocial
   Copyright (C) 2021 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package util

import (
	"fmt"
	"strings"
)

// DeriveMentionsFromStatus takes a plaintext (ie., not html-formatted) status,
// and applies a regex to it to return a deduplicated list of accounts
// mentioned in that status.
//
// It will look for fully-qualified account names in the form "@user@example.org".
// or the form "@username" for local users.
// The case of the returned mentions will be lowered, for consistency.
func DeriveMentionsFromStatus(status string) []string {
	mentionedAccounts := []string{}
	for _, m := range mentionFinderRegex.FindAllStringSubmatch(status, -1) {
		mentionedAccounts = append(mentionedAccounts, m[1])
	}
	return lower(unique(mentionedAccounts))
}

// DeriveHashtagsFromStatus takes a plaintext (ie., not html-formatted) status,
// and applies a regex to it to return a deduplicated list of hashtags
// used in that status, without the leading #. The case of the returned
// tags will be lowered, for consistency.
func DeriveHashtagsFromStatus(status string) []string {
	tags := []string{}
	for _, m := range hashtagFinderRegex.FindAllStringSubmatch(status, -1) {
		tags = append(tags, m[1])
	}
	return lower(unique(tags))
}

// DeriveEmojisFromStatus takes a plaintext (ie., not html-formatted) status,
// and applies a regex to it to return a deduplicated list of emojis
// used in that status, without the surround ::. The case of the returned
// emojis will be lowered, for consistency.
func DeriveEmojisFromStatus(status string) []string {
	emojis := []string{}
	for _, m := range emojiFinderRegex.FindAllStringSubmatch(status, -1) {
		emojis = append(emojis, m[1])
	}
	return lower(unique(emojis))
}

// ExtractMentionParts extracts the username @test_user and the domain @example.org
// from a mention string like @test_user@example.org.
//
// If no domain is provided, it will return just the username part.
//
// If nothing is matched, it will return an error.
func ExtractMentionParts(mention string) (username, domain string, err error) {
	matches := mentionNameRegex.FindStringSubmatch(mention)
	if matches == nil {
		err = fmt.Errorf("could't match mention %s", mention)
		return
	}
	fmt.Println(matches)
	username = matches[1]
	if len(matches) == 2 {
		domain = matches[2]
	}
	return
}

// unique returns a deduplicated version of a given string slice.
func unique(s []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// lower lowercases all strings in a given string slice
func lower(s []string) []string {
	new := []string{}
	for _, i := range s {
		new = append(new, strings.ToLower(i))
	}
	return new
}

// HTMLFormat takes a plaintext formatted status string, and converts it into
// a nice HTML-formatted string.
//
// This includes:
// - Replacing line-breaks with <p>
// - Replacing URLs with hrefs.
// - Replacing mentions with links to that account's URL as stored in the database.
func HTMLFormat(status string) string {
	// TODO: write proper HTML formatting logic for a status
	return status
}
