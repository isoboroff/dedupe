package cmd

import (
	"bufio"
	"strings"
	"regexp"
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
		// TODO: make sure num-hashes and jaccard-thresh are CLI options!!
		lsh := lib.MakeLSHForThreshold(256, 0.9)
		minhash := lib.NewMinhash(256)

		dd := lib.MakeDeduper(lsh, *minhash, read)
		dd.Dedupe(args[0])
	},
}

func read(reader *bufio.Reader, c chan lib.Document) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
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
		text := getWapoText(article)
		text = preprocess(text)
		c <- lib.Document{text, docid, title}
	}
	close(c)
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
