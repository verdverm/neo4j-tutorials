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
	clearDB()
	initDB()
}

func clearDB() {
	stmt := `
        MATCH (n)
        OPTIONAL MATCH (n)-[r]-()
        DELETE n,r
    `
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func resetDB() {
	reset.RemoveNeo4jDB()
	reset.StartNeo4jDB()
}

func initDB() {
	stmt := `
		CREATE (cats:Stuff {name:"cats"})
		CREATE (nature:Stuff {name:"nature"})
		CREATE (bikes:Stuff {name:"bikes"})
		CREATE (cars:Stuff {name:"cars"})
		CREATE (ben:User {name:"Ben"})
		CREATE (sara:User {name:"Sara"})
		CREATE (maria:User {name:"Maria"})
		CREATE (joe:User {name:"Joe"})
		CREATE sara-[:FOLLOWS]->joe
		CREATE sara-[:FOLLOWS]->ben
		CREATE sara-[:LIKES]->bikes
		CREATE sara-[:LIKES]->cars
		CREATE sara-[:LIKES]->cats
		CREATE maria-[:FOLLOWS]->joe
		CREATE maria-[:LOVES]->joe
		CREATE maria-[:LIKES]->cars
		CREATE joe-[:FOLLOWS]->sara
		CREATE joe-[:FOLLOWS]->maria
		CREATE joe-[:LOVES]->maria
		CREATE joe-[:LIKES]->bikes
		CREATE joe-[:LIKES]->nature
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {

	listUsers()
	listStuff()
	println()

	listUserRelations("Sara")
	listUserRelations("Ben")
	listUserRelations("Maria")
	listUserRelations("Joe")
	println()

	findLovers()
	println()

	// listGraphData()
}

func listUsers() {
	stmt := `
		MATCH (user:User)
		RETURN user.name
		ORDER BY user.name
	`
	res := []struct {
		Name string `json:"user.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Users:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listStuff() {
	stmt := `
		MATCH (stuff:Stuff)
		RETURN stuff.name
		ORDER BY stuff.name
	`
	res := []struct {
		Name string `json:"stuff.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Stuff:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listUserRelations(user string) {
	stmt := `
		MATCH
			(user:User)-[:FOLLOWS]->(follow),
			(user:User)-[:LIKES]->(like),
			(user:User)-[:LOVES]->(love)
		WHERE user.name = {userSub}
		RETURN
			collect(DISTINCT follow.name) AS Follows,
			collect(DISTINCT like.name) AS Likes,
			collect(DISTINCT love.name) AS Loves
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Follows []string `json:"Follows"`
		Likes   []string `json:"Likes"`
		Loves   []string `json:"Loves"`
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

	fmt.Println(user, " Relations:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  follows: %v\n", n.Follows)
			fmt.Printf("  likes:   %v\n", n.Likes)
			fmt.Printf("  loves:   %v\n", n.Loves)
		}
	} else {
		fmt.Println("No results found")
	}
}

func findLovers() {
	stmt := `
		MATCH
			(u1:User)-[love:LOVES]->(u2:User)-[:LOVES]->(u1)
		WHERE u1.name <= u2.name
		RETURN u1.name, u2.name
		ORDER BY u1.name
	`

	res := []struct {
		LoverA string `json:"u1.name"`
		LoverB string `json:"u2.name"`
	}{}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: nil,
		Result:     &res,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("Lovers:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %s \u2661 %s\n", n.LoverA, n.LoverB)
		}
	} else {
		fmt.Println("No results found")
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
