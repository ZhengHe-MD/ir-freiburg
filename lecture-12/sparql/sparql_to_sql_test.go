package main

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
	"unicode"
)

func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func TestSparQL_SparQLToSQL(t *testing.T) {
	tests := []struct {
		givenSparQL string
		wantSQL     string
	}{
		{
			`
				SELECT ?f ?p WHERE {
					?f "instance of" "film" .
					?f "award received" "Academy Award for Best Picture" .
					?f "director" ?p .
					?p "country of citizenship" "Germany"
				}
			`,
			`
				SELECT f0.subject,
                       f2.object
  				FROM   wikidata as f0,
 					   wikidata as f1,
 					   wikidata as f2,
  					   wikidata as f3
  				WHERE  f0.predicate="instance of"
  					   AND f0.object="film"
  					   AND f1.predicate="award received"
  					   AND f1.object="Academy Award for Best Picture"
  					   AND f2.predicate="director"
  					   AND f3.predicate="country of citizenship"
  					   AND f3.object="Germany"
  					   AND f2.subject=f0.subject
  					   AND f2.subject=f1.subject
					   AND f3.subject=f2.object;
            `,
		},
	}

	engine := &SparQL{}

	for _, tt := range tests {
		sql, err := engine.SparQLToSQL(tt.givenSparQL)
		assert.NoError(t, err)
		assert.Equal(t, SpaceStringsBuilder(tt.wantSQL), SpaceStringsBuilder(sql))
	}
}

func TestSparQL_ProcessSQLQuery(t *testing.T) {
	tests := []struct {
		givenDBName     string
		givenSparQLStmt string
		wantData        [][]string
	}{
		{
			"./example.db",
			`
				SELECT ?f ?p WHERE {
					?f "instance of" "film" .
					?f "award received" "Academy Award for Best Picture" .
					?f "director" ?p .
					?p "country of citizenship" "Germany"
				}	
			`,
			[][]string{
				{"The Life of Emile Zola", "William Dieterle"},
			},
		},
	}

	engine := &SparQL{}
	for _, tt := range tests {
		sqlStmt, err := engine.SparQLToSQL(tt.givenSparQLStmt)
		assert.NoError(t, err)
		data, err := engine.ProcessSQLQuery(tt.givenDBName, sqlStmt)
		assert.NoError(t, err)
		for i := 0; i < len(tt.wantData); i ++ {
			sort.Strings(tt.wantData[i])
		}
		for i := 0; i < len(data); i++ {
			sort.Strings(data[i])
		}
		assert.Equal(t, tt.wantData, data)
	}
}
