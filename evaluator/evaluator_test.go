package evaluator

import (
	"math"
	"testing"

	"github.com/dudewhocode/sushi/lexer"
	"github.com/dudewhocode/sushi/object"
	"github.com/dudewhocode/sushi/parser"
	"github.com/google/go-cmp/cmp"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"1.3454", 1.3454},
		{"-5.1235", -5.1235},
		{"-10.1112", -10.1112},
		{"5.54 + 5.22 + 5.562 + 5.323 - 10.124", 11.521000},
		{"2.25 * 2.25 * 2.25 * 2.25 * 2.25", 57.6650390625},
		{"-50.05 + 100.10 + -50.05", 0},
		{"5.5 * 2.5 + 10.5", 24.25},
		{"5.25 + 2.25 * 10.25", 28.3125},
		{"20.10 + 2 * -10.05", 0},
		{"50.5 / 2.5 * 2.5 + 10", 60.5},
		{"2.75 * (5.99 + 10.99)", 46.695},
		{"3.1 * 3.1 * 3.1 + 10.1", 39.891000},
		{"3.25 * (3.555 * 3.124) + 10.00", 46.093915},
		{"(5.01 + 10.02 * 2.03 + 15.04 / 3.05) * 2.09 + -10.08", 53.208852},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"!!5.12", true},
		{"1.1 < 2.5", true},
		{"1.125 > 2", false},
		{"1.11 < 1.11", false},
		{"1.00 > 1", false},
		{"1.2 == 1.2", true},
		{"1.5 != 1.5", false},
		{"1.123 == 2.123", false},
		{"1.12 != 2", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 < 2) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
		{"if (1.12 < 2) {10.0001} else {20}", 10.0001},
		{"if (1.12 > 2) {10.0001} else {20}", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case float64:
			testFloatObject(t, evaluated, float64(v))
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`if (10 > 1) {
				 if (10 > 1) {
				   return 10;
			}
			return 1; }
			`, 10},
		{"return 3.14", 3.14},
		{"return 25; 25.5", 25},
		{"9; return 2 * 5.5; 9;", 11.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case float64:
			testFloatObject(t, evaluated, float64(v))
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"3.14 + true; 3.14;",
			"type mismatch: FLOAT + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
				}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name":"Monkey"}[fn(x) { x }];`,
			"unhashable key: FUNCTION",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let a = 5; a;", 5},
		{"let a = 3.14; a;", 3.14},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5.5 * 5; a;", 27.5},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 25.5; let b = a; b;", 25.5},
		{"let a = 5.5; let b = a; let c = a + b + 5.5; c;", 16.5},
	}

	for _, tt := range tests {
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, testEval(tt.input), int64(v))
		case float64:
			testFloatObject(t, testEval(tt.input), float64(v))
		default:
			testNullObject(t, testEval(tt.input))
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2;}"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("functoin has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { x; }; identity(3.14);", 3.14},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(3.14);", 3.14},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let double = fn(x) { x * 2; }; double(5.5);", 11.0},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(1.1, 1.1);", 2.2},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"let add = fn(x, y) { x + y; }; add(5.5 + 5.5, add(5.5, 5.5));", 22.0},
		{"fn(x) {x; }(5)", 5},
	}

	for _, tt := range tests {
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, testEval(tt.input), int64(v))
		case float64:
			testFloatObject(t, testEval(tt.input), float64(v))
		default:
			testNullObject(t, testEval(tt.input))
		}
	}
}

func TestClosures(t *testing.T) {
	input := `
   let newAdder = fn(x) {
     fn(y) { x + y };
};
   let addTwo = newAdder(2);
   addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got =%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len(3.14)`, "argument to `len` not supported, got FLOAT"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([1, 2, 3])`, 3},
		{`len([1, 2, 3.14])`, 3},
		{`len([])`, 0},
		{`first([1, 2, 3])`, 1},
		{`first([1.0001, 2, 3])`, 1.0001},
		{`first([])`, nil},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`first(1.0)`, "argument to `first` must be ARRAY, got FLOAT"},
		{`last([1, 2, 3])`, 3},
		{`last([1, 2, 3.14])`, 3.14},
		{`last([])`, nil},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`last(1.1)`, "argument to `last` must be ARRAY, got FLOAT"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([1.24, 2.25, 3.14])`, []float64{2.25, 3.14}},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push([], 1.4)`, []float64{1.4}},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
		{`push(1.1, 1)`, "argument to `push` must be ARRAY, got FLOAT"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case float64:
			testFloatObject(t, evaluated, float64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		case []float64:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testFloatObject(t, array.Elements[i], float64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestFloatArrayLiterals(t *testing.T) {
	input := "[1.1, 2 * 2.25, 3.14 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testFloatObject(t, result.Elements[0], 1.1)
	testFloatObject(t, result.Elements[1], 4.5)
	testFloatObject(t, result.Elements[2], 6.14)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"[1.2, 2.2, 3.3][0]",
			1.2,
		},
		{
			"[1.1, 2.2, 3.3][1]",
			2.2,
		},
		{
			"[1.1, 2.2, 3.3][2]",
			3.3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"let i = 0; [1.223][i];",
			1.223,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"[1.1, 2.25, 3.14][1 + 1];",
			3.14,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1.1, 2.2, 3.3]; myArray[2];",
			3.3,
		},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{"let myArray = [1.1, 2.2, 3.3]; myArray[0] + myArray[1] + myArray[2];",
			6.6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"let myArray = [1.1, 2.2, 3.3]; let i = myArray[0]; myArray[i]",
			"ERROR: index operator not supported: ARRAY",
		},
		{"[1, 2, 3][3]",
			nil},
		{"[1.1, 2.2, 3.3][3]",
			nil},
		{
			"[1, 2, 3][-1]",
			nil,
		},
		{
			"[1.1, 2.2, 3.3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case float64:
			testFloatObject(t, evaluated, float64(v))
		case string:
			if v != evaluated.Inspect() {
				t.Errorf("expected error: %s, got %s", evaluated.Inspect(), v)
			}
		default:
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"floating": 10 + 11.25,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6,
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]interface{}{
		(&object.String{Value: "one"}).HashKey():      1,
		(&object.String{Value: "two"}).HashKey():      2,
		(&object.String{Value: "floating"}).HashKey(): 21.25,
		(&object.String{Value: "three"}).HashKey():    3,
		(&object.Integer{Value: 4}).HashKey():         4,
		TRUE.HashKey():                                5,
		FALSE.HashKey():                               6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedVal := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		switch v := expectedVal.(type) {
		case int:
			testIntegerObject(t, pair.Value, int64(v))
		case float64:
			testFloatObject(t, pair.Value, float64(v))
		}

	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5.5}["foo"]`,
			5.5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`let key = "foo"; {"foo": 5.5}[key]`,
			5.5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5},
		{
			`{5: 5.5}[5]`,
			5.5},
		{
			`{true: 5}[true]`,
			5,
		}, {
			`{false: 5}[false]`,
			5,
		},
		{
			`{true: 5.5}[true]`,
			5.5,
		}, {
			`{false: 5.5}[false]`,
			5.5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch v := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(v))
		case float64:
			testFloatObject(t, evaluated, float64(v))
		default:
			testNullObject(t, evaluated)
		}
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
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	const tolerance = .00001
	opt := cmp.Comparer(func(x, y float64) bool {
		diff := math.Abs(x - y)
		mean := math.Abs(x+y) / 2.0
		if math.IsNaN(diff / mean) {
			return true
		}
		return (diff / mean) < tolerance
	})

	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}

	if !cmp.Equal(result.Value, expected, opt) {
		t.Errorf("object has wrong value. got=%g, want=%g", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
