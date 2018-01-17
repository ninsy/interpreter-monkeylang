package evaluator

import (
	"math"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5", 15},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"3 * ( 3 * 3 ) + 10", 37},
		{"3 + 3 * 3 / 3", 6},
		{"-50 + 200 - 100", 50},
		{"30 * -1 + 15 * 2", 0},
		{"50 / 2 * 3 - 5", 70},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, want=%d",
			obj, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 < 2", true},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(2 < 1) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestMinusPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"-10", -10},
		{"-5", -5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testMinusPrefixExpression(t, evaluated, tt.expected)
	}
}

func testMinusPrefixExpression(t *testing.T, obj object.Object, expected interface{}) bool {
	switch which := expected.(type) {
	case int64:
		return testIntegerObject(t, obj, int64(which))
	default:
		return false
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL, got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 *5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if(10 > 1) {
				if(10 > 1) {
					return 10;
				}
				return 1;
			}
			`,
			10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`{"k":"v"}[fn(x) { return x; }];`,
			"Can't hash object of type FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T (%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters, should have 1, got=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("body is not '(x + 2)', got=%q", fn.Body.String())
	}
}

func TestFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let ident = fn(x) { return x; }; ident(5);", 5},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5,5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5+5, add(5,5));", 20},
		{"fn(x) { return x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string, got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has wrong value: %s", str.Value)
	}
}

func TestStringConcatOperator(t *testing.T) {
	input := `"Hello" + ", " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("evaluated object is not string, got=%q (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, World!" {
		t.Errorf("String has wrong value, got=%q", str.Value)
	}
}

func TestBuiltinParserMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"parseFloat('abc')", math.NaN},
		{"parseFloat('3.14abc')", 3.14},
		{"parseFloat('a3.14')", math.NaN},
		{"parseFloat(NaN)", math.NaN},
		{"parseFloat('-4.12')", 4.12},
		{"parseInt(' 0xF')", 16},
		{"parseInt('1111')", 1111},
		{"parseInt('1111', 2)", 15},
		{"parseInt('15 * 3', 10)", 15},
		{"parseInt('Hello')", math.NaN},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		// add support for float, nans
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not an array, got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("wrong number of elements, should have 3, got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
	
}

func TestArrayIndexExprs(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{
			"[1,2,3][0]",
			'a',
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"let i = 3; [1,2,3,4,5][i]",
			4,
		},
		{
			"[1,2,3][5]",
			nil,
		},
		{
			"[1,2,3][-5]",
			nil,
		},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, eval, int64(integer))
		} else {
			testNullObject(t, eval)
		}
	}
}

func TestObjectLiterals(t *testing.T) {
	input := `
		{
			"one": 2 - 1,
			"two": 4 - 2,
			3: 3,
			false: 4,
		}
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.ObjectLiteral)
	if !ok {
		t.Fatalf("eval has not returned object.ObjectLiteral, but %T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64 {
		(&object.String{Value: "one"}).HashKey(): 1,
		(&object.String{Value: "two"}).HashKey(): 2,
		(&object.Integer{Value: 3}).HashKey(): 3,
		FALSE.HashKey(): 4,		
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Object literal has different amount of pairs than expected, got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("Cannot find matching value for given key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input		 string
		expected interface{}
	}{
		{
			`{"key": 1}["key"]`,
			1,
		},
		{
			`{"key": 1}["null"]`,
			nil,
		},
		{
			`let k = "key"; {"key": 3}[k]`,
			3,
		},
		{
			`{}["null"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 6}[true]`,
			6,
		},
		{
			`{false: 7}[!true]`,
			7,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}