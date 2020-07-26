# download data: Because the wikidata-entities dataset is too large,
# push that to github is a waste of space, please execute this line
# in the data directory if you don't have "wikidata-entities.tsv" locally
sh download.sh

# demo
go run cmd/demo/main.go ../data/wikidata-entities.tsv