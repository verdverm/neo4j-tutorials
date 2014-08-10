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
		CREATE (himym:TVShow { name: "How I Met Your Mother" })
		CREATE (himym_s1:Season { name: "HIMYM Season 1" })
		CREATE (himym_s1_e1:Episode { name: "Pilot" })
		CREATE (ted:Character { name: "Ted Mosby" })
		CREATE (joshRadnor:Actor { name: "Josh Radnor" })
		CREATE UNIQUE (joshRadnor)-[:PLAYED_CHARACTER]->(ted)
		CREATE UNIQUE (himym)-[:HAS_SEASON]->(himym_s1)
		CREATE UNIQUE (himym_s1)-[:HAS_EPISODE]->(himym_s1_e1)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(ted)
		CREATE (himym_s1_e1_review1 { title: "Meet Me At The Bar In 15 Minutes & Suit Up",
		  content: "It was awesome" })
		CREATE (wakenPayne:User { name: "WakenPayne" })
		CREATE (wakenPayne)-[:WROTE_REVIEW]->(himym_s1_e1_review1)<-[:HAS_REVIEW]-(himym_s1_e1)
		MATCH (himym:TVShow { name: "How I Met Your Mother" }),(himym_s1:Season),
		  (himym_s1_e1:Episode { name: "Pilot" }),
		  (himym)-[:HAS_SEASON]->(himym_s1)-[:HAS_EPISODE]->(himym_s1_e1)
		CREATE (marshall:Character { name: "Marshall Eriksen" })
		CREATE (robin:Character { name: "Robin Scherbatsky" })
		CREATE (barney:Character { name: "Barney Stinson" })
		CREATE (lily:Character { name: "Lily Aldrin" })
		CREATE (jasonSegel:Actor { name: "Jason Segel" })
		CREATE (cobieSmulders:Actor { name: "Cobie Smulders" })
		CREATE (neilPatrickHarris:Actor { name: "Neil Patrick Harris" })
		CREATE (alysonHannigan:Actor { name: "Alyson Hannigan" })
		CREATE UNIQUE (jasonSegel)-[:PLAYED_CHARACTER]->(marshall)
		CREATE UNIQUE (cobieSmulders)-[:PLAYED_CHARACTER]->(robin)
		CREATE UNIQUE (neilPatrickHarris)-[:PLAYED_CHARACTER]->(barney)
		CREATE UNIQUE (alysonHannigan)-[:PLAYED_CHARACTER]->(lily)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(marshall)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(robin)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(barney)
		CREATE UNIQUE (himym_s1_e1)-[:FEATURED_CHARACTER]->(lily)
		CREATE (himym_s1_e1_review2 { title: "What a great pilot for a show :)",
		  content: "The humour is great." })
		CREATE (atlasredux:User { name: "atlasredux" })
		CREATE (atlasredux)-[:WROTE_REVIEW]->(himym_s1_e1_review2)<-[:HAS_REVIEW]-(himym_s1_e1)
		CREATE (er:TVShow { name: "ER" })
		CREATE (er_s7:Season { name: "ER S7" })
		CREATE (er_s7_e17:Episode { name: "Peter's Progress" })
		CREATE (tedMosby:Character { name: "The Advocate " })
		CREATE UNIQUE (er)-[:HAS_SEASON]->(er_s7)
		CREATE UNIQUE (er_s7)-[:HAS_EPISODE]->(er_s7_e17)
		WITH er_s7_e17
		MATCH (actor:Actor),(episode:Episode)
		WHERE actor.name = "Josh Radnor" AND episode.name = "Peter's Progress"
		WITH actor, episode
		CREATE (keith:Character { name: "Keith" })
		CREATE UNIQUE (actor)-[:PLAYED_CHARACTER]->(keith)
		CREATE UNIQUE (episode)-[:FEATURED_CHARACTER]->(keith)
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {

	listGraphData()
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
