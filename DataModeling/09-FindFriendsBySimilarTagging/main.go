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
		CREATE (animal:Tag {name:"Animals"})
		CREATE (hobby:Tag {name:"Hobby"})
		CREATE (bikes:Stuff {name:"Bikes"})
		CREATE (surfing:Stuff {name:"Surfing"})
		CREATE (cats:Stuff {name:"Cats"})
		CREATE (dogs:Stuff {name:"Dogs"})
		CREATE (horses:Stuff {name:"Horses"})
		CREATE (unicorns:Stuff {name:"Unicorns"})
		CREATE (sara:User {name:"Sara"})
		CREATE (derrick:User {name:"Derrick"})
		CREATE (joe:User {name:"Joe"})
		CREATE dogs-[:tagged]->animal
		CREATE cats-[:tagged]->animal
		CREATE horses-[:tagged]->animal
		CREATE unicorns-[:tagged]->animal
		CREATE surfing-[:tagged]->hobby
		CREATE bikes-[:tagged]->hobby
		CREATE sara-[:favorite]->bikes
		CREATE sara-[:favorite]->horses
		CREATE sara-[:favorite]->unicorns
		CREATE derrick-[:favorite]->cats
		CREATE derrick-[:favorite]->bikes
		CREATE joe-[:favorite]->cats
		CREATE joe-[:favorite]->dogs
		CREATE joe-[:favorite]->surfing
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
	listStuff()
	println()

	listUserFavorites("Sara")
	listUserFavorites("Derrick")
	listUserFavorites("Joe")
	println()

	listStuffTags("Surfing")
	listStuffTags("Cats")
	println()

	listSimilarTaggings("Sara")
	listSimilarTaggings("Derrick")
	listSimilarTaggings("Joe")

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
		fmt.Println("No results found")
	}
}

func listUserFavorites(user string) {
	stmt := `
		MATCH (user:User)-[:favorite]->(stuff:Stuff)
		WHERE user.name = {userSub}
		RETURN stuff.name
		ORDER BY stuff.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Stuff string `json:"stuff.name"`
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

	fmt.Println(user, " Favorites:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Stuff)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listStuffTags(item string) {
	stmt := `
		MATCH (item:Stuff)-[:tagged]->(tag:Tag)
		WHERE item.name = {itemSub}
		RETURN tag.name
		ORDER BY tag.name
	`

	params := neoism.Props{"itemSub": item}

	res := []struct {
		Tag string `json:"tag.name"`
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

	fmt.Println(item, " Tags:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Tag)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listSimilarTaggings(user string) {
	stmt := `
		MATCH (me)-[:favorite]->(myFavorites)-[:tagged]->(tag)<-[:tagged]-(theirFavorites)<-[:favorite]-(people)
		WHERE me.name = {userSub} AND NOT me=people
		RETURN
			people.name AS name,
			collect(myFavorites.name) AS myFavs,
			collect(theirFavorites.name) AS theirFavs,
			count(*) AS fcount
		ORDER BY fcount DESC
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Name      string   `json:"name"`
		Fcount    int      `json:"fcount"`
		MyFavs    []string `json:"myFavs"`
		TheirFavs []string `json:"theirFavs"`
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

	fmt.Println("Mutuals of", user, ":")
	if len(res) > 0 {
		fmt.Printf("%-12s  %-24s\n", "Name", "Similar Favorites")
		for _, n := range res {
			fmt.Printf("  %-10s  %d  %v  %v\n", n.Name, n.Fcount, n.MyFavs, n.TheirFavs)
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
