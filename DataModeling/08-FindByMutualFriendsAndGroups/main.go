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
		create (_1:User {name:"Bob"})
		create (_2:User {name:"Jill"})
		create (_3:User {name:"Joe"})
		create (_4:User {name:"Alice"})
		create (_5:Group {name:"Group1"})
		create (_6:Group {name:"Group2"})
		create _0-[:member_of_group]->_5
		create _2-[:member_of_group]->_5
		create _3-[:member_of_group]->_5
		create _4-[:member_of_group]->_5
		create _1-[:member_of_group]->_6
		create _3-[:member_of_group]->_6
		create _4-[:member_of_group]->_6
		create _1-[:knows]->_0
		create _2-[:knows]->_0
		create _3-[:knows]->_0
		create _3-[:knows]->_1
		create _0-[:knows]->_1
		create _0-[:knows]->_2
		create _0-[:knows]->_3
		create _0-[:knows]->_4
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)
}

func main() {

	listUsers()
	listGroups()
	println()

	listGroupMembers("Group1")
	listGroupMembers("Group2")
	println()

	listUserGroups("Bob")
	listUserGroups("Alice")
	println()

	listUserFriends("Bill")
	listUserFriends("Joe")
	println()

	findMutualsForUser("Joe", "Jill")
	findMutualsForUser("Joe", "Bob")
	findMutualsForUser("Joe", "Alice")

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

func listGroups() {
	stmt := `
		MATCH (group:Group)
		RETURN group.name
		ORDER BY group.name
	`
	res := []struct {
		Name string `json:"group.name"`
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

	fmt.Println("Groups:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listGroupMembers(group string) {
	stmt := `
		MATCH (user:User)-[:member_of_group]->(group:Group)
		WHERE group.name = {groupSub}
		RETURN user.name
		ORDER BY user.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		User string `json:"user.name"`
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

	fmt.Println(group, " Users:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.User)
		}
	} else {
		fmt.Println("No results found")
	}

}

func listUserGroups(user string) {
	stmt := `
		MATCH (user:User)-[:member_of_group]->(group:Group)
		WHERE user.name = {userSub}
		RETURN group.name
		ORDER BY group.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Group string `json:"group.name"`
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

	fmt.Println(user, " Groups:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Group)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listUserFriends(user string) {
	stmt := `
		MATCH (user:User)-[:knows]->(friend:User)
		WHERE user.name = {userSub}
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

	fmt.Println(user, " Friends:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Friend)
		}
	} else {
		fmt.Println("No results found")
	}
}

func findMutualsForUser(user string, other string) {
	stmt := `
		MATCH (me:User { name: {userSub} }),(other:User { name: {otherSub} })
		OPTIONAL MATCH pGroups=(me)-[:member_of_group]->(mg)<-[:member_of_group]-(other)
		OPTIONAL MATCH pMutualFriends=(me)-[:knows]->(mf)<-[:knows]-(other)
		RETURN other.name AS name, collect(mg.name) AS mutualGroups,
		  collect(DISTINCT mf.name) AS mutualFriends
		ORDER BY length(collect(DISTINCT mf.name)) DESC
	`

	// othersStr := "['" + others[0]
	// for i, o := range others {
	// 	if i < 1 {
	// 		continue
	// 	}
	// 	othersStr += "', '" + o
	// }
	// othersStr += "']"
	// fmt.Println("others: ", othersStr)

	// ERROR  othersStr doesn't substitute well
	// WHERE other.name IN {othersSub} --> ['Jill', 'Bob']
	// params := neoism.Props{"userSub": user, "othersSub": othersStr}
	params := neoism.Props{"userSub": user, "otherSub": other}

	res := []struct {
		Name    string   `json:"name"`
		Groups  []string `json:"mutualGroups"`
		Friends []string `json:"mutualFriends"`
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
		fmt.Printf("%-12s  %-14s  %-14s\n", "Name", "Groups", "Friends")
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Name, n.Groups, n.Friends)
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
