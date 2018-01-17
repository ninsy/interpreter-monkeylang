package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello world"}
	hello2 := &String{Value: "Hello world"}
	
	otherHello1 := &String{Value: "Hi, World!"}
	otherHello2 := &String{Value: "Hi, World!"}	

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("equal strings should have equal hashes!")
	}

	if otherHello1.HashKey() != otherHello2.HashKey() {
		t.Errorf("equal strings should have equal hashes!")
	}

	if hello1.HashKey() == otherHello1.HashKey() {
		t.Errorf("not equal strings should not have equal hashes!")		
	}

	if hello2.HashKey() == otherHello2.HashKey() {
		t.Errorf("not equal strings should not have equal hashes!")		
	}	
}

func TestIntHashKey(t *testing.T) {
	int1 := &Integer{Value: 1}
	int2 := &Integer{Value: 1}
	
	otherInt1 := &Integer{Value: 2}
	otherInt2 := &Integer{Value: 2}

	if int1.HashKey() != int2.HashKey() {
		t.Errorf("equal ints should have equal hashes!")
	}

	if otherInt1.HashKey() != otherInt2.HashKey() {
		t.Errorf("equal ints should have equal hashes!")
	}

	if int1.HashKey() == otherInt1.HashKey() {
		t.Errorf("not equal ints should not have equal hashes!")		
	}

	if int2.HashKey() == otherInt2.HashKey() {
		t.Errorf("not equal ints should not have equal hashes!")		
	}	
}

func TestBoolHashKey(t *testing.T) {
	bool1 := &Boolean{Value: true}
	bool2 := &Boolean{Value: true}
	
	otherBool1 := &Boolean{Value: false}
	otherBool2 := &Boolean{Value: false}

	if bool1.HashKey() != bool2.HashKey() {
		t.Errorf("equal bools should have equal hashes!")
	}

	if otherBool1.HashKey() != otherBool2.HashKey() {
		t.Errorf("equal bools should have equal hashes!")
	}

	if bool1.HashKey() == otherBool1.HashKey() {
		t.Errorf("not equal bools should not have equal hashes!")		
	}

	if bool2.HashKey() == otherBool2.HashKey() {
		t.Errorf("not equal bools should not have equal hashes!")		
	}	
}