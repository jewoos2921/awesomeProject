package synonyms

import (
	"awesomeProject/domainfinder/lib/thesaurus"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{
		APIkey: apiKey,
	}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for "+word, err)
		}
		if len(syns) == 0 {
			log.Fatalln("Couldn't find any synonyms for" + word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
