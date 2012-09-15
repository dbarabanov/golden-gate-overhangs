package goldengate

import (
//    . "github.com/dbarabanov/golden_gate_overhangs/utils"
    "fmt"
    )

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
    Data byte
    Cost byte
    Level byte
    Index byte
    SignalCounter int
}

type Source struct {
    Output []chan Signal
}

func RunSinc(sinc *Sinc, nil_threshold byte) {
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
   for {
       received := <-node.Input
       if received != nil {
//           fmt.Printf("Node (%v, %v) received: %v\n", node.Level, node.Index, received)
           received.Cost += node.Cost
           received.Overhangs = append(received.Overhangs, node.Data)
//           received.Path[node.Level] = node.Index
           received.Path = append(received.Path, node.Index)
           BroadcastSignal(node, received)
           node.SignalCounter ++
       } else {
           fmt.Printf("Node (%v, %v) final SingalCounter: %v\n", node.Level, node.Index, node.SignalCounter)
           BroadcastSignal(node, nil)
           return
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

func CreateNode(data byte, cost byte, level byte, index byte) *Node {
    node := Node{}
    node.Data = data
    node.Cost = cost
    node.Input = make(chan *Signal)
    node.Level = level
    node.Index = index
    return &node
}

func CreateInitialSignal(max_levels byte) *Signal {
    sig := Signal{}
    sig.Cost = 0
    sig.Overhangs = []byte{} 
    sig.Path = make([]byte, 0, max_levels)
    return &sig
}

func CreateSinc() *Sinc {
    return &Sinc{make(chan *Signal), make(chan string), 0}
}

func CreateSource(channels_count byte) Source {
    return Source{make([]chan Signal, channels_count, channels_count)}
}
