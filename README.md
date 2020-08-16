# Information Retrieval Course Notes

[![Build Status](https://travis-ci.org/ZhengHe-MD/ir-freiburg.svg?branch=master)](https://travis-ci.org/ZhengHe-MD/ir-freiburg)

I happened to found this course about Information Retrieval(IR) from Youtube long time ago, and save it to 
my playlist, as usual. Recently, I've just got started to work with ElasticSearch and want to get some background knowledge.
So I take out this course from my playlist, and it turns out to be a great course.   

All original materials are collected from GitHub and the course website. I rewrite the provided codes using Go because
I use it for my daily work. If you found this repo useful, feel free to fork and play around with it. All paper and pencil solutions
are written in markdown by myself, and I also provide pdf version for each. Note that the markdown editor I use for solutions
is [typora](typora.io).

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

### Lecture-04 ✅

* Compression
  * Motivation
  * Memory & Disk
  * Gap encoding
    * Idea: stores differences (always positive integers) instead of raw doc ids.
    * Binary representation
    * Prefix-free codes: decoding from left to right is unambiguous.
* Codes
  * Elias-Gamma (1975)
  * Elias-Delta (1975)
  * Golomb (1966)
  * Variable-Bytes (VB)
  * Other codes: ANS
* Entropy (theory part)
  * Motivation: Which code compresses the best? It depends on the distributions. The intuition is that **more frequent less bits**.
  * Shannon's source coding theorem (1948)
    * In plain words: **No code can be better than the entropy, and there is always a code that is (almost) as good**.
    
### Lecture-05 ✅

* Fuzzy Search
* Edit distance
* q-Gram Index

Result Table:
| Query | Time | #PED | #RES | Language | Processor/RAM |
|-------|------|------|------|----------|---------------|
| the   | 199ms| 139  | 139  | Go       | 2.7 GHz Dual-Core Intel Core i5 |
| breib | 830ms| 189857 | 279 | Go      | 2.7 GHz Dual-Core Intel Core i5 |
| the BIG lebauski | 133ms | 669 | 1 | Go | 2.7 GHz Dual-Core Intel Core i5 |

### Lecture 06-07 ❌

Lecture 06 and 07 are mostly about the html, javascript and css stuff, which I've been familiar with. So I decide to skip these two lectures. There is a very clear and intuitive discussion about UTF-8 in lecture 07, the dominant encoding scheme in the web, and it's worth reading. Though the content is located in the slides of lecture 07, the teacher actually walks though that in the beginning of lecture 08.

### Lecture 08 ✅

* Vector Space Model (VSM)
  * dot product of query vector and term-document matrix
  * sparse matrices (inverted index itself)
    * row-major
    * column-major
  * normalization (similar to idf part from BM25)
    * L1-norm: sum of the absolutes of the entries
    * L2-norm: sum of the squares of the entries

I use the library [james-bowman/sparse](github.com/james-bowman/sparse) to do the sparse matrix stuff. The following list shows the benchmarking results of different config setup from best to worst: (BM params are fixed to b = 0.75, k = 1.25). 

| Score Type | Normalization | MAP   |
| ---------- | ------------- | ----- |
| BM25       | None          | 0.318 |
| BM25WithoutIDF | Row-wise L2 | 0.295 |
| BM25WithoutIDF | Row-wise L1 | 0.210 |
| BM25       | L2           | 0.191 |
| TF.IDF     | None          | 0.190 |
| BM25WithoutIDF | None | 0.175 |
| TF.IDF     | L2            | 0.157 |
| BM25       | L1            | 0.054 |
| TF.IDF     | L1            | 0.050 |

The row-wise normalization is like the idf part of BM25, and the column-wise normalization is like the part manipulating k and b, that try to take documents of different lengths equally. So I only tried BM25WithoutIDF with row-wise normalization and TF.IDF with column-wise normalization.

The benchmarking result shows that BM25 without normalization is still the best one. I think it's because BM25 takes the length of documents into account while VSM doesn't do that well, and [the Google paper](http://infolab.stanford.edu/~backrub/google.html) also claims that VSM tends to rank shorter documents higher.

### Lecture 09 ✅

* Clustering
  * The number of clusters is given as part of the input (k)
  * Goal:
    * Intra-cluster distances are as small as possible 
    * Inter-cluster distances are as large as possible
  * Centroid: intuitively, a centroid is a single element from the metric space that "represents" the cluster
  * RSS: residual sum of squares
* K-Means
  * Idea: find a local optimum of the RSS by greedily minimizing it in every step
  * Steps:
    * Initialization: pick a set of centroids (for example, pick a random subset from the input)
    * Alternate between the following 2 steps
        * A: Assign each element to its nearest centroid
        * B: Compute new centroids as average of elems assigned to it
  * Termination conditions, options
    * when no more change in clustering (optimal, but this can take a very long time)
    * after a fixed number of iterations (easy, but how to guess the right number)
    * when RSS falls below a given threshold (reasonable, but RSS may never fall below that threshold)
    * when decrease in RSS falls below a given threshold (reasonable, stop when when we are close to convergence)
  * Choice of a good k
    * choose the k with smallest RSS (bad idea, because RSS is minimized for k = n)
    * choose the k with smallest RSS + λ * k (make sense)
  * When is K-Means a good clustering algorithm
    * K-Means tends to produce compact clusters of about equal size
  * Alternatives
    * K-Medoids
    * Fuzzy k-means
    * EM-Algorithm
* K-Means for Text Documents

### Lecture 10 ✅
* Latent Semantic Indexing
  * Motivation: synonyms and polysem, add the missing synonyms to the documents automagically.
  * Goal: 
    * Given a term-document matrix A and k < rank(A)
    * Then find a matrix A' of (column) rank k such that the difference between A' and A is as small as possible.
  * How: SVD (singular value decomposition)
    * A = U * S * V
    * For a given k < rank(A) let
      * U_k = the first k columns of U, column-orthonormal
      * S_k = the upper k x k part of S
      * V_k = the first k rows of V
    * A_k = U_k * S_k * V_k
    * Complexity: Lanczos method has complexity O(k*nnz), k is the rank and nnz means the number of non-zero values in the matrix.
* Computing the SVD
  * EVD
  * EVD => SVD
* Using LSI for better Retrieval
  * Variants
    * Variant 1: work with A_k instead of A, but A_k is a dense matrix, O(m*m*n)
    * Variant 2: work with V_k instead of A
      * map the query to concept space
      * work with v_k instead of A
    * Variant 3: expand the original documents
  * Linear combination with original scores in practice

### Lecture 11 ✅

* Classification
  * Objects/Classes
  * Training/Prediction
  * Difference to K-means
    * no cluster names
    * no learning/training phase
    * soft clustering
  * Evaluation: Precision/Recall/F-measure
* Probability recap
  * Maximum Likelihood Estimation (MLE)
  * Conditional probabilities, Bayes Theorem
* Naive Bayes
  * Theoretial: Assumptions/Learning phase/Prediction
  * Practical: Smoothing/Numerical stability/Linear Algebra (LA)

| Name | Genres, Variant 1 | Genres, Variant 2 | Ratings, Variant 1 | Ratings, Variant 2 | Implemented Refinements |
|------|-------------------|-------------------|--------------------|--------------------|-------------------------|
| benchmark (ck1028) | F = 77.49% | F = 78.46% | F = 47.47% | F = 47.48% | None |
| ZhengHe-MD | F = 76.91% | F = 79.19% | F = 47.49% | F = 47.18% | None |

#### Top 20 Words

##### Genres, Variant 1

Horror  young house find woman family old group man night friends years soon dead life evil people death back horror girl 
Drama   life young family father man mother old son wife story woman years lives time home girl daughter day finds year 
Documentary     life documentary world people story years time family war journey year lives history young women many own day man work 
Comedy  life man young wife family find father time old money friends girl daughter friend e day woman school decides way 
Western gang town sheriff man ranch men father gold killed finds money cattle find brother marshal tom outlaw help son murder

##### Genres, Variant 2

Horror  young house woman family find man night old group years evil friends soon dead life death people back horror town 
Western gang town sheriff man ranch men father gold killed finds cattle marshal brother bill outlaw money help steve find johnny 
Comedy  man life wife young family money time find father old house decides son girl daughter wants finds friend town back 
Drama   life young man old father family son mother wife story girl home year years time woman world lives school day 
Documentary     life world documentary years people story time family war old journey history children young year day lives man women way

##### Ratings, Variant 1

R       life young man find family wife woman father old time finds years world friend back story friends own help way 
PG-13   life man family young world find father time school old story mother years home back finds wife help woman year 
PG      life young world man father old find family time town story help back finds boy year mother way years friends

##### Ratings, Variant 2

PG-13   life family man young world find father old story time mother school years finds son home wife back woman year 
R       life man young find family finds father wife time old woman back world years friend help town way police soon 
PG      life young father man old world find family time story town finds help year back boy way school years friends


## References

* [videos](https://www.youtube.com/playlist?list=PLfgMNKpBVg4V8GtMB7eUrTyvITri8WF7i)
* [course wiki](https://ad-wiki.informatik.uni-freiburg.de/teaching/InformationRetrievalWS1718)
* [lecture svn repo](https://daphne.informatik.uni-freiburg.de/ws1718/InformationRetrieval/svn-public/public/): codes, slides and exercises can be found here.
* Text books
  * [Introduction to Information Retrieval](https://nlp.stanford.edu/IR-book/information-retrieval-book.html)
