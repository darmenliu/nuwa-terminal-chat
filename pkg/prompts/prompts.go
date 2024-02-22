package prompts

import "strings"

const (
	// prompt for generate code format
	FileFormatPrompt string = "You will output the content of each file necessary to achieve the goal, including ALL code.\n" +
		"Represent files like so:\n\n" +

		"@FILENAME@\n" +
		"```\n" +
		"CODE\n" +
		"```\n\n" +

		"The following tokens must be replaced like so:\n" +
		"FILENAME is the lowercase combined path and file name including the file extension\n" +
		"CODE is the code in the file\n\n" +

		"Example representation of a file:\n\n" +

		"@cmd/hello_world.go@\n" +
		"```\n" +
		"package main\n\n" +
		"import \"fmt\"\n\n" +
		"func main() {\n" +
		"    fmt.Println(\"Hello, World!\")\n" +
		"}\n" +
		"```\n\n" +

		"Do not comment on what every file does. Please note that the code should be fully functional. No placeholders."

	CodeGeneratorPrompt string = `Think step by step and reason yourself to the correct decisions to make sure we get it right.
First lay out the names of the core classes, functions, methods that will be necessary, As well as a quick comment on their purpose.

FILE_FORMAT

You will start with the "entrypoint" file, then go to the ones that are imported by that file, and so on.
Please note that the code should be fully functional. No placeholders.

Follow Golang and framework appropriate best practice file naming convention.
Make sure that files contain all imports, types etc.  The code should be fully functional. Make sure that code in different files are compatible with each other.
Ensure to implement all code, if you are unsure, write a plausible implementation.
Include module dependency or package manager dependency definition file.
Before you finish, double check that all parts of the architecture is present in the files.

When you are done, write finish with "this concludes a fully working implementation".`

	PhilosophyPrompt string = `Almost always put different classes in different files.
Always use Golang as the programming language.
Always add a comment briefly describing the purpose of the function definition.
Add comments explaining very complex bits of logic.
Always follow the best practices for the Golang for folder/file structure and how to package the project.


Python toolbelt preferences:
- pytest
- dataclasses`

	RoadmapPrompt string = `You will get instructions for code to write.
You will write a very long answer. Make sure that every detail of the architecture is, in the end, implemented as code.`
)

func GetCodeGeneratorPrompt(fileFormat string) string {
	return strings.Replace(CodeGeneratorPrompt, "FILE_FORMAT", fileFormat, 1)
}

func GetSysPrompt() string {
	return RoadmapPrompt + "\n\n" + PhilosophyPrompt + "\n\n" + GetCodeGeneratorPrompt(FileFormatPrompt)
}

func GetUserPrompt(userPrompt string) string {
	return GetSysPrompt() + "\n\n" + userPrompt
}
