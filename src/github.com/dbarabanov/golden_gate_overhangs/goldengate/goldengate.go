package goldengate

import (
    "fmt"
    )

const CUTS byte = 5

type Signal struct {
    Overhangs []byte
    Path []byte
    Cost byte
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
    Cost byte
    Level byte
    Index byte
    SignalCounter int
}

func RunSinc(sinc *Sinc) {
    nil_threshold := CUTS
    var nil_counter byte
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
    var nil_counter byte
    for {
        received := <-node.Input
        if received != nil {
//           fmt.Printf("Node (%v, %v) received: %v\n", node.Level, node.Index, received)
            sig := CreateInitialSignal(byte(len(received.Path)))
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
                fmt.Printf("Node (%v, %v) received nil. \n", node.Level, node.Index)
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

func CreateNode(overhang byte, cost byte, level byte, index byte) *Node {
    node := Node{}
    node.Overhang = overhang
    node.Cost = cost
    node.Input = make(chan *Signal)
    node.Level = level
    node.Index = index
    return &node
}

func CreateLayer(level byte) []*Node {
    nodes := make([]*Node, CUTS, CUTS)
    for i := range nodes {
        nodes[i] = CreateNode(byte(i)*level, byte(i), level, byte(i))
    }
    return nodes
}

func CreateInitialSignal(max_levels byte) *Signal {
    sig := Signal{}
    sig.Cost = 0
    sig.Overhangs = make([]byte, 0, max_levels)
    sig.Path = make([]byte, 0, max_levels)
    return &sig
}

func CreateSinc() *Sinc {
    return &Sinc{make(chan *Signal), make(chan string), 0}
}

func SendInitialSignals(nodes []*Node, level_depth byte) {
    for i:= range nodes {
        nodes[i].Input<-CreateInitialSignal(level_depth)
        for _ = range nodes{
            nodes[i].Input<-nil
        }
    }
}
