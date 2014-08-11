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
	// resetDB()
	var err error
	db, err = neoism.Connect("http://localhost:7474/db/data")
	if err != nil {
		panic(err)
	}
	// initDB()
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
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)

	listGraphData()
}

func main() {
	getOtherMoviesViaActors("The Matrix")
	getCoActingFromMovie("The Matrix")
	getActorPaths("Keanu Reeves", "Carrie-Anne Moss")
}

func getOtherMoviesViaActors(movie string) {
	stmt := `
		MATCH (:Movie { title: {movieSub} })<-[:ACTS_IN]-(actor)-[:ACTS_IN]->(movie)
		RETURN movie.title AS Title, collect(actor.name) AS Actors, count(*) AS Count
		ORDER BY count(*) DESC ;
	`
	params := neoism.Props{"movieSub": movie}

	// query results
	res := []struct {
		Title  string
		Actors []string
		Count  int
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

	fmt.Println("Movies: ", movie, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %24q  %d  %v\n", i, n.Title, n.Count, n.Actors)
	}
}

func getCoActingFromMovie(movie string) {
	stmt := `
		MATCH (:Movie { title: {movieSub} })<-[:ACTS_IN]-(actor)-[:ACTS_IN]->(movie)<-[:ACTS_IN]-(colleague)
		RETURN actor.name AS Actor, collect(DISTINCT colleague.name) AS Actors;
	`
	params := neoism.Props{"movieSub": movie}

	// query results
	res := []struct {
		Actor  string
		Actors []string
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

	fmt.Println("Movies: ", movie, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %24q  %v\n", i, n.Actor, n.Actors)
	}
}

func getActorPaths(actor1, actor2 string) {
	stmt := `
		MATCH p =(:Actor { name: {actor1Sub} })-[:ACTS_IN*0..5]-(:Actor { name: {actor2Sub} })
		RETURN extract(n IN nodes(p)| coalesce(n.title,n.name)) AS List, length(p) AS Len
		ORDER BY length(p)
		LIMIT 10;
	`
	params := neoism.Props{"actor1Sub": actor1, "actor2Sub": actor2}

	// query results
	res := []struct {
		List []string
		Len  int
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

	fmt.Println("Paths: ", actor1, actor2)
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %d  %v\n", i, n.Len, n.List)
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
		n := res[i]
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
		N neoism.Node
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
