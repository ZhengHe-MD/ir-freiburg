// Copyright 2016, University of Freiburg,
// Chair of Algorithms and Data Structures.
// Authors: Hannah Bast <bast@cs.uni-freiburg.de>,
//          Patrick Brosi <brosi@cs.uni-freiburg.de>

// NOTE: this file contains specifications and design suggestions in
// pseudo-code. It is not supposed to be compilable in any language. The
// specifications are mandatory, the design suggestions are not.

// Class for a simple inverted index. Copy your code from sheet-01 and extend it
// by BM25 scores, following the explanations in the lecture. In particular,
// consider the implementation advice on slide 18.
class InvertedIndex

  // Build the inverted index from a file with BM25 scores.
  // TEST CASE:
  //   read_from_file("example.txt", 1.75, 0.75)
  // RESULT:
  //   {'docum': [(1, 0.000), (2, 0.000), (3, 0.000)], 'first': [(1, 1.885)], 'second': [(2, 2.325)], 'third': [(3, 2.521)]}
  read_from_file(String file_name, bm25_k, bm25_b)

  // Process a query.
  //
  // !! IMPORTANT: You have to implement ONE of the two following test cases.
  //    You do NOT have to implement both! !!
  //
  // TEST CASE 1:
  //   InvertedIndex ii
  //   ii.inverted_lists = {"bla": [(1, 0.2), (3, 0.6)], "blubb": [(2, 0.4), (3, 0.1), (4, 0.8)]}
  //   process_query("bla blubb")
  // RESULT:
  //   [(4, 0.800), (3, 0.700), (2, 0.400), (1, 0.200)]
  //
  // TEST CASE 2:
  //   read_from_file("example.txt", 1.75, 0.75)
  //   process_query("docum third")
  // RESULT:
  //   [(3, 2.521), (1, 0.000), (2, 0.000)]
  Array<Tuple<int, double>> process_query(String query)

  // Merge two inverted lists.
  // TEST CASE:
  //   merge([(2, 0), (5, 2), (7, 7), (8, 6)], [(4, 1), (5, 3), (6, 3), (8, 3), (9, 8)])
  // RESULT:
  //   [(2, 0), (4, 1), (5, 5), (6, 3), (7, 7), (8, 9), (9, 8)]
  Array<Tuple<int, double>> merge(Array<Tuple<int, double>> a, Array<Tuple<int, double>> b)

// Class for evaluating a given benchmark.
class EvaluateBenchmark

  // Compute the P@k for a given result list and a given set of relevant docs.
  // See the explanation in the lecture, i.p. the example on slide 21.
  // Note that this can also be used to compute P@R by simply taking k as the
  // number of relevant ids.
  // TEST CASE:
  //   precision_at_k([0, 1, 2, 5, 6], [0, 2, 5, 6, 7, 8], 4)
  // RESULT:
  //   0.750
  double precision_at_k(Array<int> results_ids, Set relevant_ids, int k)

  // Compute the AP (avergae precision) of a given result list and a given set
  // of relevant docs. See the explanation in the lecture, i.p. the example on
  // slide 22.
  // TEST CASE:
  //   average_precision([582, 17, 5666, 10003, 10], [10, 582, 877, 10003])
  // RESULT:
  //   0.525
  double average_precision(Array<int> results_ids, Set relevant_ids)

  // Evaluate the given benchmark (in a file in the format of the one from the
  // Wiki) and return the requested measures (MP@3, MP@R, MAP).
  Tripe<double, double, double> evaluate_benchmark(String file_name)