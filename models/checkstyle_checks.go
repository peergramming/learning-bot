package models

var Checks = map[string]CheckDesc{
	"ArrayTrailingComma": CheckDesc{Category: "coding",
		Description: `A comma should be inserted at the end of the last element of the array if there are no left or right curly bracket on the same line.`,
		Rationale:   `Putting a comma at the end of each element allows you to easily change the order of the array, or add new elements at the end without changing the surrounding lines.`,
		Example: `int[] numbers = {
  1,
  2  // Violation: Missing comma.
};`},
	"EmptyStatement": CheckDesc{Category: "coding",
		Description: "Code should not contain empty statements.",
		Rationale:   "Empty statements may introduce bugs and can be hard to spot.",
		Example: `if (someCondition);
  doConditional(); // This will always run no matter the value of 'someCondition'
doUnconditional();`},
	"EqualsHashCode": CheckDesc{Category: "coding",
		Description: "Any class which overriders either equals() or hashcode() must override the other.",
		Rationale:   "Both equals() and hashcode() should depend on the same set of fields, so you can use your class in hash-based collections"},
	"IllegalCatch": CheckDesc{Category: "coding",
		Description: "Catch statements should not handle exception types like 'Exception', 'RuntimeException', or 'Throwable'.",
		Rationale:   "It is never acceptable to catch these types of exception superclasses, as these may lead to catching unexpected errors such as NullPointerException or OutOfMemoryException"},
}

type CheckDesc struct {
	Category    string
	Description string
	Rationale   string
	Example     string
}
