package main

import (
    . "github.com/dbarabanov/golden_gate_overhangs/utils"
    . "github.com/dbarabanov/golden_gate_overhangs/goldengate"
    "fmt"
    "time"
)

func main() {
    var TOTAL_JUNCTIONS byte = 8
//    var TOTAL_CUTS byte = 5
//    var DEFAULT_CHANNEL_BUFFER_SIZE byte = 10

    fmt.Printf("\nSinc starting....\n")
    sinc := CreateSinc()
    
    node7 := CreateNode(EncodeOverhang("CGAT"), 7, 0, 7)
    WireSinc(node7, &sinc)
    go RunNode(node7)

    node8 := CreateNode(EncodeOverhang("GGAT"), 3, 0, 3)
    WireSinc(node8, &sinc)
    go RunNode(node8)

    go RunSinc(&sinc, 2)

    node7.Input<-CreateInitialSignal(TOTAL_JUNCTIONS)
    node7.Input<-nil

    sig := CreateInitialSignal(TOTAL_JUNCTIONS)
    node8.Input<-sig
    node8.Input<-sig
    node8.Input<-nil

    fmt.Printf("Closing Sinc: %v\n", <-sinc.Output)


    time.Sleep(100)
}

func Test_utils() {
    seq := "TAGC"
    fmt.Printf("Encoding: %v\n", seq)
    encoded4 := Encode4(seq)
    fmt.Printf("Encoded: %0#8b\n", encoded4)
    fmt.Printf("Decoded: %v\n", Decode4(encoded4))
    fmt.Printf("A: %c\n", 'A')
    fmt.Printf("Encoded A: %c\n", Encode1("A"))
    fmt.Printf("Encoded A: %b\n", Encode1("A"))
    fmt.Printf("Encoded A: %0#8b\n", EncodeBase('A'))
    fmt.Printf("Encoded C: %0#8b\n", EncodeBase('C'))
    overhang := "ATCG" 
    fmt.Println("Overhang: "+overhang)
    encoded := EncodeOverhang(overhang)
    fmt.Printf("Encoded: %0#8b\n", encoded)
    decoded := DecodeOverhang(encoded)
    fmt.Printf("Decoded: %v\n", decoded)
    o1 := "ATCT"
    o2 := "ACAC"
    b1 := EncodeOverhang(o1)
    b2 := EncodeOverhang(o2)
//    fmt.Printf("IsCompatible("+o1+", "+o2+"): %v\n", IsCompatible(b1, b2))
    fmt.Printf("Complementary("+o1+"): %v\n", ComplementaryOverhang(o1))
    fmt.Printf("Reverse("+o1+"): %v\n", ReverseOverhang(o1))
    fmt.Printf("Partner("+o1+"): %v\n", PartnerOverhang(o1))
    fmt.Printf("GetRepeatCount("+o1+", "+o2+"): %v\n", GetRepeatCount(b1, b2))
    fmt.Printf("AreOverhangsCompatible("+o1+", "+o2+", 2): %v\n", AreOverhangsCompatible(o1, o2, 2))
    fmt.Printf("IsOverhangSelfCompatible("+o1+", 2): %v\n", IsOverhangSelfCompatible(o1, 2))
    var max_tolerable_repeats byte = 1
//    var total_self_incompatible = 0
    var i byte
//    for i=0; i<=255; i++ {
//        is_self_compatible := IsSelfCompatible(i, max_tolerable_repeats)
//        if !is_self_compatible {total_self_incompatible++}
//        fmt.Printf("%v, %v\n", DecodeOverhang(i), is_self_compatible)
//        if i==255 {break}
//    }
//    fmt.Printf("total_self_incompatible: %v\n", total_self_incompatible)
    var j byte
    var incompatible_row bool
    for i=0; i<5; i++ {
        incompatible_row = !IsSelfCompatible(i, max_tolerable_repeats)
        fmt.Printf("%v: ", DecodeOverhang(i))
        count := 0
        for j=0; j<70; j++ {
            if incompatible_row  || !IsSelfCompatible(j, max_tolerable_repeats) {
                fmt.Printf("0")
            } else if AreCompatible(i, j, max_tolerable_repeats) {
        //                fmt.Printf("%#6v ", )
                fmt.Printf("1")
                count++
            } else {fmt.Printf("0")}
            if j==255 {break}
        }
        fmt.Printf(" %v", 255-count)
        fmt.Printf("\n")
    }
    max_tolerable_repeats = 2
    overhangs := []string{"AAAA","CCCC", "ATAA"}
    fmt.Printf("IsOverhangCompatibleWithMany("+o1+", %v, "+string(max_tolerable_repeats)+"): %v\n", overhangs,  IsOverhangCompatibleWithMany(o1, overhangs, max_tolerable_repeats))
}
