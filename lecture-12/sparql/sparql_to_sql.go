package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ZhengHe-MD/wheels/collection"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var triplePattern = regexp.MustCompile(`\s?(\?[^\s]+|"[^"]+")\s+("[^"]+")\s+(\?[^\s]+|"[^"]+")`)

type Triple struct {
	Subject    string
	Predicate  string
	Object     string
	TableAlias string
}

type SelectClause struct {
	TableAlias string
	Field      string
}

type FromClause struct {
	Source     string
	TableAlias string
}

type WhereClause struct {
	Left  string
	Op    string
	Right string
}

type LimitClause struct {
	Limit int64
}

type VarCondition struct {
	TableAlias  string
	TripleField string
}

const (
	SparQLSubject   = "subject"
	SparQLPredicate = "predicate"
	SparQLObject    = "object"
)

type SparQL struct {
}

func (m *SparQL) SparQLToSQL(sparQL string) (sqlStmt string, err error) {
	// Transform all letters to lower cases.
	sparQLL := strings.ToLower(sparQL)
	// Find all variables in the SparQL between the SELECT and WHERE clause.
	selectStart := strings.Index(sparQLL, "select ") + 7
	selectEnd := strings.Index(sparQLL, " where")
	variables := strings.Fields(sparQLL[selectStart:selectEnd])

	// Find all triples between "WHERE {" and "}"
	whereStart := strings.Index(sparQLL, "{") + 1
	whereEnd := strings.LastIndex(sparQLL, "}")
	whereText := sparQL[whereStart:whereEnd]
	tripleTexts := strings.Split(whereText, ".")
	var triples []Triple
	for _, tripleText := range tripleTexts {
		matches := triplePattern.FindStringSubmatch(strings.TrimSpace(tripleText))
		subj := strings.ReplaceAll(matches[1], "\"", "")
		pred := strings.ReplaceAll(matches[2], "\"", "")
		obj := strings.ReplaceAll(matches[3], "\"", "")
		triples = append(triples, Triple{
			Subject:   subj,
			Predicate: pred,
			Object:    obj,
		})
	}

	// Find the (optional) LIMIT clause
	hasLimit := true
	limitStart := strings.Index(sparQLL, " limit ") + 7
	limitText := strings.TrimSuffix(strings.TrimSpace(sparQLL[limitStart:]), ";")
	limit, err := strconv.ParseInt(limitText, 10, 64)
	if err != nil {
		hasLimit = false
		err = nil
	}

	var selectList []SelectClause
	var fromList []FromClause
	var whereList []WhereClause
	var limitList []LimitClause

	variableSet := collection.StringSliceToSet(variables)
	variableToConditions := make(map[string][]VarCondition)
	for i := 0; i < len(triples); i++ {
		tableAlias := fmt.Sprintf("f%d", i)
		fromList = append(fromList, FromClause{
			Source:     "wikidata",
			TableAlias: tableAlias,
		})
		triples[i].TableAlias = tableAlias

		triple := triples[i]

		if _, ok := variableSet[triple.Subject]; ok {
			variableToConditions[triple.Subject] = append(
				variableToConditions[triple.Subject],
				VarCondition{TableAlias: tableAlias, TripleField: SparQLSubject})
		} else {
			whereList = append(whereList, WhereClause{
				Left:  fmt.Sprintf("%s.%s", tableAlias, SparQLSubject),
				Op:    "=",
				Right: fmt.Sprintf("\"%s\"", triple.Subject),
			})
		}

		if _, ok := variableSet[triple.Predicate]; ok {
			variableToConditions[triple.Predicate] = append(
				variableToConditions[triple.Predicate],
				VarCondition{TableAlias: tableAlias, TripleField: SparQLPredicate})
		} else {
			whereList = append(whereList, WhereClause{
				Left:  fmt.Sprintf("%s.%s", tableAlias, SparQLPredicate),
				Op:    "=",
				Right: fmt.Sprintf("\"%s\"", triple.Predicate),
			})
		}

		if _, ok := variableSet[triple.Object]; ok {
			variableToConditions[triple.Object] = append(
				variableToConditions[triple.Object],
				VarCondition{TableAlias: tableAlias, TripleField: SparQLObject})
		} else {
			whereList = append(whereList, WhereClause{
				Left:  fmt.Sprintf("%s.%s", tableAlias, SparQLObject),
				Op:    "=",
				Right: fmt.Sprintf("\"%s\"", triple.Object),
			})
		}
	}

	for _, varConds := range variableToConditions {
		selectList = append(selectList, SelectClause{
			TableAlias: varConds[0].TableAlias,
			Field:      varConds[0].TripleField,
		})

		lastCond := varConds[len(varConds)-1]
		for i := 0; i < len(varConds)-1; i++ {
			currCond := varConds[i]

			whereList = append(whereList, WhereClause{
				Left:  fmt.Sprintf("%s.%s", lastCond.TableAlias, lastCond.TripleField),
				Op:    "=",
				Right: fmt.Sprintf("%s.%s", currCond.TableAlias, currCond.TripleField),
			})
		}
	}

	if hasLimit {
		limitList = append(limitList, LimitClause{Limit: limit})
	}

	sqlBuilder := strings.Builder{}

	sqlBuilder.WriteString("SELECT")
	// NOTE: make tests stable
	sort.Slice(selectList, func(i, j int) bool {
		return strings.Compare(selectList[i].TableAlias + selectList[i].Field,
			                   selectList[j].TableAlias + selectList[j].Field) == 0
	})
	for i, selectClause := range selectList {
		if i == len(selectList)-1 {
			sqlBuilder.WriteString(fmt.Sprintf("\t%s.%s\n", selectClause.TableAlias, selectClause.Field))
		} else {
			sqlBuilder.WriteString(fmt.Sprintf("\t%s.%s,\n", selectClause.TableAlias, selectClause.Field))
		}
	}

	sqlBuilder.WriteString("FROM")
	// NOTE: make tests stable
	sort.Slice(fromList, func(i, j int) bool {
		return strings.Compare(fromList[i].TableAlias + fromList[i].Source,
			fromList[j].TableAlias + fromList[j].Source) == 0
	})
	for i, fromClause := range fromList {
		if i == len(fromList)-1 {
			sqlBuilder.WriteString(fmt.Sprintf("\t%s as %s\n", fromClause.Source, fromClause.TableAlias))
		} else {
			sqlBuilder.WriteString(fmt.Sprintf("\t%s as %s,\n", fromClause.Source, fromClause.TableAlias))
		}
	}

	sqlBuilder.WriteString("WHERE")
	sort.Slice(whereList, func(i, j int) bool {
		return strings.Compare(whereList[i].Left + whereList[i].Op + whereList[i].Right,
			whereList[j].Left + whereList[j].Op + whereList[j].Right) == 0
	})

	for i, whereClause := range whereList {
		if i == 0 {
			sqlBuilder.WriteString(fmt.Sprintf("\t%s%s%s\n",
				whereClause.Left,
				whereClause.Op,
				whereClause.Right))
		} else {
			sqlBuilder.WriteString(fmt.Sprintf("\tAND %s%s%s\n",
				whereClause.Left,
				whereClause.Op,
				whereClause.Right))
		}
	}

	for _, limitClause := range limitList {
		sqlBuilder.WriteString(fmt.Sprintf("LIMIT %d", limitClause.Limit))
	}

	sqlBuilder.WriteByte(';')
	sqlStmt = sqlBuilder.String()
	return
}

func (m *SparQL) ProcessSQLQuery(dbName string, sqlStmt string) (data [][]string, err error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return
	}
	defer db.Close()

	row, err := db.Query(sqlStmt)
	if err != nil {
		return
	}
	defer row.Close()
	for row.Next() {
		var columns []string
		columns, err = row.Columns()
		if err != nil {
			return
		}

		values := make([]string, len(columns))
		dest := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			dest[i] = &values[i]
		}

		err = row.Scan(dest...)
		if err != nil {
			return
		} else {
			data = append(data, values)
		}
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: cmd <dbname>")
		os.Exit(-1)
	}

	dbName := os.Args[1]

	engine := &SparQL{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter SparQL query: ")

	var multiLineInput []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "EOF" {
			fmt.Println("See you.")
			break
		}

		if line != "" {
			multiLineInput = append(multiLineInput, line)
			continue
		}

		rawSql := strings.Join(multiLineInput, "\n")
		multiLineInput = nil

		sqlStmt, err := engine.SparQLToSQL(rawSql)
		fmt.Println("got sqlStmt:")
		fmt.Println(sqlStmt)
		if err != nil {
			fmt.Println("SparQL syntax error")
			continue
		}

		data, err := engine.ProcessSQLQuery(dbName, sqlStmt)
		if err != nil {
			fmt.Println("process sqlStmt err", err)
			continue
		}
		fmt.Printf("found %d rows\n", len(data))
		bs, _ := json.Marshal(data)
		fmt.Println(string(bs))
		fmt.Println("Enter SparQL query: ")
	}

	os.Exit(0)
}
