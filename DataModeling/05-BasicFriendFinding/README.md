DM.05 BasicFriendFinding
================

In this section we will explore some basic friend finding activities.

### Friend List

Here is a function to get a user's friend list.

```Go
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
```

### Recommending friends

Here's a function to get a list of friend recommendations.
It works by finding all `friends_of_friend` which the user
is not already friends with, and includes
the number of connections to the `friend_of_friend`.

```Go
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
```
