/*
Copyright © 2019 Ian Soboroff <ian.soboroff@nist.gov>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"os"
	"log"
	"strings"
	"regexp"
	"fmt"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"nist.local/isoboroff/dedupe/lib"
)

// wapoCmd represents the wapo command
var wapoCmd = &cobra.Command{
	Use:   "wapo",
	Short: "Dedupe the Washington Post collection",
	Long: `Compute near-duplicate hashes for documents in the Washington Post
collection.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		log.Println("--- First pass, indexing documents")
		
		reader := bufio.NewReader(file)
		// TODO: make sure num-hashes and jaccard-thresh are CLI options!!
		lsh := lib.MakeLSHForThreshold(256, 0.9)
		minhash := lib.NewMinhash(256)
		titles := make(map[string]string)
		linecount := 0
		
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			linecount++
			if (linecount % 10000) == 0 {
				log.Println(linecount, "docs")
			}
			article := gjson.Parse(line)
			docid := article.Get("id").String()
			title := article.Get("title").String()
			title = strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return ' '
				} else {
					return r
				}
			}, title)
			titles[docid] = title
			text := getWapoText(article)
			text = preprocess(text)

			sigs := lib.Shingle(text)
			sigs = minhash.Hash(sigs)
			lsh.Insert(docid, sigs)
		}

		log.Println("--- Second pass, identifying duplicates")
		
		file.Seek(0, 0)
		reader = bufio.NewReader(file)
		id2cluster := make(map[string]string, linecount)
		linecount = 0

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			linecount++
			if (linecount % 10000) == 0 {
				log.Println(linecount, "docs")
			}
			article := gjson.Parse(line)
			docid := article.Get("id").String()
			cluster, ok := id2cluster[docid]
			if ok {
				fmt.Println(cluster, docid, titles[docid])
				continue
			}

			id2cluster[docid] = docid
			fmt.Println(docid, docid, titles[docid])

			text := getWapoText(article)
			text = preprocess(text)

			sigs := lib.Shingle(text)
			sigs = minhash.Hash(sigs)
			dupes := lsh.Query(sigs)
			for _, d := range dupes {
				if _, ok := id2cluster[d]; !ok {
					id2cluster[d] = docid
				}
			}
		}

	},
}

func getWapoText(obj gjson.Result) string {
	var textbuf strings.Builder
	obj.Get("contents").ForEach(func(key, val gjson.Result) bool {
		if strings.HasPrefix(val.Get("mime").String(), "text/") {
			textbuf.WriteString(val.Get("content").String())
			textbuf.WriteRune(' ')
		}
		return true
	})
	return textbuf.String()
}

var nonword_re *regexp.Regexp
var html_re *regexp.Regexp

func preprocess(text string) string {
	result := strings.TrimSpace(text)
	result = strings.ToLower(result)
	result = html_re.ReplaceAllLiteralString(result, " ")
	result = nonword_re.ReplaceAllLiteralString(result, " ")
	return result
}

func init() {
	nonword_re = regexp.MustCompile(`[^\pL]+`)
	html_re = regexp.MustCompile(`<[^>]+?>`)
	
	rootCmd.AddCommand(wapoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wapoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wapoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
