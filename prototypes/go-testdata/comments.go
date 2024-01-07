package main

// Single line comment on interface declaration
type TestInterface interface {
	InlineCommentedFunc() // This is an inline comment

	// This is a single line comment
	SingleLineCommentedFunc(string) error

	// First line
	// Second line
	MultiLineCommentedFunc(string) error

	/*
		This is a block comment
	*/

	BlockCommentedFunc(string) error
}
