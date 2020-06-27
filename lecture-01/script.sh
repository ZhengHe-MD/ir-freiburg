# generate words+frequencies.txt
go run inverted_index.go movies.txt | sort -k2,2rc > words+frequencies.txt

# plot, x = docId, y = wordFreq
gnuplot -e "plot 'words+frequencies.txt' using 2; pause -1;"

# plot, x = log(docId), y = log(wordFreq)
gnuplot -e "set logscale xy; plot 'words+frequencies.txt' using 2; pause -1;"