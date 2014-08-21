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
		create (_0:User {name:"Sara"})
		create (_1:Thing {name:"Cats"})
		create (_2:User {name:"Derrick"})
		create (_3:Thing {name:"Bikes"})
		create (_4:User {name:"Jill"})
		create (_5:User {name:"Joe"})
		create _0-[:favorite]->_1
		create _0-[:favorite]->_3
		create _2-[:favorite]->_1
		create _2-[:favorite]->_3
		create _4-[:favorite]->_3
		create _5-[:favorite]->_1
		create _5-[:favorite]->_3
		create _5-[:friend]->_0
		create _0-[:friend]->_5
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
	listThings()
	println()

	getFavoriteList("Joe")
	println()

	findSameFavoriteUsersOfUser("Joe")
	println()
	findSameFavoriteFriendsOfUser("Joe")
	println()

	findThingFavoriteList("Cats")
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
		fmt.Println("  No results found")
	}
}

func listThings() {
	stmt := `
		MATCH (thing:Thing)
		RETURN thing.name
		ORDER BY thing.name
	`
	res := []struct {
		Name string `json:"thing.name"`
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

	fmt.Println("Things:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func getFavoriteList(user string) {
	stmt := `
		MATCH (user:User {name: {userSub}})-[:favorite]->(thing:Thing)
		RETURN thing.name
		ORDER BY thing.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Favorite string `json:"thing.name"`
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

	fmt.Println(user, "favorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Favorite)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func findSameFavoriteUsersOfUser(user string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing)<-[:favorite]-(other:User)
		WHERE user.name = {userSub}
		RETURN other.name, count(thing) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Other string `json:"other.name"`
		Count int    `json:"ocount"`
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

	fmt.Println(user, " users    #cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func findSameFavoriteFriendsOfUser(user string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing)<-[:favorite]-(other:User),
			(user)-[:friend]->(other)
		WHERE user.name = {userSub}
		RETURN other.name, count(thing) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Other string `json:"other.name"`
		Count int    `json:"ocount"`
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

	fmt.Println(user, "friends    #cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func findThingFavoriteList(thing string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(thing:Thing {name: {thingSub}})
		RETURN user.name
		ORDER BY user.name
	`
	params := neoism.Props{"thingSub": thing}

	res := []struct {
		Name string `json:"user.name"`
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

	fmt.Println(thing, "favorited by:", len(res))
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("  No results found")
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
