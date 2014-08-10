package main

import (
	"fmt"
	"github.com/jmcvetta/neoism"
	"github.com/verdverm/neo4j-tutorials/common/reset"
)

var (
	db *neoism.Database
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	resetDB()
	var err error
	db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	initDB()
}

func resetDB() {
	reset.RemoveNeo4jDB()
	reset.StartNeo4jDB()
}

func initDB() {
	stmt := `
		CREATE (root { name: 'ROOT' })-[:LINK]->(root)
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {

	insertNode(25)
	insertNode(10)
	insertNode(14)
	insertNode(30)
	insertNode(35)
	insertNode(20)
	insertNode(05)

	listListData()

	deleteNode(05)
	deleteNode(35)
	deleteNode(20)

	listListData()

}

func insertNode(val int) {
	stmt := `
		MATCH (root)-[:LINK*0..]->(before),(after)-[:LINK*0..]->(root),(before)-[old:LINK]->(after)
		WHERE root.name = 'ROOT' AND (before.value <  {valSub} OR before = root) AND ( {valSub} < after.value OR after = root)
		CREATE UNIQUE (before)-[:LINK]->({ value: {valSub} })-[:LINK]->(after)
		DELETE old
	`
	params := neoism.Props{"valSub": val}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)
}

func deleteNode(val int) {
	stmt := `
		MATCH
		  (root)-[:LINK*0..]->(before),
		  (before)-[delBefore:LINK]->(del)-[delAfter:LINK]->(after),
		  (after)-[:LINK*0..]->(root)
		WHERE root.name = 'ROOT' AND del.value = {valSub}
		CREATE UNIQUE (before)-[:LINK]->(after)
		DELETE del, delBefore, delAfter
	`
	params := neoism.Props{"valSub": val}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

}

func listListData() {
	// query results
	res := []struct {
		From neoism.Node
		Rel  neoism.Relationship
		To   neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: `
			MATCH (n)-[r]->(m)
			RETURN n AS From, r AS Rel, m AS To
			ORDER BY n.value
		`,
		Result: &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("List Nodes: ", len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("%02d: %v\n", i, n.From.Data["value"])
	}
	fmt.Println()
}

func listGraphData() {
	// query results
	res := []struct {
		From neoism.Node
		Rel  neoism.Relationship
		To   neoism.Node
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement: `
			MATCH (n)-[r]->(m)
			RETURN n AS From, r AS Rel, m AS To;
		`,
		Result: &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Graph Data: ", len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %+v -> %+v -> %+v\n", i, n.From.Data, n.Rel.Data, n.To.Data)
	}
}
