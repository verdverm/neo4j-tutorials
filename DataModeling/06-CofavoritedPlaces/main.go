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
		CREATE (_0:Place {name:"SaunaX"})
		CREATE (_1:Place {name:"CoffeeShop1"})
		CREATE (_2:Place {name:"MelsPlace"})
		CREATE (_3:Place {name:"CoffeeShop3"})
		CREATE (_4:Tag {name:"Cool"})
		CREATE (_5:Place {name:"CoffeeShop2"})
		CREATE (_6:Place {name:"CoffeShop2"})
		CREATE (_7:Tag {name:"Cosy"})
		CREATE (_8:User {name:"Jill"})
		CREATE (_9:User {name:"Joe"})
		CREATE _1-[:tagged]->_4
		CREATE _1-[:tagged]->_7
		CREATE _2-[:tagged]->_7
		CREATE _2-[:tagged]->_4
		CREATE _3-[:tagged]->_7
		CREATE _5-[:tagged]->_4
		CREATE _8-[:favorite]->_1
		CREATE _8-[:favorite]->_2
		CREATE _8-[:favorite]->_6
		CREATE _9-[:favorite]->_1
		CREATE _9-[:favorite]->_0
		CREATE _9-[:favorite]->_2
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {
	listUsers()
	listTags()
	println()
	listPlaces()
	println()

	getFavoriteList("Jill")
	getFavoriteList("Joe")
	println()

	cofavoritedPlaces("CoffeeShop1")
	cotaggedPlaces("CoffeeShop1")

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

func listPlaces() {
	stmt := `
		MATCH (place:Place)
		RETURN place.name
		ORDER BY place.name
	`
	res := []struct {
		Name string `json:"place.name"`
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

	fmt.Println("Places:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func listTags() {
	stmt := `
		MATCH (tag:Tag)
		RETURN tag.name
		ORDER BY tag.name
	`
	res := []struct {
		Name string `json:"tag.name"`
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

	fmt.Println("Tags:")
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
		MATCH (user:User {name: {userSub}})-[:favorite]->(place:Place)
		RETURN place.name
		ORDER BY place.name
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Favorite string `json:"place.name"`
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

func cofavoritedPlaces(place string) {
	stmt := `
		MATCH (place:Place)<-[:favorite]-(person:User)-[:favorite]->(other:Place)
		WHERE place.name = {placeSub}
		RETURN other.name, count(*) AS ocount
		ORDER BY ocount DESC, other.name
	`
	params := neoism.Props{"placeSub": place}

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

	fmt.Println(place, "cofavorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %4d\n", n.Other, n.Count)
		}
	} else {
		fmt.Println("  No results found")
	}
}

func cotaggedPlaces(place string) {
	stmt := `
		MATCH (place:Place)-[:tagged]->(tag:Tag)<-[:tagged]-(other:Place)
		WHERE place.name = {placeSub}
		RETURN other.name, collect(tag.name) as tags
		ORDER BY length(collect(tag.name)) DESC, other.name
	`
	params := neoism.Props{"placeSub": place}

	res := []struct {
		Other string   `json:"other.name"`
		Tags  []string `json:"tags"`
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

	fmt.Println(place, "cotagged:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s  %q\n", n.Other, n.Tags)
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
