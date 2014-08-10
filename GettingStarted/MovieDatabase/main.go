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

func resetDB() {
	reset.RemoveNeo4jDB()
	reset.StartNeo4jDB()
}

func initDB() {
	stmt := `
		CREATE (matrix1:Movie { title : 'The Matrix', year : '1999-03-31' })
		CREATE (matrix2:Movie { title : 'The Matrix Reloaded', year : '2003-05-07' })
		CREATE (matrix3:Movie { title : 'The Matrix Revolutions', year : '2003-10-27' })
		CREATE (keanu:Actor { name:'Keanu Reeves' })
		CREATE (laurence:Actor { name:'Laurence Fishburne' })
		CREATE (carrieanne:Actor { name:'Carrie-Anne Moss' })
		CREATE (keanu)-[:ACTS_IN { role : 'Neo' }]->(matrix1)
		CREATE (keanu)-[:ACTS_IN { role : 'Neo' }]->(matrix2)
		CREATE (keanu)-[:ACTS_IN { role : 'Neo' }]->(matrix3)
		CREATE (laurence)-[:ACTS_IN { role : 'Morpheus' }]->(matrix1)
		CREATE (laurence)-[:ACTS_IN { role : 'Morpheus' }]->(matrix2)
		CREATE (laurence)-[:ACTS_IN { role : 'Morpheus' }]->(matrix3)
		CREATE (carrieanne)-[:ACTS_IN { role : 'Trinity' }]->(matrix1)
		CREATE (carrieanne)-[:ACTS_IN { role : 'Trinity' }]->(matrix2)
		CREATE (carrieanne)-[:ACTS_IN { role : 'Trinity' }]->(matrix3)
	`

	// construct query
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)
}

func init() {
	// resetDB()
	var err error
	db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	// initDB()
}

func main() {
	countNodes()

	showAllActors()
	getActorByName("Laurence Fishburne")

	showAllMovies()
	getMovieByName("The Matrix")

	listGraphData()
}

// Create a node with neoism function
func countNodes() {
	res := queryNodes("(n)", "", "", "n", "")
	fmt.Println("countNodes()", len(res))
}

func countNodesByType(typ string) {
	match := "(n:" + typ + ")"
	res := queryNodes(match, "", "", "n", "")
	fmt.Println("countNodes()", len(res))
}

func showAllActors() {
	res := queryNodes("(n:Actor)", "", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
}

func getActorByName(name string) {
	res := queryNodes("(n:Actor)", "n.name = '"+name+"'", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
}

func showAllMovies() {
	res := queryNodes("(n:Movie)", "", "", "n", "")
	fmt.Println("Movies: ", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Movie[%d] %+v\n", i, n.Data)
	}
}

func getMovieByName(title string) {
	res := queryNodes("(n:Movie)", "n.title = '"+title+"'", "", "n", "")
	fmt.Println("Actors: ", len(res))
	for i, _ := range res {
		n := res[i].N // Only one row of data returned
		fmt.Printf("  Actor[%d] %+v\n", i, n.Data)
	}
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
		n := res[i] // Only one row of data returned
		fmt.Printf("  [%d] %+v -> %+v -> %+v\n", i, n.From.Data, n.Rel.Data, n.To.Data)
	}
}

// careful...
func queryNodes(MATCH, WHERE, CREATE, RETURN, ORDERBY string) []struct{ N neoism.Node } {
	stmt := ""
	if MATCH != "" {
		stmt += "MATCH " + MATCH + " "
	}
	if WHERE != "" {
		stmt += "WHERE " + WHERE + " "
	}
	if CREATE != "" {
		stmt += "CREATE " + CREATE + " "
	}
	if RETURN != "" {
		stmt += "RETURN " + RETURN + " "
	}
	if ORDERBY != "" {
		stmt += "ORDERBY " + ORDERBY + " "
	}
	stmt += ";"
	// params
	params := neoism.Props{
		"MATCH":   MATCH,
		"WHERE":   WHERE,
		"CREATE":  CREATE,
		"RETURN":  RETURN,
		"ORDERBY": ORDERBY,
	}

	// query results
	res := []struct {
		N neoism.Node // Column "n" gets automagically unmarshalled into field N
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	return res
}
