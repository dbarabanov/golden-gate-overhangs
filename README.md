Fast implementation of Golden Gate DNA assembly method in Go. 
For background information refer to J5 paper: http://pubs.acs.org/doi/full/10.1021/sb2000116
This repository is focused on algorithm S4: "Search for the optimal set of Golden Gate assembly piece overhangs". It introduces a few ideas for speedup and parallelization and implements them efficiently in Go.

To run:

1. Install Go (http://golang.org/doc/install)
2. Supply your set of junctions in junctions.txt file (or use sample)
3. ./run.sh
4. See results in results.txt

Speedup ideas:

1. Prebuilt compatibility tables
There are only 4**4=256 possible 4-bp overhangs. We can pre-build a 256x256 bit compatibility table between all overhangs for different values of MAXIMUM_IDENTITIES_GOLDEN_GATE_OVERHANGS_COMPATIBLE setting. It only requires ~8Kb of memory to hold the bitmasks and we can do lookups very fast. 
Every 4-bp overhang is encoded by 1 byte and compatibility is determined by bit arithmetic.

2. Gradually building sets of inter-compatible overhangs
We can progressively build a set of inter-compatible overhangs: 
a. Take a first junction. Find all self-compatible overhangs in it. Put each one and it's partner in a new set.
b. Take the next junction. Select all self-compatible overhangs. Compare each overhang and it's partner with every set of overhangs from previous step. If they are compatible, create a new set with itself and partner. If not, discard.
c. Repeat b for all remaining junctions.
d. Select a set of overhangs with the best score and report.
If all overhangs were inter-compatible this algorithm would result in combinatorial explosion: every new junction would increase the size of sets by a factor of 5. I.e. 10 junctions would result in ~10mln sets. But because of compatibility constraints the size of sets never pushes beyond a few thousands which is manageable. In fact it starts dropping after a certain point because as the number of junctions goes up it gets progressively difficult to insert a new overhang into a set. 
A nice side effect of this approach is that a user no longer has to select MAXIMUM_IDENTITIES_GOLDEN_GATE_OVERHANGS_COMPATIBLE setting. We can quickly eliminate lower setting values because their compatibility factors are low.

3. Parallelization
Previous idea immediately suggests parallelization. A new junction can evaluate itself against previous junction's results independently. Go's concurrency primitives make this very easy to implement. We just wire a bunch of channels between adjacent junctions and send initial signal to the first layer. Parallelization is handled by a compiler. I get a 3x performance improvement by moving from 1 to 5 cores on my Mac. And the overall execution time is less than a second for 20 junctions with MAXIMUM_IDENTITIES_GOLDEN_GATE_OVERHANGS_COMPATIBLE=2
