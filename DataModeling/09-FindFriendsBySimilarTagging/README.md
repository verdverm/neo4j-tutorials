DM.09 Find Friends by Similar Tagging
=====================================================

### Listing Functions

Standard practice

```Go
func listUsers() {...}
func listStuff() {...}
func listTags() {...}
```

Sublists

```Go
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

	...
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

	...
}
```

### Finding by Similar Tagging

Here, we are trying to find the number of common favorites
of a user based on similar tagging of the items.

This function finds and counts the number of paths between
two users which traverse items having common tags.
If
`userA` likes `item1` and `item2`
AND
`userB` likes `item3` and `item4`
Then
`4 paths` will be returned.

```Go
func listSimilarTaggings(user string) {
	stmt := `
		MATCH (me)-[:favorite]->(myFavorites)-[:tagged]->(tag)<-[:tagged]-(theirFavorites)<-[:favorite]-(people)
		WHERE me.name = {userSub} AND NOT me=people
		RETURN people.name AS name, count(*) AS similar_favs
		ORDER BY similar_favs DESC
	`

	params := neoism.Props{"userSub": user}

	res := []struct {
		Name string `json:"name"`
		Favs int    `json:"similar_favs"`
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
			fmt.Printf("  %-10s  %d\n", n.Name, n.Favs)
		}
	} else {
		fmt.Println("No results found")
	}
}
```
