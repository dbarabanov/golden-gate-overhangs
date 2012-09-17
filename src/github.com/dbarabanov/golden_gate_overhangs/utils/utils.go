package utils

import (
    "fmt"
    "math/rand"
    "io/ioutil"
    "strings"
    "time"
    )

func Encode1(letter string) (encoded byte) {
//    letter := strings.ToUpper(base)
    switch {
        case letter == "A":
            return 0
        case letter == "C":
            return 1
        case letter == "G":
            return 2
        case letter == "T":
            return 3
    }
    panic(fmt.Sprintf("Invalid argument to Encode1(): %v. Must be one of A,C,T,G.", letter))
}

func EncodeBase(letter byte) (encoded byte) {
    switch {
        case letter == 'A':
            return 0
        case letter == 'C':
            return 1
        case letter == 'G':
            return 2
        case letter == 'T':
            return 3
    }
    panic(fmt.Sprintf("Invalid argument to Encode1(): %v. Must be one of A,C,T,G.", letter))
}

func Decode1(encoded byte) (letter string) {
     switch {
        case encoded == 0:
            return "A"
        case encoded == 1: 
            return "C"
        case encoded == 2:
            return "G"
        case encoded == 3:
            return "T"
    }
    panic(fmt.Sprintf("Invalid argument to Decode1(): %v. Must be one of 0,1,2,3.", letter))
}

func Encode4(seq string) (encoded byte) {
    if len(seq) != 4 {
        panic(fmt.Sprintf("Invalid argument to Encode4(): %v. Must be a 4-letter string of A,C,T,G.", seq))
    }
//#    first := Encode1(seq[0])
    first := Encode1(string(seq[0]))
    second := Encode1(string(seq[1]))
    third := Encode1(string(seq[2]))
    forth := Encode1(string(seq[3]))
    return (forth)|(third<<2)|(second<<4)|(first<<6)
}

func EncodeOverhang1(seq string) (encoded byte) {
    if len(seq) != 4 {
        panic(fmt.Sprintf("Invalid argument to Encode4(): %v. Must be a 4-letter string of A,C,T,G.", seq))
    }
    first := EncodeBase(seq[0])
    second := EncodeBase(seq[1])
    third := EncodeBase(seq[2])
    forth := EncodeBase(seq[3])
    return (forth)|(third<<2)|(second<<4)|(first<<6)
}

func EncodeOverhang(seq string) (encoded byte) {
    if len(seq) != 4 {
        panic(fmt.Sprintf("Invalid argument to EncodeOverhang(): \"%v\". Must be a 4-letter string of A,C,T,G.", seq))
    }
    return (EncodeBase(seq[0])<<6)|(EncodeBase(seq[1])<<4)|(EncodeBase(seq[2])<<2)|(EncodeBase(seq[3]))
}

func Decode4(b byte) (s string) {
    return Decode1(b>>6) + Decode1(b<<2>>6) + Decode1(b<<4>>6) + Decode1(b<<6>>6)
}

func DecodeOverhang(b byte) (s string) {
    return Decode4(b)
}

func AreOverhangsCompatible(overhang1 string, overhang2 string, max_repeats byte) bool {
    return AreCompatible(EncodeOverhang(overhang1), EncodeOverhang(overhang2), max_repeats)
}

func AreCompatible(b1 byte, b2 byte, max_repeats byte) bool {
    return GetRepeatCount(b1, b2) <= max_repeats 
}

func ComplementaryOverhang(overhang string) string {
    return DecodeOverhang(Complementary(EncodeOverhang(overhang)))
}

func Complementary(b byte) byte {
    return 255-b 
}

func ReverseOverhang(overhang string) string {
    return DecodeOverhang(Reverse(EncodeOverhang(overhang)))
}

func Reverse(b byte) byte {
    return b&3<<6 + b&12<<2 + b&48>>2 + b&192>>6
}

func IsOverhangSelfCompatible(overhang string, max_repeats byte) bool{
    return IsSelfCompatible(EncodeOverhang(overhang), max_repeats)
}

func IsSelfCompatible(b byte, max_repeats byte) bool {
    return AreCompatible(b, Partner(b), max_repeats)
}

func PartnerOverhang(overhang string) string {
    return DecodeOverhang(Partner(EncodeOverhang(overhang)))
}

func Partner(b byte) byte {
    return Reverse(Complementary(b)) 
}

func IsOverhangCompatibleWithMany(overhang string, overhangs []string, max_repeats byte) bool {
    return IsCompatibleWithMany(EncodeOverhang(overhang), EncodeOverhangs(overhangs), max_repeats)
}

func EncodeOverhangs(overhangs []string) []byte {
    encoded_overhangs := make ([]byte, len(overhangs))
    for i := range encoded_overhangs {
        encoded_overhangs[i] = EncodeOverhang(overhangs[i])
    }
    return encoded_overhangs
}

func IsCompatibleWithMany(overhang byte, overhangs []byte, max_repeats byte) bool {
    for i := range overhangs {
        if !AreCompatible(overhang, overhangs[i], max_repeats) {
            return false
        }
    }
    return true
}

func GetRepeatCount(b1 byte, b2 byte) byte {
    return 4-(
        (((b1^b2)&1)|((b1^b2)&2>>1)) + 
        (((b1^b2)&4>>2)|((b1^b2)&8>>3)) + 
        (((b1^b2)&16>>4)|((b1^b2)&32>>5)) + 
        (((b1^b2)&64>>6)|((b1^b2)&128>>7)))
}

func GenerateRandomGrid(cuts int, levels int) [][]byte {
    rand.Seed(time.Now().Unix())
//    fmt.Printf("Seeding random number disabled. %v\n", time.Now())
//    rand.Seed(3)
    grid := make([][]byte, levels, levels)
    for i := range grid {
        grid[i] = make([]byte, cuts, cuts)
        for j := range grid[i]{
            grid[i][j] = byte(rand.Int())
        }
    }
    return grid
}

func GridFromFile(filename string) (overhangs []string, grid[][]byte) {
    fmt.Printf("Reading file \""+filename+"\".\n")
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        panic("unable to read file \"" + filename+"\".\n")
    }
    lines := strings.Split(string(content), "\n")

    if len(lines[len(lines)-1]) == 0 {
        lines = lines[:len(lines)-1]
    }

    grid = make([][]byte, len(lines), len(lines))

    for i, line := range lines {
    //    fmt.Printf("%v line: %v\n", i, line)
        if len(line) != 8 {
            panic("Every line must have exactly 8 letters. I.e. \"ATTCGTGT\". Found: " + line)
        }
        grid[i] = make([]byte, 5, 5)
        for j := 0; j<5; j++ {
            overhang := line[j:j+4]
            grid[i][j] = EncodeOverhang(overhang)
        }
    }
    return lines, grid
}

func Write_results_to_file(path []byte, overhangs []string, stats string) {
//    filename := "results_"+time.Now().Format("15-04-05")
    filename := "results.txt"
    fmt.Printf("Writing results to \""+filename+"\".\n")
    results := ""
    results += stats
    results += "\nOverhangs: \n\n"
    if len(path) == len(overhangs) {
        for i, line := range overhangs {
            results = results + line[:path[i]]+" "+line[path[i]:]+"\n"
        }
    }
    ioutil.WriteFile(filename, []byte(results), 0777)
}
