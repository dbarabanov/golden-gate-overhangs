package utils

import (
    "fmt"
    "math/rand"
    "io/ioutil"
    "strings"
    "time"
    )

//Encodes nucleotide letter (A,C,T,G) into byte (0,1,2,3). Panics if input is not one of A,C,T,G.
//
//  EncodeBase('C') = 1
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
    panic(fmt.Sprintf("Invalid argument to EncodeBase(): '%c'. Must be one of A,C,T,G.", letter))
}

//Decode byte(0,1,2,3) into nucleotide letter (A,C,T,G). Panics if input is not 0,1,2, or 3.
//
//  DecodeBase(1) = "C"
func DecodeBase(encoded byte) (letter string) {
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
    panic(fmt.Sprintf("Invalid argument to DecodeBase(): %v. Must be one of 0,1,2,3.", encoded))
}

//EncodeOverhang encodes 4-bp overhang into byte. 2 most significant bits in the resulting byte correspond to first base pair in overhang. Similar to least signigicat bits and 2 middle bit pairs. See also EncodeBase.
//
//  EncodeOverhang("AAAC") = 1
//  EncodeOverhang("TTTT") = 255
func EncodeOverhang(seq string) (encoded byte) {
    if len(seq) != 4 {
        panic(fmt.Sprintf("Invalid argument to EncodeOverhang(): %v. Must be a 4-letter string of A,C,T,G.", seq))
    }
    return (EncodeBase(seq[0])<<6)|(EncodeBase(seq[1])<<4)|(EncodeBase(seq[2])<<2)|(EncodeBase(seq[3]))
}

//DecodeOverhang decodes byte into 4-bp overhang. It's a reverse function of EncodeOverhang.
//
//  DecodeOverhang(1) = "AAAC"
//  DecodeOverhang(255) = "TTTT" 
func DecodeOverhang(b byte) (s string) {
    return DecodeBase(b>>6) + DecodeBase(b<<2>>6) + DecodeBase(b<<4>>6) + DecodeBase(b<<6>>6)
}

//ComplementaryOverhang returns an overhang that is complementary to the overhang passed as argument.
//
//  ComplementaryOverhang("ACCG") = "TGGC"
func ComplementaryOverhang(overhang string) string {
    return DecodeOverhang(Complementary(EncodeOverhang(overhang)))
}

//Complementary returns byte-encoded complementary overhang. See ComplementaryOverhang.
func Complementary(b byte) byte {
    return 255-b 
}

//ReverseOverhang returns an overhang that is a reverse of the overhang passed as an argument.
//
//  ReverseOverhang("ACCG") = "GCCA"
func ReverseOverhang(overhang string) string {
    return DecodeOverhang(Reverse(EncodeOverhang(overhang)))
}

//Reverse returns byte-encoded reverse overhang. See ReverseOverhang.
func Reverse(b byte) byte {
    return b&3<<6 + b&12<<2 + b&48>>2 + b&192>>6
}

//PartnerOverhang returns reverse-complement overhang.
//
//  PartnerOverhang("ACCG") = CGGT
func PartnerOverhang(overhang string) string {
    return DecodeOverhang(Partner(EncodeOverhang(overhang)))
}

//Partner returns byte-encoded reverse-complement overhang. See PartnerOverhang.
func Partner(b byte) byte {
    return Reverse(Complementary(b)) 
}

//GetOverhangRepeatCount returns the number of nucleotides that are the same in the same positions in 2 overhangs.
//
//  GetOverhangRepeatCount("AAAT", "AAAA") = 3
//  GetOverhangRepeatCount("CCGG", "GGCC") = 0
func GetOverhangRepeatCount(overhang1 string, overhang2 string) byte{
    return GetRepeatCount(EncodeOverhang(overhang1), EncodeOverhang(overhang2))
}

//GetRepeatCount returns the number of repeats between 2 byte-encoded overhangs. See GetOverhangRepeatCount.
func GetRepeatCount(b1 byte, b2 byte) byte {
//+1 for same base at each of 4 position.
    return isZero((b1^b2)&3) +
        isZero((b1^b2)&12) +
        isZero((b1^b2)&48)+
        isZero((b1^b2)&192)
//    return 4-(
//        (((b1^b2)&1)|((b1^b2)&2>>1)) + 
//        (((b1^b2)&4>>2)|((b1^b2)&8>>3)) + 
//        (((b1^b2)&16>>4)|((b1^b2)&32>>5)) + 
//        (((b1^b2)&64>>6)|((b1^b2)&128>>7)))
}

//can't cast bool to byte in Go. This little function to the rescue.
func isZero(b byte) byte {
    if b==0 {return 1}
    return 0
}

func AreOverhangsCompatible(overhang1 string, overhang2 string, max_repeats byte) bool {
    return AreCompatible(EncodeOverhang(overhang1), EncodeOverhang(overhang2), max_repeats)
}

func AreCompatible(b1 byte, b2 byte, max_repeats byte) bool {
    return GetRepeatCount(b1, b2) <= max_repeats 
}

func IsOverhangSelfCompatible(overhang string, max_repeats byte) bool{
    return IsSelfCompatible(EncodeOverhang(overhang), max_repeats)
}

func IsSelfCompatible(b byte, max_repeats byte) bool {
    return AreCompatible(b, Partner(b), max_repeats)
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
