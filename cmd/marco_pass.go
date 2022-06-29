package cmd

import (
	"bufio"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"nist.local/isoboroff/dedupe/lib"
)

// betterCmd represents the BETTER command
var marco_passCmd = &cobra.Command{
	Use:   "marco_pass [JSON lines file]",
	Short: "Dedupe the MSMARCO v2 passage collection",
	Long:  `Compute near-duplicate hashes for passages in the MSMARCO collection.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lshThresh := viper.GetFloat64("lsh.threshold")
		lshBuckets := viper.GetInt("lsh.buckets")
		minhashSize := viper.GetInt("minhash.size")
		lsh := lib.MakeLSHForThreshold(lshBuckets, lshThresh)
		minhash := lib.NewMinhash(minhashSize)

		dd := lib.MakeDeduper(lsh, *minhash, marco_read)
		dd.Dedupe(args[0])
	},
}

func marco_read(reader *bufio.Reader, c chan lib.Document) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		article := gjson.Parse(line)
		docid := article.Get("pid").String()
		text := article.Get("passage").String()

		var title string
		if len(text) < 50 {
			title = text[:]
		} else {
			title = text[:25]
		}
		title = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return ' '
			} else {
				return r
			}
		}, title)

		text = marco_preprocess(text)
		c <- lib.Document{text, docid, title}
	}
	close(c)
}

func marco_preprocess(text string) string {
	result := strings.TrimSpace(text)
	result = strings.ToLower(result)
	result = html_re.ReplaceAllLiteralString(result, " ")
	result = nonword_re.ReplaceAllLiteralString(result, " ")
	return result
}

func init() {
	nonword_re = regexp.MustCompile(`[^\pL]+`)
	html_re = regexp.MustCompile(`<[^>]+?>`)

	rootCmd.AddCommand(marco_passCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wapoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wapoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
