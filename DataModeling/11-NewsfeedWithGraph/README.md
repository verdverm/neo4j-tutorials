DM.11 Implementing Newsfeed with a Graph
========================================

### Listing Functions

Standard practice

```Go
func listUsers() {...}
```

Sublists

```Go
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

	...
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

	...
}
```

### Get status for all users

```Go
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
```

### Get newsfeed for a user

```Go
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
```

### Add a new post for a user

```Go
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
```
