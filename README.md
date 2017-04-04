# Random Beacon
[![Build Status](https://dfinity.build/job/beacon/badge/icon)](https://dfinity.build/job/beacon)

The Random Beacon is produced by a relay mechanism between groups where each group has set up a BLS threshold-signature key pair. The next randomness is derived from the unique deterministic threshold group signature on the previous randomness as the message. The mechanism is called a "Threshold-Relay Chain" and was developed for the DFINITY network. It is described in greater detail, e.g., here: http://dfinity.network/pdfs/viewer.html?file=../library/threshold-relay-blockchain-stanford.pdf
Some more background can be found in the "Technical" tab of this FAQ: http://dfinity.network/faq/

## Simulation Environment
You should first clone the `beacon` repository, and then run the `dfinity/build-env` docker container from inside the repository.
```
git clone git@github.com:dfinity/random-beacon.git
cd random-beacon
docker pull dfinity/build-env
docker run \
	--interactive \
	--rm \
	--tty \
	--volume $(pwd):/go/src/dfinity/beacon \
	--workdir /go/src/dfinity/beacon \
	dfinity/build-env:latest
```
Or short: `docker run --rm -it -v $(pwd):/go/src/dfinity/beacon -w /go/src/dfinity/beacon dfinity/build-env:latest`

## Run Simulation
### With default parameters

`go run main.go`

The default parameters are named in the output as follows:  
N = number of individual processes simulated  
m = number of groups formed  
n = group size (number of processes per group)  
k = group threshold  
l = number of blocks simulated after genesis block  
seed = numerical seed derived from seed string  

Also shown in the output are (all abbreviated to first two bytes):  
sec = secret key of process or group (for groups this is nil unless --bist is enabled)  
seed = internal randomness seed of process  
addr = address of process or group (as registered on the blockchain)  
pub = pubkey of process or group (as registered on the blockchain)  
mem = list of addresses of all group members  
sig = group signature created by the currently active group  
rnd = random beacon output produced by the currently active group  
grp = address of the group to be selected next  

Sample output:
```
BlkCh: (n)3 (k)2 (seed)d69198ea1c42e06a
--- Process setup: (N)8
Proc: (sec)3881 (seed)249e Node: (addr)b5da (pub)2 0xa4c3
Proc: (sec)1502 (seed)c809 Node: (addr)6042 (pub)2 0x20b2
Proc: (sec)3287 (seed)0f54 Node: (addr)2ac6 (pub)3 0x230b
Proc: (sec)2754 (seed)d88b Node: (addr)98b0 (pub)2 0x146b
Proc: (sec)1136 (seed)9aa4 Node: (addr)a991 (pub)2 0x20a9
Proc: (sec)1015 (seed)a598 Node: (addr)7d80 (pub)3 0x1488
Proc: (sec)2537 (seed)3f67 Node: (addr)ffc2 (pub)2 0x8e74
Proc: (sec)3080 (seed)6487 Node: (addr)75c9 (pub)3 0x3329
--- Group setup: (m)5
GrpP: (sec)<nil GrpR: (addr)66f9 (pub)2 0xa957 (n)3 (k)2 (mem)[ffc2,a991,6042]
GrpP: (sec)<nil GrpR: (addr)7819 (pub)2 0x1256 (n)3 (k)2 (mem)[6042,7d80,a991]
GrpP: (sec)<nil GrpR: (addr)253b (pub)2 0x3006 (n)3 (k)2 (mem)[6042,98b0,2ac6]
GrpP: (sec)<nil GrpR: (addr)2315 (pub)2 0x1f42 (n)3 (k)2 (mem)[75c9,6042,2ac6]
GrpP: (sec)<nil GrpR: (addr)2394 (pub)3 0x11a4 (n)3 (k)2 (mem)[75c9,98b0,ffc2]
--- Genesis block
1: Stat: (sig) (rnd)c5d2 (N)8 (m)5 (grp)253b
    0. Node: (addr)2ac6 (pub)3 0x230b
    1. Node: (addr)6042 (pub)2 0x20b2
    2. Node: (addr)75c9 (pub)3 0x3329
    3. Node: (addr)7d80 (pub)3 0x1488
    4. Node: (addr)98b0 (pub)2 0x146b
    5. Node: (addr)a991 (pub)2 0x20a9
    6. Node: (addr)b5da (pub)2 0xa4c3
    7. Node: (addr)ffc2 (pub)2 0x8e74
    0. GrpR: (addr)2315 (pub)2 0x1f42 (n)3 (k)2 (mem)[75c9,6042,2ac6]
    1. GrpR: (addr)2394 (pub)3 0x11a4 (n)3 (k)2 (mem)[75c9,98b0,ffc2]
    2. GrpR: (addr)253b (pub)2 0x3006 (n)3 (k)2 (mem)[6042,98b0,2ac6]
    3. GrpR: (addr)66f9 (pub)2 0xa957 (n)3 (k)2 (mem)[ffc2,a991,6042]
    4. GrpR: (addr)7819 (pub)2 0x1256 (n)3 (k)2 (mem)[6042,7d80,a991]
--- Blockchain states: (l)20
  2: Stat: (sig)2 0x22b2 (rnd)28f1 (N)8 (m)5 (grp)253b
  3: Stat: (sig)3 0x7f08 (rnd)928c (N)8 (m)5 (grp)2394
  4: Stat: (sig)3 0x8070 (rnd)e106 (N)8 (m)5 (grp)7819
  5: Stat: (sig)3 0x6903 (rnd)3af4 (N)8 (m)5 (grp)7819
  6: Stat: (sig)2 0x18d6 (rnd)db16 (N)8 (m)5 (grp)2315
  7: Stat: (sig)2 0x97a9 (rnd)dd57 (N)8 (m)5 (grp)2394
  8: Stat: (sig)2 0x1b74 (rnd)1ca2 (N)8 (m)5 (grp)66f9
  9: Stat: (sig)2 0x146f (rnd)acfb (N)8 (m)5 (grp)2315
 10: Stat: (sig)2 0x1678 (rnd)9f28 (N)8 (m)5 (grp)7819
 11: Stat: (sig)3 0x19b8 (rnd)500d (N)8 (m)5 (grp)7819
 12: Stat: (sig)3 0x1f97 (rnd)1147 (N)8 (m)5 (grp)253b
 13: Stat: (sig)2 0x216e (rnd)10da (N)8 (m)5 (grp)2394
 14: Stat: (sig)3 0x5084 (rnd)f634 (N)8 (m)5 (grp)2315
 15: Stat: (sig)2 0x1490 (rnd)4f95 (N)8 (m)5 (grp)2315
 16: Stat: (sig)2 0x18ff (rnd)a2ad (N)8 (m)5 (grp)253b
 17: Stat: (sig)3 0x1aea (rnd)fcd5 (N)8 (m)5 (grp)7819
 18: Stat: (sig)3 0x17f2 (rnd)5e9b (N)8 (m)5 (grp)66f9
 19: Stat: (sig)2 0xfcd8 (rnd)31b2 (N)8 (m)5 (grp)66f9
 20: Stat: (sig)3 0x1e9a (rnd)16f8 (N)8 (m)5 (grp)253b
 21: Stat: (sig)2 0x1082 (rnd)29e4 (N)8 (m)5 (grp)253b
 ```

### With custom parameters
Example:
`go run main.go -l=100 -N=50 -n=40 -k=20 -timing=true`

List of parameters:

* `-l` number of random outputs (blocks) to create (default 20)
* `-N` number of processes in pool (default 5)
* `-n` group size (default 3)
* `-k` group threshold (default 2)
* `-m` number of groups in pool (default 2)
* `-timing` flag to output timing information (default false)
* `-vvec` flag to run validation of verification vectors (default false)
* `-bist` flag to run built-in self tests (default false)

## Run test

`go test ./...`
 
## Run Benchmark

`go test --bench=. ./...`

Sample output for BN254:
```
BenchmarkPubkeyFromSeckey-4       	    5000	    313638 ns/op
BenchmarkSigning-4                	   10000	    105368 ns/op
BenchmarkValidation-4             	    2000	    583468 ns/op
BenchmarkDeriveSeckeyShare500-4   	  100000	     17605 ns/op
BenchmarkRecoverSeckey100-4       	    3000	    582981 ns/op
BenchmarkRecoverSeckey200-4       	    1000	   1770932 ns/op
BenchmarkRecoverSeckey500-4       	     200	   9178534 ns/op
BenchmarkRecoverSeckey1000-4      	      50	  33492710 ns/op
BenchmarkRecoverSignature100-4    	     200	   7813759 ns/op
BenchmarkRecoverSignature200-4    	     100	  15824394 ns/op
BenchmarkRecoverSignature500-4    	      30	  47239189 ns/op
BenchmarkRecoverSignature1000-4   	      10	 104978743 ns/op
```

The benchmark tests the speed of the underlying elliptic curve and pairing implementation by Shigeo Mitsunari (https://github.com/herumi/mcl).

Notably, we see the __signature validation time is 0.8 ms__ which involves a pairing evaluation.

We also see that __combining 500 signature shares into a group signature takes 60 ms__ (which would be used at a group size of 1000).


Sample output for BN382_1:
```
BenchmarkPubkeyFromSeckey-4       	    2000	    922435 ns/op
BenchmarkSigning-4                	    5000	    273455 ns/op
BenchmarkValidation-4             	    1000	   1717128 ns/op
BenchmarkDeriveSeckeyShare500-4   	   50000	     30183 ns/op
BenchmarkRecoverSeckey100-4       	    2000	    979008 ns/op
BenchmarkRecoverSeckey200-4       	     500	   3070395 ns/op
BenchmarkRecoverSeckey500-4       	     100	  15512509 ns/op
BenchmarkRecoverSeckey1000-4      	      20	  59512824 ns/op
BenchmarkRecoverSignature100-4    	     100	  21427153 ns/op
BenchmarkRecoverSignature200-4    	      30	  44704722 ns/op
BenchmarkRecoverSignature500-4    	      10	 118874472 ns/op
BenchmarkRecoverSignature1000-4   	       5	 268353411 ns/op
```

## Dependencies

The dependencies below are all met in the docker image above.

### go-ethereum

Code currently depends on `github.com/ethereum/go-ethereum/common` being present in the `src` directory.

### cgo bindings

For cgo, which is transitioning in, we need the environment variables set:

`export LIBRARY_PATH=/build/herumi/bls/lib:/build/herumi/mcl/lib:$LIBRARY_PATH`

`export CPATH=/build/herumi/bls/include:$CPATH`
