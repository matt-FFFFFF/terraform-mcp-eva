package prompt

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddSolveAvmIssuePrompt(s *mcp.Server) {
	s.AddPrompt(&mcp.Prompt{
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "issue_number",
				Description: "The issue number to be processed, for example: `1234`. If not provided, the prompt will try to infer it from the context or ask the user to provide it.",
			},
			{
				Name:        "category",
				Description: "The category of the issue, for example: `bug`, `feature`, `doc`, 'chore', etc. If not provided, the prompt will try to infer it from the context or ask the user to provide it.",
			},
		},
		Description: "If you're processing changes to repo that contains `terraform` and `avm` in repo's name, and there is `avm` or `avm.bat` file in the root directory, you should use this prompt to get instructions on how to process the changes. The prompt will return a list of instructions that you can follow to process the changes.",
		Name:        "solve_avm_issue",
	}, func(ctx context.Context, session *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
		issueNumber := params.Arguments["issue_number"]
		category := params.Arguments["category"]
		return &mcp.GetPromptResult{
			Meta:        mcp.Meta{},
			Description: "",
			Messages: []*mcp.PromptMessage{
				{
					Content: &mcp.TextContent{
						Annotations: &mcp.Annotations{
							Audience: []mcp.Role{
								"assistant",
							},
							Priority: 1.0,
						},
						Text: fmt.Sprintf(`As an Azure Verified Modules development expert, you must strictly follow these steps:
Analyze the user's request: and extract the issue number from it.
The issue number is %s, and the category is %s.
Use git checkout -b <category>/<issue-number> to create and switch to a new branch.
Create a new file named 'todo.md' in the root directory of the repository, write down your analysis of the issue, and provide a detailed plan on how to resolve it. Make sure to refer to this plan to track progress.
If you want to create or update Terraform blocks, you must consult the mcp server to get the latest Terraform schema and provider information first. When you want to query the schema and document, try tools have 'query_' prefix first. If you want to query azapi provider's resource schema, try tools have 'query_azapi_' prefix first.'
Remember to update the 'todo.md' file with the progress you made.
Only create new examples if there is a significant new feature being added to the module. Don't forget to add '_footer.md' and '_header.md' files like other examples.'
[IMPORTANT REFERENCE] The Azure Verified Modules specification index is available at this location: https://azure.github.io/Azure-Verified-Modules/llms.txt. Download this file and retrieve any relevant information from it. The specification references starting with TF* are pertinent. Any BC* should be ignored.
[CRITICAL STEP] After all changes are complete, you must execute in a bash shell:
1. ./avm pre-commit

If checks succeeds too then you should:

1. commit the changes with proper commit message, DO NOT commit 'todo.md' file, you can remove this.
2. propose creating a Pull Request (PR). If it fails, report the failure message, try to solve the issues with best effort.

Now, please begin execution.`, issueNumber, category,
						)},
					Role: "user",
				},
			},
		}, nil
	})
}
