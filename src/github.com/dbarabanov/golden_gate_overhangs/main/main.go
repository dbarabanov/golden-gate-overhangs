package main

import (
    . "github.com/dbarabanov/golden_gate_overhangs/utils"
    . "github.com/dbarabanov/golden_gate_overhangs/goldengate"
    "fmt"
    "time"
    "math"
    "runtime"
)

func main() {

    const TOTAL_JUNCTIONS = 3

    runtime.GOMAXPROCS(6)
    const CUTS = 5
    const INFINITE_COST = 999999
    const DATA_FROM_FILE = true

//    input byte[][] := {{1,2,3,4,5},{6,7,8,9,10}}
//    input := [][]byte{{11,12,13,14,15},{26, 27,28,29,210}}
    var overhangs []string    
    var input [][]byte

    input = GenerateRandomGrid(CUTS, TOTAL_JUNCTIONS)
    if DATA_FROM_FILE {
        overhangs, input = GridFromFile("junctions.txt")
    }
    

    fmt.Printf("input: %v\n", input)

    Set_infinite_cost(INFINITE_COST)
    best_cost := INFINITE_COST
    var Max_repeats byte = 0
    var sig Signal
    var sinc_signals int

    start_time := time.Now()
    for best_cost == INFINITE_COST && Max_repeats < 4 {
        Set_max_repeats(Max_repeats)
        fmt.Printf("Starting evaluation for MAX_REPEATS = %v\n", Max_repeats)
    //    grid, sinc := BuildRandomGrid(TOTAL_JUNCTIONS)
        grid, sinc := BuildGrid(input)

        SendInitialSignals(grid[0], len(grid))

        sig = <-sinc.Best_signal
        sinc_signals = <-sinc.Total_signals
        Max_repeats += 1
        best_cost = sig.Cost
    }
    run_time := time.Since(start_time)
    fmt.Printf("Best signal: %v\n", sig)
    fmt.Printf("Execution time: %v\n", run_time)
    fmt.Printf("Total signals: %v\n", sinc_signals)
    max_signals_possible := int64(math.Pow(CUTS, TOTAL_JUNCTIONS))
    fmt.Printf("Max possible signal: %v\n", max_signals_possible)
    fmt.Printf("Signal throughput: %v\n", float64(sinc_signals)/float64(max_signals_possible))
    fmt.Printf("Best neutrality score: %v\n", sig.Cost)
    fmt.Printf("BEST_MAX_REPEATS: %v\n", Max_repeats-1)

    stats := ""
    if sig.Cost == INFINITE_COST {
        stats += "Input set of junctions is incompatible."
    } else {
        stats += fmt.Sprintf("Best neutrality score: %v\n", sig.Cost)
        stats += fmt.Sprintf("MAXIMUM_IDENTITIES_GOLDEN_GATE_OVERHANGS_COMPATIBLE: %v\n", Max_repeats-1)
        stats += fmt.Sprintf("Execution time: %v\n", run_time)
    }
    if DATA_FROM_FILE {
        Write_results_to_file(sig.Path, overhangs, stats)
    }

    time.Sleep(1000)
//    Test_utils()
}

func Test_utils() {
    seq := "TAGC"
    fmt.Printf("Encoding: %v\n", seq)
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
