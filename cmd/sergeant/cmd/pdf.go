package cmd

import (
	"bytes"
	"fmt"
	"time"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/albatross-org/sergeant"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// pdfCmd represents the 'pdf' command.
var pdfCmd = &cobra.Command{
	Use:   "pdf --set-name [set name] -n [number of questions]",
	Short: "Generate a PDF containing questions",
	Long: `pdf lets you generate a PDF containing questions from a specific set.
	
For example:

	$ sergeant pdf --set "school" --user "ben" --num 10 --view "bayesian" | pdflatex -jobname "PDF Questions $(date +%Y-%m-%d)"
	# or, using the shorter version of flags
	$ sergeant pdf -s "differentiation" -u "olly" -n 10 -v "random" | pdflatex -jobname "PDF Questions $(date +%Y-%m-%d)"
	`,

	Run: func(cmd *cobra.Command, args []string) {
		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}

		config, err := sergeant.LoadConfig(configPath)
		if err != nil {
			logrus.Fatal(err)
		}

		underlyingStore, err := albatross.FromConfig(config.Store)
		if err != nil {
			logrus.Fatal(err)
		}

		setName, err := cmd.Flags().GetString("set")
		if err != nil {
			logrus.Fatal(err)
		}

		n, err := cmd.Flags().GetInt("num")
		if err != nil {
			logrus.Fatal(err)
		}

		user, err := cmd.Flags().GetString("user")
		if err != nil {
			logrus.Fatal(err)
		}

		viewName, err := cmd.Flags().GetString("view")
		if err != nil {
			logrus.Fatal(err)
		}

		store := sergeant.NewStore(underlyingStore, config)

		set, warnings, err := store.Set(setName)
		if err != nil {
			logrus.Fatal(err)
		}

		for path, warning := range warnings {
			logrus.Warningf("Malformed card: %s -> %s", path, warning)
		}

		completed := 0
		for _, card := range set.Cards {
			if card.TotalCompletions() > 0 {
				completed++
			}
		}

		logrus.Infof("Loaded %d cards from %q, %d completed. (%.2f%%)", len(set.Cards), setName, completed, 100*float64(completed)/float64(len(set.Cards)))

		questions := []*sergeant.Card{}

		for i := 0; i < n; i++ {
			question := sergeant.DefaultViews[viewName].Next(set, user)
			questions = append(questions, question)
		}

		source := generatePDFSource(questions, user, viewName)
		fmt.Println(source)
	},
}

func generatePDFSource(questions []*sergeant.Card, user, viewName string) string {
	var doc bytes.Buffer

	doc.WriteString(fmt.Sprintf(`\documentclass[a4paper]{report}
\usepackage[utf8]{inputenc}
\usepackage{graphicx}
\usepackage{url}
\pagenumbering{gobble}

\title{\huge Sergeant Questions}
\author{For %s (%s)}
\date{%s}
\setlength{\parindent}{0pt}
\begin{document}
\maketitle
\tableofcontents
\newpage`, user, "``"+viewName+"\"", time.Now().Format("January 2nd, 2006")))

	for i, question := range questions {
		doc.WriteString(fmt.Sprintf(`\section*{Question (%d)}
\noindent\makebox[\linewidth]{\rule{\paperwidth}{0.4pt}}

\smallskip
\huge Start:

\medskip
\large ID: \texttt{%s}

Path: \path{%s}

\noindent\makebox[\linewidth]{\rule{\paperwidth}{0.4pt}}

\medskip

\includegraphics[width=\textwidth,height=\textheight,keepaspectratio]{%s}

\newpage\phantom{A}
\newpage\phantom{A}
\newpage\phantom{A}

\section*{Answer (%d)}
\noindent\makebox[\linewidth]{\rule{\paperwidth}{0.4pt}}

\smallskip
\huge End:

\medskip
\large ID: \texttt{%s}

Path: \path{%s}

\begin{itemize}
    \item Perfect
    \item Major
    \item Minor
\end{itemize}

\noindent\makebox[\linewidth]{\rule{\paperwidth}{0.4pt}}

\medskip

\includegraphics[width=\textwidth,height=\textheight,keepaspectratio]{%s}

\newpage`, i+1, question.ID, question.Path, question.QuestionPath, i+1, question.ID, question.Path, question.AnswerPath))
	}

	doc.WriteString(`\end{document}`)

	return doc.String()
}

func init() {
	pdfCmd.Flags().StringP("set", "s", "all", "set to get questions from")
	pdfCmd.Flags().StringP("user", "u", "", "user to get questions for")
	pdfCmd.Flags().StringP("view", "v", "bayesian", "view to get questions from (options: random, unseen, bayesian)")
	pdfCmd.Flags().IntP("num", "n", 10, "number of questions to grab")

	rootCmd.AddCommand(pdfCmd)
}
