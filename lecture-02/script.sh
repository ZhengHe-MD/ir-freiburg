# exercise 3: run benchmark on movies dataset
go run cmd/benchmark/main.go ../data/movies.txt ../data/movies-benchmark-minus-1.txt

# without any refinements

# b = 0.75, k = 1.25
# MP@3: 0.444
# MP@R: 0.327
# MAP: 0.318

# b = 0.1, k = 0.75
# MP@3: 0.556
# MP@R: 0.473
# MAP: 0.488

# b = 0.3, k = 1.0
# MP@3: 0.481
# jMP@R: 0.416
# MAP: 0.468

# b = 0.34, k = 1.35
# MP@3: 0.481
# MP@R: 0.415
# MAP: 0.420

# b = 0.11, k = 0.77
# MP@3: 0.593
# MP@R: 0.473
# MAP: 0.491
