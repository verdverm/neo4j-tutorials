DM.10 Multirelational (Social) Graphs
=====================================================

### Listing Functions

Standard practice

```Go
func listUsers() {...}
func listStuff() {...}
```

Sublists

```Go
func listUserRelations(user string) {
	stmt := `
		MATCH
			(user:User)-[:FOLLOWS]->(follow),
			(user:User)-[:LIKES]->(like),
			(user:User)-[:LOVES]->(love)
		WHERE user.name = {userSub}
		RETURN
			collect(DISTINCT follow.name) AS Follows,
			collect(DISTINCT like.name) AS Likes,
			collect(DISTINCT love.name) AS Loves
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Follows []string `json:"Follows"`
		Likes   []string `json:"Likes"`
		Loves   []string `json:"Loves"`
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

	fmt.Println(user, " Relations:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  follows: %v\n", n.Follows)
			fmt.Printf("  likes:   %v\n", n.Likes)
			fmt.Printf("  loves:   %v\n", n.Loves)
		}
	} else {
		fmt.Println("No results found")
	}
}
```


### Find Lovers

This function returns people who are in love,
where there is a two-way `[:LOVES]` relationship.

```Go
func findLovers() {
	stmt := `
		MATCH
			(u1:User)-[love:LOVES]->(u2:User)-[:LOVES]->(u1)
		WHERE u1.name <= u2.name
		RETURN u1.name, u2.name
		ORDER BY u1.name
	`

	res := []struct {
		LoverA string `json:"u1.name"`
		LoverB string `json:"u2.name"`
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

	fmt.Println("Lovers:")
	if len(res) > 0 {
		for _, n := range res {
			fmt.Printf("  %s \u2661 %s\n", n.LoverA, n.LoverB)
		}
	} else {
		fmt.Println("No results found")
	}
}
```

The `WHERE` clause eliminates the duplicate, symmetric results like:

```
Lovers:
  Joe ♡ Maria
  Maria ♡ Joe
```

when we only want

```
Lovers:
  Joe ♡ Maria
```

It's a bit of a hack and you should think about situations where lovers share the same name. You can also think about other features on which
duplicates can be removed such as sex or age.
