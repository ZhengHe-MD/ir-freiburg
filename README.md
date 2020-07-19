# Information Retrieval Course Notes

[![Build Status](https://travis-ci.org/ZhengHe-MD/ir-freiburg.svg?branch=master)](https://travis-ci.org/ZhengHe-MD/ir-freiburg)

I happened to found this course about Information Retrieval(IR) from Youtube long time ago, and save it to 
my playlist, as usual. Recently, I've just got started to work with ElasticSearch and want to get some background knowledge.
So I take out this course from my playlist, and it turns out to be a great course.   

All original materials are collected from GitHub and the course website, and I rewrite them using Go because
I use it for my daily work. If you found this repo useful, feel free to fork and play around with it. 

## Details

### Lecture-01 ✅

##### Topics

* Keyword Search
* Inverted Index
* One, Two and More Words
* Zipf's Law

In-class demo and exercise code can be found in [lecture-01 directory](./lecture-01). The [script.sh](./lecture-01/script.sh) contains all runnable examples you need. 

### Lecture-02 ✅

##### Topics

* Ranking
  * Term Frequency (tf)
  * Document Frequency (df)
  * tf.idf
  * BM25 (best match)
* Evaluation
  * Precision (P@K)
  * Average Precision (AP)
  * Mean Precisions (MP@k, MP@R, MAP)
  * Discounted Cumulative Gain (DCG)
  * Binary Preference (bpref)

In-class demo and exercies code can be found in [lecture-02 directory](./lecture-02). The [script.sh](./lecture-02/script.sh) contains the command to benchmark on movies dataset. It's counter-intuitive that the provided [movies-benchmark.txt](./data/movies-benchmark.txt) start counting docID at 2, which conflicts with the provided unit test cases in [TIP file](./lecture-02/sheet-02.TIP) either. So I write a [script](./data/process_movies_benchmark.go) to process the movies-benchmark.txt, make it start counting docID at 1, the result benchmark file [movies-benchmark-minus-1.txt](./data/movies-benchmark-minus-1.txt) is also provided in the [data directory](./data).

### Lecture-03 ✅

* List intersection
  * Intersection and merge
  * Time measurement (repeat 5 times at least)
* Non-algorithmic improvements
  * Naive arrays, ArrayList in Java, std::vector in C++
  * Predictable branches (cpu pipelining, reduce guessing)
  * Sentinels (avoid condition testing)
* Algorithmic improvements
  * Preliminaries: smaller list A with k elements, longer list B with n elements
  * Binary search in the longer list, θ(klogn)
  * Binary search in the remainder of longer list, best case θ(k+logn), worst case θ(klogn), average case θ(klogn)
  * Galloping search, O(klog(1+n/k))
  * Skip Pointers

| optimization strategy | ns/op |
|:----------------------|:------|
|BenchmarkIntersectBasic-4                              |10584688 ns/op |
|BenchmarkIntersectWithLessConditionalParts-4           |7265918 ns/op |
|BenchmarkIntersectWithSentinels-4                      |7008293 ns/op |
|BenchmarkIntersectWithBinarySearchInLongerRemainder-4  |7649609 ns/op |
|BenchmarkIntersectWithGallopingSearch-4                |6083991 ns/op |
|BenchmarkIntersectWithSkipPointer-4                    |10583605 ns/op |
|BenchmarkIntersectHybrid-4                             |5345517 ns/op |

title of the film: [The Big Lebowski](https://en.wikipedia.org/?curid=29782).

### Lecture-04 🚧

* Compression
  * Motivation
    * Large inverted list in web-scale collection. Try a few keyword searches on Google.
    * Compression can saves space and time.
  * Memory & Disk
    * Index in memory
      * Compression saves memory.
      * More index can be fit into memory, which is good because memory ops are much faster than disk ops.
    * Index on disk
      * Compression saves disk space.
      * Compression saves query time. Spend less time on i/o, but more time on cpu, the later usually outweights the former.
  * Gap encoding
    * Idea: stores differences (always positive integers) instead of raw doc ids.
    * Binary representation
      * takes $(\left \lfloor{log_2{x}}\right \rfloor + 1)$ bits for each number.
      * but there are multiple ways to decode an inverted list with binary representation.
    * Prefix-free codes: decoding from left to right is unambiguous.
* Codes
  * Elias-Gamma (1975)
    * code: $\left \lfloor{log_2{x}}\right \rfloor$ zeros + x in binary.
    * length: exactly $(2\left \lfloor{log_2{x}}\right \rfloor + 1)$ bits.
    * why prefix-free?
  * Elias-Delta (1975)
    * code: $\left \lfloor{log_2{x}}\right \rfloor + 1$ in Elias-Gamma + x in binary without the leading 1
    * length: $\left \lfloor{log_2{x}}\right \rfloor + 2log_2{log_2{x}} + O(1)$ bits
    * why prefix-free? Elias-Gamma code tells exactly how many bits of the code come afterwards
  * Golomb (1966)
    * parameter: M (integer), called modulus
    * code:
      * $x = q*M + r$, $q = x \space div \space M$, $r = x \space mod \space M$
      * concatenation of:
        * q written in unary with 0s
        * a single 1 (as a delimiter)
        * r written in binary
    * length: $\left \lfloor{\frac{x}{M}}\right \rfloor + 1 + \left \lceil{log_2{M}}\right \rceil$
    * why prefix-free?
  * Variable-Bytes (VB)
    * idea
      * Use whole bytes and avoid bit fidding, just like UTF-8 encoding.
      * use one bit of each byte to indicate whether this is the last byte in the current code or not.
  * Other codes: ANS
* Entropy (theory part)
  * Motivation: Which code compresses the best? It depends on the distributions. The intuition is that **more frequent less bits**.
  * Definition: 
    * Entropy of a discrete random variable $X$, which generates the symbols of a message.
    * $H(x) = -\sum_{i}p_{i}log_2{p_i}$, where $p_i = Prob(X = i)$, $H(x)$ is the optimal number of bits to encode a random symbol generated according to $X$.
  * Shannon's source coding theorem (1948)
    * For an arbitrary prefix-free (PF) encoding, let $L_i$ be the length of the code for $i \in range(X)$
      * For any PF encoding it holds: $E(L_x) \ge H(x)$
      * There is a PF encoding with: $E(L_x) \le H(x) + 1$
    * In plain words: **No code can be better than the entropy, and there is always a code that is (almost) as good**.


## References
* [videos](https://www.youtube.com/playlist?list=PLfgMNKpBVg4V8GtMB7eUrTyvITri8WF7i)
* [course wiki](https://ad-wiki.informatik.uni-freiburg.de/teaching/InformationRetrievalWS1718)
* [lecture svn repo](https://daphne.informatik.uni-freiburg.de/ws1718/InformationRetrieval/svn-public/public/)
* Text books
  * [Introduction to Information Retrieval](https://nlp.stanford.edu/IR-book/information-retrieval-book.html)
