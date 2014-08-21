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
		CREATE (_0 {name:"U2G2R34"})
		CREATE (_1 {name:"U1G3R34"})
		CREATE (_2:User {name:"User2"})
		CREATE (_3:User {name:"User1"})
		CREATE (_4:Role {name:"Role6"})
		CREATE (_5 {name:"U1G2R23"})
		CREATE (_6:Role {name:"Role4"})
		CREATE (_7:Role {name:"Role5"})
		CREATE (_8 {name:"U2G1R25"})
		CREATE (_9:Group {name:"Group1"})
		CREATE (_10:Role {name:"Role2"})
		CREATE (_11:Group {name:"Group2"})
		CREATE (_12:Role {name:"Role3"})
		CREATE (_13:Group {name:"Group3"})
		CREATE (_14 {name:"U1G1R12"})
		CREATE (_15:Role {name:"Role1"})
		CREATE (_16 {name:"U2G3R56"})
		CREATE _0-[:hasGroup]->_11
		CREATE _0-[:hasRole]->_12
		CREATE _0-[:hasRole]->_6
		CREATE _1-[:hasGroup]->_13
		CREATE _1-[:hasRole]->_12
		CREATE _1-[:hasRole]->_6
		CREATE _2-[:hasRoleInGroup]->_8
		CREATE _2-[:hasRoleInGroup]->_0
		CREATE _2-[:hasRoleInGroup]->_16
		CREATE _3-[:hasRoleInGroup]->_14
		CREATE _3-[:hasRoleInGroup]->_5
		CREATE _3-[:hasRoleInGroup]->_1
		CREATE _5-[:hasGroup]->_11
		CREATE _5-[:hasRole]->_10
		CREATE _5-[:hasRole]->_12
		CREATE _8-[:hasGroup]->_9
		CREATE _8-[:hasRole]->_10
		CREATE _8-[:hasRole]->_7
		CREATE _14-[:hasGroup]->_9
		CREATE _14-[:hasRole]->_15
		CREATE _14-[:hasRole]->_10
		CREATE _16-[:hasGroup]->_13
		CREATE _16-[:hasRole]->_7
		CREATE _16-[:hasRole]->_4
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

	listGroupRolesUser("Group1")
	listGroupRolesUser("Group2")
	listGroupRolesUser("Group3")
	println()

	listUserGroupRoles("User1")
	listUserGroupRoles("User2")
	println()

	listUserRolesInGroup("User1", "Group1")
	println()

	findCommonGroupsForUsers("User1", "User2")
	println()

	findCommonRolesForUsers("User1", "User2")
	println()

	findCommonGroupRoles()

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

	if len(res) > 0 {
		fmt.Println("Groups:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listGroupRoles(group string) {
	stmt := `
		MATCH (hyperedge)-[:hasGroup]->(group)
		WHERE group.name = {groupSub}
		MATCH (hyperedge)-[:hasRole]->(role)
		RETURN role.name
		ORDER BY role.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		Name string `json:"role.name"`
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
		fmt.Println(group, " Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}

}

func listGroupRolesUser(group string) {
	stmt := `
		MATCH (hyperedge)-[:hasGroup]->(group)
		WHERE group.name = {groupSub}
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)-[:hasRole]->(role)
		WITH user, role
		ORDER BY user.name
		RETURN role.name, collect(user.name) AS users
		ORDER BY role.name
	`

	params := neoism.Props{"groupSub": group}

	res := []struct {
		Role  string   `json:"role.name"`
		Users []string `json:"users"`
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
		fmt.Println(group, " Role / User:")
		for _, n := range res {
			fmt.Printf("  %-24s  %v\n", n.Role, n.Users)
		}
	} else {
		fmt.Println("No results found")
	}

}

func listUserGroupRoles(user string) {
	stmt := `
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)
		WHERE user.name = {userSub}
		MATCH (hyperedge)-[:hasGroup]->(group), (hyperedge)-[:hasRole]->(role)
		RETURN group.name, collect(role.name) AS Roles
		ORDER BY group.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Group string   `json:"group.name"`
		Role  []string `json:"Roles"`
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
		fmt.Println(user, " Groups / Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s  %-v\n", n.Group, n.Role)
		}
	} else {
		fmt.Println("No results found")
	}

}

func listUserRolesInGroup(user, group string) {
	stmt := `
		MATCH (user:User)-[:hasRoleInGroup]->(hyperedge)-[:hasGroup]->(group:Group)
		WHERE user.name = {userSub} AND group.name = {groupSub}
		MATCH (hyperedge)-[:hasRole]->(role)
		RETURN role.name
		ORDER BY role.name
	`

	params := neoism.Props{"userSub": user, "groupSub": group}

	res := []struct {
		Role string `json:"role.name"`
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
		fmt.Println(user, "-", group, " Roles:")
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Role)
		}
	} else {
		fmt.Println("No results found")
	}

}

func findCommonGroupsForUsers(user1, user2 string) {
	stmt := `
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasGroup]->(group:Group),
			(hyper1)-[:hasRole]->(role1),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasGroup]->(group:Group),
			(hyper2)-[:hasRole]->(role2)
		WHERE u1.name = {user1Sub} AND u2.name = {user2Sub}
		RETURN group.name, collect(DISTINCT role1.name) AS role1s, collect(DISTINCT role2.name) AS role2s
		ORDER BY group.name
	`

	params := neoism.Props{"user1Sub": user1, "user2Sub": user2}

	res := []struct {
		Group string   `json:"group.name"`
		Role1 []string `json:"role1s"`
		Role2 []string `json:"role2s"`
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
		fmt.Printf("%-12s  %-14s  %-14s\n", "Group", user1, user2)
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Group, n.Role1, n.Role2)
		}
	} else {
		fmt.Println("No results found")
	}
}

func findCommonRolesForUsers(user1, user2 string) {
	stmt := `
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasRole]->(role:Role),
			(hyper1)-[:hasGroup]->(g1),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasRole]->(role:Role),
			(hyper2)-[:hasGroup]->(g2)
		WHERE u1.name = {user1Sub} AND u2.name = {user2Sub}
		RETURN role.name, collect(DISTINCT g1.name) AS group1s, collect(DISTINCT g2.name) AS group2s
		ORDER BY role.name
	`

	params := neoism.Props{"user1Sub": user1, "user2Sub": user2}

	res := []struct {
		Role   string   `json:"role.name"`
		Group1 []string `json:"group1s"`
		Group2 []string `json:"group2s"`
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
		fmt.Printf("%-12s  %-16s  %-16s\n", "Role", user1, user2)
		for _, n := range res {
			fmt.Printf("  %-10s  %v  %v\n", n.Role, n.Group1, n.Group2)
		}
	} else {
		fmt.Println("No results found")
	}
}

func findCommonGroupRoles() {
	stmt := `
		MATCH
			(user:User)-[:hasRoleInGroup]->(hyper)-[:hasRole]->(role:Role),
			(hyper)-[:hasGroup]->(group:Group)
		WITH user
		ORDER BY user.name
		MATCH
			(u1:User)-[:hasRoleInGroup]->(hyper1)-[:hasRole]->(role:Role),
			(hyper1)-[:hasGroup]->(group:Group),
			(u2:User)-[:hasRoleInGroup]->(hyper2)-[:hasRole]->(role:Role),
			(hyper2)-[:hasGroup]->(group:Group)
		WHERE u1.name <> u2.name
		RETURN group.name, role.name, collect(DISTINCT user.name) AS users
		ORDER BY group.name, role.name


	`

	res := []struct {
		Group string   `json:"group.name"`
		Role  string   `json:"role.name"`
		Users []string `json:"users"`
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
		fmt.Printf("%-12s  %-12s  %-12s\n", "Group", "Role", "Users")
		for _, n := range res {
			fmt.Printf("  %-10s  %-12s  %v\n", n.Group, n.Role, n.Users)
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
