package goldengate

import (
    "fmt"
    )

var CUTS int = 5

type Signal struct {
    Overhangs []byte
    Path []byte
    Cost int
}

type Sinc struct {
    Input chan *Signal
    Output chan string
    SignalCounter int
}

type Node struct {
    Input chan *Signal
    Output []chan *Signal
    Overhang byte
    Cost int
    Level int
    Index byte
    SignalCounter int
}

func RunSinc(sinc *Sinc) {
    nil_threshold := CUTS
    var nil_counter int 
    for {
        sig := <-sinc.Input
        if sig != nil {
            fmt.Printf("sinc received: %v\n", *sig)
            sinc.SignalCounter ++
        } else {
            fmt.Printf("sinc received: %v\n", sig)
            nil_counter++
            if nil_counter >= nil_threshold {
                fmt.Printf("sinc is closed.\n")
                sinc.Output<-fmt.Sprintf("Sinc total signals: %v", sinc.SignalCounter)
                return
            }
        }
    }
}

func RunNode(node *Node) {
    nil_threshold := CUTS
    var nil_counter int
    for {
        received := <-node.Input
        if received != nil {
//           fmt.Printf("Node (%v, %v) received: %v\n", node.Level, node.Index, received)
            sig := CreateInitialSignal(len(received.Path))
            sig.Cost = received.Cost + node.Cost
            sig.Overhangs = make([]byte, len(received.Overhangs), len(received.Overhangs))
            copy(sig.Overhangs, received.Overhangs)
            sig.Overhangs = append(sig.Overhangs, node.Overhang)
            sig.Path = make([]byte, len(received.Path), len(received.Path))
            sig.Path = append(sig.Path, node.Index)
            copy(sig.Path, received.Path)
            node.SignalCounter ++
            BroadcastSignal(node, sig)
       } else {
            nil_counter++
//                fmt.Printf("Node (%v, %v) received nil. \n", node.Level, node.Index)
            if nil_counter >= nil_threshold {
                fmt.Printf("Node (%v, %v) final SingalCounter: %v\n", node.Level, node.Index, node.SignalCounter)
                BroadcastSignal(node, nil)
                return
            }
       }
   }
}

func BroadcastSignal(node *Node, sig *Signal) {
    for i := range node.Output {
        node.Output[i]<-sig
    }
}

func RunNodes(nodes []*Node) {
    for i := range nodes {
        go RunNode(nodes[i])
    }
}

func WireSinc(node *Node, sinc *Sinc) {
   node.Output = make([]chan *Signal, 1, 1)
   node.Output[0] = sinc.Input
}

func WireNodesToSinc(nodes []*Node, sinc *Sinc) {
    for i := range nodes {
        WireSinc(nodes[i], sinc)
    }
}

func WireNodeToNodes(sender *Node, receiver []*Node) {
   sender.Output = make([]chan *Signal, len(receiver), len(receiver))
   for i := range sender.Output {
       sender.Output[i] = receiver[i].Input
   }
}

func WireLevels(lower []*Node, higher[]*Node) {
    for i:= range lower {
        WireNodeToNodes(lower[i], higher)
    }
}

func CreateNode(overhang byte, cost int, level int, index byte) *Node {
    node := Node{}
    node.Overhang = overhang
    node.Cost = cost
    node.Input = make(chan *Signal)
    node.Level = level
    node.Index = index
    return &node
}

func CreateRandomLevel(level int) []*Node {
    nodes := make([]*Node, CUTS, CUTS)
    for i := range nodes {
        nodes[i] = CreateNode(byte(i*level), i, level, byte(i))
    }
    return nodes
}

func CreateLevel(overhangs []byte, level int) []*Node {
    nodes := make([]*Node, len(overhangs), len(overhangs))
    for i := range nodes {
        nodes[i] = CreateNode(overhangs[i], Cost_by_index(i), level, byte(i))
    }
    return nodes
}

func Cost_by_index (i int) int{
    if i<CUTS/2 {
        return CUTS/2-i
    }
    return i-CUTS/2
}

func CreateInitialSignal(max_levels int) *Signal {
    sig := Signal{}
    sig.Cost = 0
    sig.Overhangs = make([]byte, 0, max_levels)
    sig.Path = make([]byte, 0, max_levels)
    return &sig
}

func CreateSinc() *Sinc {
    return &Sinc{make(chan *Signal), make(chan string), 0}
}

func SendInitialSignals(nodes []*Node, level_depth int) {
    for i:= range nodes {
        nodes[i].Input<-CreateInitialSignal(level_depth)
        for _ = range nodes{
            nodes[i].Input<-nil
        }
    }
}

func BuildRandomGrid(total_levels int) ([][]*Node, *Sinc){
    levels := make([][]*Node, total_levels, total_levels)
    for i := range levels {
        levels[i] = CreateRandomLevel(i+1)
    }

    sinc := CreateSinc()
    
    for i := range levels{
        if i < total_levels-1 {
            WireLevels(levels[i], levels[i+1])
        } else {
            WireNodesToSinc(levels[i], sinc)
        }
    }

    for i := range levels {
        RunNodes(levels[i])
    }

    go RunSinc(sinc)
    return levels, sinc
}

func BuildGrid(input [][]byte) ([][]*Node, *Sinc){
    CUTS = len(input[0])
    for _, level := range(input){
        if len(level) != CUTS {
            panic("Bad input grid")
        }
    }
    total_levels := len(input)
    levels := make([][]*Node, total_levels, total_levels)
    for i := range levels {
        levels[i] = CreateLevel(input[i], i+1)
    }

    sinc := CreateSinc()
    
    for i := range levels{
        if i < total_levels-1 {
            WireLevels(levels[i], levels[i+1])
        } else {
            WireNodesToSinc(levels[i], sinc)
        }
    }

    for i := range levels {
        RunNodes(levels[i])
    }

    go RunSinc(sinc)
    return levels, sinc
}
