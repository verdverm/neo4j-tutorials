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
	// stmt := `
	// 	CREATE (_0 { name: "U1G2R1"})
	// 	CREATE (_1 { name: "Role2"})
	// 	CREATE (_2 { name: "Group1"})
	// 	CREATE (_3 { name: "Group2"})
	// 	CREATE (_4 { name: "Role1"})
	// 	CREATE (_5 { name: "Role"})
	// 	CREATE (_6 { name: "User1"})
	// 	CREATE (_7 { name: "U1G1R2"})
	// 	CREATE (_8 { name: "Group"})
	// 	CREATE _0-[:hasRole]->_4
	// 	CREATE _0-[:hasGroup]->_3
	// 	CREATE _1-[:isA]->_5
	// 	CREATE _2-[:canHave]->_4
	// 	CREATE _2-[:canHave]->_1
	// 	CREATE _2-[:isA]->_8
	// 	CREATE _3-[:canHave]->_1
	// 	CREATE _3-[:canHave]->_4
	// 	CREATE _3-[:isA]->_8
	// 	CREATE _4-[:isA]->_5
	// 	CREATE _6-[:in]->_2
	// 	CREATE _6-[:in]->_3
	// 	CREATE _6-[:hasRoleInGroup]->_0
	// 	CREATE _6-[:hasRoleInGroup]->_7
	// 	CREATE _7-[:hasRole]->_1
	// 	CREATE _7-[:hasGroup]->_2
	// `
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
	listGroupRoles("Group1")
	listGroupRoles("Group2")
	listGroupRoles("Group3")
	// listUserGroupsRoles("User1")

	// listGraphData()
}

// func getShowInfo(show string) {
// 	stmt := `
// 		MATCH (tvShow:TVShow)-[:HAS_SEASON]->(season)-[:HAS_EPISODE]->(episode)
// 		WHERE tvShow.name = {showSub}
// 		RETURN season.name, episode.name
// 	`
// 	params := neoism.Props{"showSub": show}

// 	res := []struct {
// 		Season  string `json:"season.name"`
// 		Episode string `json:"episode.name"`
// 	}{}

// 	// construct query
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}

// 	// execute query
// 	err := db.Cypher(&cq)
// 	panicErr(err)

// 	if len(res) > 0 {
// 		fmt.Println("Show: ", show)
// 		fmt.Println("  ", res[0].Season, res[0].Episode)
// 	} else {
// 		fmt.Println("No results found")
// 	}

// }

// func getShowInfoWithComments(show string) {
// 	stmt := `
// 		MATCH (tvShow:TVShow)-[:HAS_SEASON]->(season)-[:HAS_EPISODE]->(episode)
// 		WHERE tvShow.name = {showSub}
// 		WITH season, episode
// 		OPTIONAL MATCH (episode)-[:HAS_REVIEW]->(review)
// 		RETURN season.name, episode.name, collect(review.content) AS Reviews
// 	`
// 	params := neoism.Props{"showSub": show}

// 	res := []struct {
// 		Season  string `json:"season.name"`
// 		Episode string `json:"episode.name"`
// 		Reviews []string
// 	}{}

// 	// construct query
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}

// 	// execute query
// 	err := db.Cypher(&cq)
// 	panicErr(err)

// 	if len(res) > 0 {

// 		fmt.Println("Show & Reviews: ", show, len(res[0].Reviews))
// 		fmt.Println("Show: ", show)
// 		fmt.Println("  ", res[0].Season, res[0].Episode)
// 		for i, r := range res[0].Reviews {
// 			fmt.Printf("     %d:  %s\n", i, r)
// 		}
// 	} else {
// 		fmt.Println("No results found")
// 	}
// }

// func getCharacterList(show string) {
// 	stmt := `
// 		MATCH (tvShow:TVShow)-[:HAS_SEASON]->()-[:HAS_EPISODE]->()-[:FEATURED_CHARACTER]->(character)
// 		WHERE tvShow.name = {showSub}
// 		RETURN DISTINCT character.name
// 	`
// 	params := neoism.Props{"showSub": show}

// 	res := []struct {
// 		Name string `json:"character.name"`
// 	}{}

// 	// construct query
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}

// 	// execute query
// 	err := db.Cypher(&cq)
// 	panicErr(err)

// 	if len(res) > 0 {
// 		fmt.Println("Show: ", show)
// 		for _, n := range res {
// 			fmt.Println(" ", n.Name)
// 		}
// 	} else {
// 		fmt.Println("No results found")
// 	}

// }

// func getActorList(show string) {
// 	stmt := `
// 		MATCH (tvShow:TVShow)-[:HAS_SEASON]->()-[:HAS_EPISODE]->()-[:FEATURED_CHARACTER]->(character)<-[:PLAYED_CHARACTER]-(actor)
// 		WHERE tvShow.name = {showSub}
// 		RETURN DISTINCT actor.name, character.name
// 	`
// 	params := neoism.Props{"showSub": show}

// 	res := []struct {
// 		Name string `json:"actor.name"`
// 		Char string `json:"character.name"`
// 	}{}

// 	// construct query
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}

// 	// execute query
// 	err := db.Cypher(&cq)
// 	panicErr(err)

// 	if len(res) > 0 {
// 		fmt.Println("Show: ", show)
// 		for _, n := range res {
// 			fmt.Printf("  %-24s  %-24s\n", n.Name, n.Char)
// 		}
// 	} else {
// 		fmt.Println("No results found")
// 	}
// }

// func getActorInfo(name string) {
// 	stmt := `
// 		MATCH (actor:Actor)-[:PLAYED_CHARACTER]->(character)<-[:FEATURED_CHARACTER]-(episode), (episode)<-[:HAS_EPISODE]-(season)<-[:HAS_SEASON]-(tvshow)
// 		WHERE actor.name = {nameSub}
// 		RETURN tvshow.name AS Show, season.name AS Season, episode.name AS Episode, character.name AS Character
// 	`
// 	params := neoism.Props{"nameSub": name}

// 	res := []struct {
// 		Show      string
// 		Season    string
// 		Episode   string
// 		Character string
// 	}{}

// 	// construct query
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}

// 	// execute query
// 	err := db.Cypher(&cq)
// 	panicErr(err)

// 	if len(res) > 0 {
// 		fmt.Println("Actor: ", name)
// 		for _, n := range res {
// 			fmt.Printf("  %-24s  %-16s  %-24s  %-16s\n", n.Show, n.Season, n.Episode, n.Character)
// 		}
// 	} else {
// 		fmt.Println("No results found")
// 	}
// }

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
		MATCH (perm)-[:hasGroup]->(group)
		WHERE group.name = {groupSub}
		MATCH (perm)-[:hasRole]->(role)
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
		fmt.Println("Group Roles:", group)
		for _, n := range res {
			fmt.Printf("  %-24s\n", n.Name)
		}
	} else {
		fmt.Println("No results found")
	}

}

func listUserGroupsRoles(user string) {
	stmt := `
		MATCH (group:Group)-[:hasRole]->(role)
		WHERE group.name = {groupSub}
		RETURN group.name, role.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Group string `json:"group.name"`
		Role  string `json:"role.name"`
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
		fmt.Println("User Group Roles:", user)
		for _, n := range res {
			fmt.Printf("  %-24s  %-24s\n", n.Group, n.Role)
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
