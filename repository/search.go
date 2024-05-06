package repository

import (
	"strings"
)

/*

	Pour test en local et démontrer skill algo / back

	le top serait d'avoir une data struct type trie
	bonne piste dans algo monster, mais il faudra intégrer la restitution du mot en plus des prefix search


	Solution maison :

	map et slice

	la map contient les occurences uniques des mots
		"abc":		1
		"ab": 		2
		"azerty":	1
		"azer": 	3

	la slice permet de parcourir les clesf plus rapidement ?

		["ab", "abc", "azer", "azerty"]

		il faut veiller à ce qu'elle soit 1a jour et sorted

			maintenir le sort au moment de l'insertion uniquement ?
			bissect => trouver le meilleur index pour insérer l'elem
			...

			sinon juste append et sort, c'est mieux.


	algo:

		findPrefixes

			prefix = "ab"

			1) bsearch pour trouver index de la première occurence
			2) à partir de la liste pour chaque elem qui match le prefix ou bien le mot entier
				conserver la clef avec la valeur max (qui est stockée dans le dict)
				retourner la clef qui correspond.

*/

// return the index of the smallest prefix match from the words list
// -1 if there are no match at all
func bsearch(prefix string, words []string) int {
	left := 0
	right := len(words) - 1
	pivot := (left + right) / 2
	for left <= right {
		if words[pivot] < prefix {
			left = pivot + 1
		} else if words[pivot] > prefix {
			right = pivot - 1
		} else {
			return pivot
		}
		pivot = (left + right) / 2
	}
	// not found
	if pivot == len(words)-1 || pivot == 0 {
		return -1
	}
	if !strings.HasPrefix(words[pivot+1], prefix) {
		return -1
	}
	// if the boundaries are not reached, the prefix is at n+1
	return pivot + 1

}
