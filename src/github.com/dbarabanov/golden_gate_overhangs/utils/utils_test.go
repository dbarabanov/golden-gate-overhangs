package utils
import (
    "testing"
    "fmt"
    )

func TestEncode1(t *testing.T) {
    encodings := map[string] byte {"A": 0, "C": 1, "G": 2, "T": 3}
    for key, value := range encodings {
        actual := Encode1(key)
        if  value != actual {
            t.Errorf("Encode1(\"%v\") = %v, want %v", key, actual, value)
        }
    }
}

func TestEncode1Panic(t *testing.T) {
    letter := "x"
    defer func() {
        expected := fmt.Sprintf("Invalid argument to Encode1(): %v. Must be one of A,C,T,G.", letter)
        actual := recover()
        if  expected != actual {
            t.Errorf("Encode1(\"%v\") panic msg:\n \"%v\"\n    want: \n\"%v\"", letter, actual, expected)
        }

        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    Encode1(letter)
}
