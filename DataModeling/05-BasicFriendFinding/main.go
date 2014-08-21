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
		create (_0:User {name:"Bill"})
		create (_1:User {name:"Sara"})
		create (_2:User {name:"Derrick"})
		create (_3:User {name:"Ian"})
		create (_4:User {name:"Jill"})
		create (_5:User {name:"Joe"})
		create _0-[:Knows]->_2
		create _0-[:Knows]->_3
		create _1-[:Knows]->_0
		create _1-[:Knows]->_3
		create _1-[:Knows]->_4
		create _5-[:Knows]->_0
		create _5-[:Knows]->_1
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {
	listUsers()
	println()

	getFriendList("Joe")
	println()

	getFriendRecommendations("Joe")
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

	if len(res) > 0 {
		fmt.Println("Users:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func getFriendList(user string) {
	stmt := `
		MATCH (user:User {name: {userSub}})-[:Knows]->(friend:User)
		RETURN friend.name
		ORDER BY friend.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Friend string `json:"friend.name"`
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

	if len(res) > 0 {
		fmt.Println(user, "friends:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Friend)
		}
	} else {
		fmt.Println("No results found")
	}
}

func getFriendRecommendations(user string) {
	stmt := `
		MATCH (user:User {name: {userSub}})-[:Knows*2..2]->(friend_of_friend:User)
		WHERE NOT (user)-[:Knows]-(friend_of_friend)
		RETURN friend_of_friend.name, COUNT(*) as fcount
		ORDER BY fcount DESC, friend_of_friend.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Friend string `json:"friend_of_friend.name"`
		Count  int    `json:"fcount"`
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

	if len(res) > 0 {
		fmt.Println(user, "friend recommendations:")
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Friend, n.Count)
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
