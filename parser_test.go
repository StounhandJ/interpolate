package interpolate_test

import (
	"reflect"
	"testing"

	"github.com/mfridman/interpolate"
)

func TestParser(t *testing.T) {
	var testCases = []struct {
		String   string
		Expected []interpolate.ExpressionItem
	}{
		{
			String: `Buildkite... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}`,
			Expected: []interpolate.ExpressionItem{
				{Text: "Buildkite... "},
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
					},
				},
				{Text: " "},
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "ANOTHER_VAR",
						Expressions: []interpolate.Expansion{
							interpolate.EmptyValueExpansion{
								Identifier: "ANOTHER_VAR",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "🏖"},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `${TEST1:- ${TEST2:-$TEST3}}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "TEST1",
						Expressions: []interpolate.Expansion{
							interpolate.EmptyValueExpansion{
								Identifier: "TEST1",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: " "},
									{
										Expansion: interpolate.VariableExpansion{
											Identifier: "TEST2",
											Expressions: []interpolate.Expansion{
												interpolate.EmptyValueExpansion{
													Identifier: "TEST2",
													Content: interpolate.Expression([]interpolate.ExpressionItem{
														{
															Expansion: interpolate.VariableExpansion{
																Identifier: "TEST3",
															},
														},
													}),
												}},
										},
									},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.UnsetValueExpansion{
								Identifier: "HELLO_WORLD",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "blah"},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `\\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Text: `\\`},
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.UnsetValueExpansion{
								Identifier: "HELLO_WORLD",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "blah"},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Text: `$`},
				{Text: `{HELLO_WORLD-blah}`},
			},
		},
		{
			String: `Test \\\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Text: `Test `},
				{Text: `\\`},
				{Text: `$`},
				{Text: `{HELLO_WORLD-blah}`},
			},
		},
		{
			String: `${HELLO_WORLD:1}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.SubstringExpansion{
								Identifier: "HELLO_WORLD",
								Offset:     1,
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD: -1}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.SubstringExpansion{
								Identifier: "HELLO_WORLD",
								Offset:     -1,
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD:-1}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.EmptyValueExpansion{
								Identifier: "HELLO_WORLD",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "1"},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD:1:7}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.SubstringExpansion{
								Identifier: "HELLO_WORLD",
								Offset:     1,
								Length:     7,
								HasLength:  true,
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD:1:-7}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.SubstringExpansion{
								Identifier: "HELLO_WORLD",
								Offset:     1,
								Length:     -7,
								HasLength:  true,
							},
						},
					},
				},
			},
		},
		{
			String: `${HELLO_WORLD?Required}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "HELLO_WORLD",
						Expressions: []interpolate.Expansion{
							interpolate.RequiredExpansion{
								Identifier: "HELLO_WORLD",
								Message: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "Required"},
								}),
							},
						},
					},
				},
			},
		},
		{
			String: `$`,
			Expected: []interpolate.ExpressionItem{
				{Text: `$`},
			},
		},
		{
			String: `\`,
			Expected: []interpolate.ExpressionItem{
				{Text: `\`},
			},
		},
		{
			String: `$(echo hello world)`,
			Expected: []interpolate.ExpressionItem{
				{Text: `$(`},
				{Text: `echo hello world)`},
			},
		},
		{
			String: `${DATA:1:3-Tuesday}`,
			Expected: []interpolate.ExpressionItem{
				{
					Expansion: interpolate.VariableExpansion{
						Identifier: "DATA",
						Expressions: []interpolate.Expansion{
							interpolate.SubstringExpansion{
								Identifier: "DATA",
								Offset:     1,
								Length:     3,
								HasLength:  true,
							},
							interpolate.UnsetValueExpansion{
								Identifier: "DATA",
								Content: interpolate.Expression([]interpolate.ExpressionItem{
									{Text: "Tuesday"},
								}),
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.String, func(t *testing.T) {
			actual, err := interpolate.NewParser(tc.String).Parse()
			if err != nil {
				t.Fatal(err)
			}

			expected := interpolate.Expression(tc.Expected)
			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Expected vs Actual: \n%s\n\n%s", expected, actual)
			}
		})
	}
}
