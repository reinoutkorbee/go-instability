# go-instability
Instability metric for go source files

According to Wikipedia: 

Instability (I): The ratio of efferent coupling (Ce) to total coupling (Ce + Ca) such 
that I = Ce / (Ce + Ca). This metric is an indicator of the package's resilience to 
change. The range for this metric is 0 to 1, with I=0 indicating a completely stable 
package and I=1 indicating a completely unstable package.

The efferent coupling, or Ce, is measured by counting the outward pointing
connections. The afferent coupling, or Ca, is measured by counting the inward
pointing connections.

Adapted for Go, this program counts exported methods as inward connections of a Go source
file. The method calls to another package's exported methods are counted as outward
connections.

The following options are available:

-debug
	enable verbose logging to stdout
	
-pkgs 
	comma separated list of package roots relative to the GOPATH to include
	
	
Usage:
	go-instability -debug -pkgs src/mycom.com/projectA,src/mycom.com/projectB src/mycom.com
	
The above example ignores all method calls to packages not in projectA and projectB. This
feature is useful if you want to only calculate the instability in relation to source 
code that is actually under your control.

The output is in CSV format on stdout.

For example:

> ./go-instability -pkgs src/github.com  src/github.com/eapache | sort -n
PATH, INST.
src/github.com/eapache/go-resiliency/batcher/batcher.go, 0.250000
src/github.com/eapache/go-resiliency/batcher/batcher_test.go, 0.791667
src/github.com/eapache/go-resiliency/breaker/breaker.go, 0.700000
src/github.com/eapache/go-resiliency/breaker/breaker_test.go, 0.915254
src/github.com/eapache/go-resiliency/deadline/deadline.go, 0.500000
src/github.com/eapache/go-resiliency/deadline/deadline_test.go, 0.777778
src/github.com/eapache/go-resiliency/retrier/backoffs.go, 0.000000
src/github.com/eapache/go-resiliency/retrier/backoffs_test.go, 0.846154
src/github.com/eapache/go-resiliency/retrier/classifier.go, 0.000000
src/github.com/eapache/go-resiliency/retrier/classifier_test.go, 0.833333
src/github.com/eapache/go-resiliency/retrier/retrier.go, 0.500000
src/github.com/eapache/go-resiliency/retrier/retrier_test.go, 0.806452
src/github.com/eapache/go-resiliency/semaphore/semaphore.go, 0.333333
src/github.com/eapache/go-resiliency/semaphore/semaphore_test.go, 0.750000
src/github.com/eapache/go-xerial-snappy/fuzz.go, 0.000000
src/github.com/eapache/go-xerial-snappy/snappy.go, 0.555556
src/github.com/eapache/go-xerial-snappy/snappy_test.go, 0.711538
src/github.com/eapache/queue/queue.go, 0.000000
src/github.com/eapache/queue/queue_test.go, 0.476190

The implementation is rather rudimentary and there are plenty of improvements to be made, for example:

1. Exclude standard library
2. Handle package aliases
3. Improve matching of selectors to package names
4. Include exported structs as well
5. Provide option to exclude certain packages


