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
		CREATE INDEX ON :Post(date)
	`
	cq := neoism.CypherQuery{
		Statement: stmt,
	}
	err := db.Cypher(&cq)
	panicErr(err)

	stmt = `
		CREATE (bob:User {name:"Bob"})
		CREATE (alice:User {name:"Alice"})
		CREATE (joe:User {name:"Joe"})
		CREATE (bp1:Post {date:1, name:"bob_s1", text:"bobs status1"})
		CREATE (bp2:Post {date:4, name:"bob_s2", text:"bobs status2"})
		CREATE (bp3:Post {date:7, name:"bob_s3", text:"bobs status3"})
		CREATE (bp4:Post {date:8, name:"bob_s4", text:"bobs status4"})
		CREATE (ap1:Post {date:2, name:"alice_s1", text:"Alices status1"})
		CREATE (ap2:Post {date:5, name:"alice_s2", text:"Alices status2"})
		CREATE (jp1:Post {date:3, name:"joe_s1", text:"Joe status1"})
		CREATE (jp2:Post {date:6, name:"joe_s2", text:"Joe status2"})
		CREATE bob-[:FRIEND {status:"CONFIRMED"}]->alice
		CREATE alice-[:FRIEND {status:"CONFIRMED"}]->bob
		CREATE joe-[:FRIEND {status:"CONFIRMED"}]->bob
		CREATE bob-[:FRIEND {status:"CONFIRMED"}]->joe
		CREATE alice-[:FRIEND {status:"PENDING"}]->joe
		CREATE bob-[:POSTED]->bp1
		CREATE bob-[:POSTED]->bp2
		CREATE bob-[:POSTED]->bp3
		CREATE bob-[:POSTED]->bp4
		CREATE bob-[:STATUS]->bp4
		CREATE alice-[:POSTED]->ap1
		CREATE alice-[:POSTED]->ap2
		CREATE alice-[:STATUS]->ap2
		CREATE joe-[:POSTED]->jp1
		CREATE joe-[:POSTED]->jp2
		CREATE joe-[:STATUS]->jp2

	`
	cq = neoism.CypherQuery{
		Statement: stmt,
	}
	err = db.Cypher(&cq)
	panicErr(err)
}

func main() {

	listUsers()
	println()

	listUserFriends("Bob")
	listUserFriends("Alice")
	listUserFriends("Joe")
	println()

	listUserPosts("Bob")
	listUserPosts("Alice")
	listUserPosts("Joe")
	println()

	listUsersStatus()
	println()

	listUsersNewsfeed("Bob")
	listUsersNewsfeed("Alice")
	listUsersNewsfeed("Joe")
	println()

	addNewPost("Alice", "new post", "My first real post", 9)
	addNewPost("Joe", "new post", "My first real post", 10)
	println()

	listUserPosts("Bob")
	listUserPosts("Alice")
	listUserPosts("Joe")
	println()

	listUsersStatus()
	println()

	listUsersNewsfeed("Bob")
	listUsersNewsfeed("Alice")
	listUsersNewsfeed("Joe")
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

func listUserFriends(user string) {
	stmt := `
		MATCH
			(user:User)-[rel:FRIEND]->(friend:User)
		WHERE user.name = {userSub}
		RETURN
			friend.name, rel.status
		ORDER BY
			rel.status, friend.name
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Name   string `json:"friend.name"`
		Status string `json:"rel.status"`
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
			fmt.Printf("  %s: %s\n", n.Name, n.Status)
		}
	} else {
		fmt.Println("No results found")
	}
}

func listUserPosts(user string) {
	stmt := `
		MATCH (user:User)-[:POSTED]-(post:Post)
		WHERE user.name = {userSub}
		RETURN
			post.date AS date,
			post.name AS title,
			post.text AS text
		ORDER BY post.date DESC
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Date  int    `json:"date"`
		Title string `json:"title"`
		Text  string `json:"text"`
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

	fmt.Println(user, " Posts:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %02d: %-12s %q\n", n.Date, n.Title, n.Text)

		}
	} else {
		fmt.Println("No results found")
	}
}

func listUsersStatus() {
	stmt := `
		MATCH (u)-[:STATUS]-(p:Post)
		RETURN
			u.name AS name,
			p.date AS date,
			p.name AS title,
			p.text AS text
		ORDER BY p.date DESC
	`

	res := []struct {
		User  string `json:"name"`
		Date  int    `json:"date"`
		Title string `json:"title"`
		Text  string `json:"text"`
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

	fmt.Println("User Status:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-8s [%02d]: %-12s  %q\n", n.User, n.Date, n.Title, n.Text)

		}
	} else {
		fmt.Println("No results found")
	}
}

func listUsersNewsfeed(user string) {
	stmt := `
		MATCH (me:User)-[rel:FRIEND]->(myfriend:User)-[:POSTED]-(post:Post)
		WHERE me.name = {userSub} AND rel.status = "CONFIRMED"
		RETURN
			myfriend.name AS name,
			post.date AS date,
			post.name AS title,
			post.text AS text
		ORDER BY post.date DESC LIMIT 3
	`
	params := neoism.Props{"userSub": user}

	res := []struct {
		Friend string `json:"name"`
		Date   int    `json:"date"`
		Title  string `json:"title"`
		Text   string `json:"text"`
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

	fmt.Println(user, " Newsfeed:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %-8s [%02d]: %-12s  %q\n", n.Friend, n.Date, n.Title, n.Text)

		}
	} else {
		fmt.Println("No results found")
	}
}

func addNewPost(user, title, text string, date int) {
	stmt := `
		MATCH (user:User)-[r:STATUS]-(post:Post)
		WHERE user.name = {userSub}
		DELETE r
		CREATE (newpost:Post {date:{dateSub}, name:{titleSub}, text:{textSub}})
		CREATE (user)-[:POSTED]->(newpost)
		CREATE (user)-[:STATUS]->(newpost)
	`

	params := neoism.Props{
		"userSub":  user,
		"titleSub": title,
		"textSub":  text,
		"dateSub":  date,
	}

	// construct query
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     nil,
	}

	// execute query
	err := db.Cypher(&cq)
	panicErr(err)

	fmt.Println("New Post:\n  ", user, date, title, text)
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
