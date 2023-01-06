package cmd

import (
       "strings"
       "sort"
)

type PullRequestTemplate struct {
     IssueLink		 string

     TypeOfChanges	 []string
     ChecklistQuestions	 []string

     Description	 string
     TypeOfChange	 int
     Checklist		 []int
}

func NewPullRequestTemplate(issueLink string) PullRequestTemplate {
     return PullRequestTemplate {
     	    IssueLink: issueLink,
     	    TypeOfChanges: []string{
				"Chore (no production code changed. eg: doc, style, tests)Chore (no production code changed. eg: doc, style, tests)",
				"Bug fix (non-breaking change which fixes an issue)",
				"New feature (non-breaking change which adds functionality)",
				"Breaking change (fix or feature that would cause existing functionality to not work as expected)",
				},
	    ChecklistQuestions: []string{
				     "I have commented on my code, particularly in hard-to-understand areas",
				     "I have made corresponding changes to the documentation",
				     "The required manual tests were executed and no issue was found",
				     },
     }
}

func (this PullRequestTemplate) Generate() string {
     var sb strings.Builder

     sb.WriteString("# Definition of Done\n\n")

     sb.WriteString("## Description\n\n")
     sb.WriteString(this.Description)
     sb.WriteString("\n\n")

     sb.WriteString("## Type of change\n\n")
     sb.WriteString(CheckedItem(this.TypeOfChanges[this.TypeOfChange]))
     sb.WriteString("\n\n")

     sb.WriteString(this.TaskHeader())
     sb.WriteString("\n\n")
     sb.WriteString("- ")
     sb.WriteString(this.IssueLink)
     sb.WriteString("\n\n")

     sb.WriteString("## Checklist:\n")
     for i, checklistQuestion := range this.ChecklistQuestions {
     	 sb.WriteString("\n")
	 sb.WriteString(Item(this.HasCheckedQuestion(i), checklistQuestion))
     }

     return sb.String()
}

func (this PullRequestTemplate) TaskHeader() string {
     if this.IsNewFeature() {
     	return "### Features"
     } else {
        return "### Fixes"
     }
}

func (this PullRequestTemplate) IsNewFeature() bool {
     return this.TypeOfChange == 2
}

func Item(checked bool, item string) string {
     if checked {
     	return CheckedItem(item)
     } else {
        return UncheckedItem(item)
     }
}

func CheckedItem(item string) string {
     return "[x] - " + item
}

func UncheckedItem(item string) string {
     return "[ ] - " + item
}

func (this PullRequestTemplate) HasCheckedQuestion(questionIndex int) bool {
     n := len(this.Checklist)
     
     indexFound := sort.Search(n, func(i int) bool {
        return this.Checklist[i] >= questionIndex
     })

     return indexFound < n && this.Checklist[indexFound] == questionIndex
}