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
}

func main() {
	listGraphData()

	createUser("Tony")
	createRating("Tony", "The Matrix", "5", "The first is usually the best")
	createRating("Tony", "The Matrix Reloaded", "4", "Good Action...")
	createRating("Tony", "The Matrix Revolutions", "3", "Getting old......")

	createUser("John")
	createRating("John", "The Matrix", "5", "Awesome!")
	createRating("John", "The Matrix Reloaded", "3", "Neo4j kicks more ass than Neo plain")
	createRating("John", "The Matrix Reloaded", "1", "shit...")

	createUser("Bob")
	createRating("Bob", "The Matrix", "4", "")
	createRating("Bob", "The Matrix Reloaded", "2", "")

	createFriendship("Tony", "John")
	createFriendship("Tony", "Bob")
	createFriendship("John", "Bob")
	createFriendship("Bob", "Tony")

	getRatingsByUser("Tony")
	getFriendsByUser("Tony")

	listGraphData()
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

func createUser(name string) {
	queryNodes("", "", "(n:User {name: '"+name+"'})", "n", "")
}

func createFriendship(user, friend string) {
	match := "(u:User {name: '" + user + "'}),(f:User {name: '" + friend + "'})"
	create := "(u)-[:FRIEND]->(f)"
	queryNodes(match, "", create, "", "")
}

func createRating(user, title, stars, comment string) {
	match := "(u:User {name: '" + user + "'}),(m:Movie {title: '" + title + "'})"
	create := "(u)-[:RATED { stars: " + stars + ", comment: '" + comment + "'}]->(m)"
	queryNodes(match, "", create, "", "")
}

func getRatingsByUser(user string) {
	stmt := `
		MATCH (u:User {name: {userSub}}),(u)-[rating:RATED]->(movie)
	    RETURN movie, rating;
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Movie  neoism.Node
		Rating neoism.Relationship
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

	fmt.Println("User Ratings: ", user, len(res))
	for i, _ := range res {
		m := res[i].Movie.Data
		r := res[i].Rating.Data.(map[string]interface{})
		fmt.Printf("  [%d] %v    %v    %v\n",
			i, m["title"], r["stars"], r["comment"])
	}
}

func getFriendsByUser(user string) {
	stmt := `
		MATCH (u:User {name: {userSub}}),(u)-[r:FRIEND]->(f)
	    RETURN type(r) AS T, f.name AS F;
	`
	params := neoism.Props{"userSub": user}

	// query results
	res := []struct {
		T string
		F string
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

	fmt.Println("User Friends: ", user, len(res))
	for i, _ := range res {
		n := res[i]
		fmt.Printf("  [%d] %q  %q\n", i, n.T, n.F)
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
