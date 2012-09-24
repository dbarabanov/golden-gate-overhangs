package utils
import (
    "testing"
    "fmt"
    )

func TestEncodeBase(t *testing.T) {
    encodings := map[byte] byte {'A': 0, 'C': 1, 'G': 2, 'T': 3}
    for key, value := range encodings {
        actual := EncodeBase(key)
        if  value != actual {
            t.Errorf("EncodeBase('%c') = %v, want %v", key, actual, value)
        }
    }
}

func TestEncodeBasePanic(t *testing.T) {
    var letter byte = 'x' 
    defer func() {
        expected := fmt.Sprintf("Invalid argument to EncodeBase(): '%c'. Must be one of A,C,T,G.", letter)
        actual := recover()
        if  expected != actual {
            t.Errorf("EncodeBase('%c') panic msg:\n \"%v\"\n    want: \n\"%v\"", letter, actual, expected)
        }
    }()
    EncodeBase(letter)
}

func TestDecodeBase(t *testing.T) {
    encodings := map[string] byte {"A": 0, "C": 1, "G": 2, "T": 3}
    for key, value := range encodings {
        actual := DecodeBase(value)
        if  key != actual {
            t.Errorf("DecodeBase(\"%v\") = %v, want %v.", value, actual, key)
        }
    }
}

func TestDecodeBasePanic(t *testing.T) {
    var encoded_letter byte = 10
    defer func() {
        expected := fmt.Sprintf("Invalid argument to DecodeBase(): %v. Must be one of 0,1,2,3.", encoded_letter)
        actual := recover()
        if  expected != actual {
            t.Errorf("DecodeBase(\"%v\") panic msg:\n \"%v\"\n    want: \n\"%v\".", encoded_letter, actual, expected)
        }
    }()
    DecodeBase(encoded_letter)
}

func TestEncodeOverhang(t *testing.T) {
    encodings := map[string] byte {"AAAA": 0, "AAAC": 1, "TTTT": 255}
    for key, value := range encodings {
        actual := EncodeOverhang(key)
        if  value != actual {
            t.Errorf("EncodeOverhang(\"%v\") = %v, want %v.", key, actual, value)
        }
    }
}

func TestEncodeOverhangPanic(t *testing.T) {
    overhang := "AATAA" //5 letters
    defer func() {
        expected := fmt.Sprintf("Invalid argument to EncodeOverhang(): %v. Must be a 4-letter string of A,C,T,G.", overhang)
        actual := recover()
        if  expected != actual {
            t.Errorf("EncodeOverhang(\"%v\") panic msg:\n \"%v\"\n    want: \n\"%v\".", overhang, actual, expected)
        }
    }()
    EncodeOverhang(overhang)
}

func TestDecodeOverhang(t *testing.T) {
    decodings := map[byte] string{0: "AAAA", 1: "AAAC", 255: "TTTT"}
    for key, value := range decodings {
        actual := DecodeOverhang(key)
        if  value != actual {
            t.Errorf("DecodeOverhang(\"%v\") = %v, want %v.", key, actual, value)
        }
    }
}

func TestComplementaryOverhang(t *testing.T) {
    overhang := "ACCG"
    expected := "TGGC"
    actual := ComplementaryOverhang(overhang)
    if  expected != actual {
        t.Errorf("ComplementaryOverhang(\"%v\") = %v, want %v.", overhang, actual, expected)
    }
}

func TestReverseOverhang(t *testing.T) {
    overhang := "ACCG"
    expected := "GCCA"
    actual := ReverseOverhang(overhang)
    if  expected != actual {
        t.Errorf("ComplementaryOverhang(\"%v\") = %v, want %v.", overhang, actual, expected)
    }
}

func TestPartnerOverhang(t *testing.T) {
    overhang := "ACCG"
    expected := "CGGT"
    actual := PartnerOverhang(overhang)
    if  expected != actual {
        t.Errorf("PartnerOverhang(\"%v\") = %v, want %v.", overhang, actual, expected)
    }
}

func TestGetOverhangRepeatCount(t *testing.T) {
    overhang1 := "AAAA"
    overhang2 := "AAAT"
    var expected byte = 3
    actual := GetOverhangRepeatCount(overhang1, overhang2)
    if  expected != actual {
        t.Errorf("GetOverhangRepeatCount(\"%v\", \"%v\") = %v, want %v.", overhang1, overhang2, actual, expected)
    }
}

