package lib

import (
	"bufio"
	"os"
	"log"
	"fmt"
)

type Deduper struct {
	lsh LSH
	minhash MinHasher
	readfn func(*bufio.Reader, chan Document)
}

func MakeDeduper(lsh LSH, minhash MinHasher, readfn func(*bufio.Reader, chan Document)) *Deduper {
	dd := new(Deduper)
	dd.lsh = lsh
	dd.minhash = minhash
	dd.readfn = readfn
	return dd
}

func (dd Deduper) Fingerprint(shingles []uint32) []uint32 {
	return dd.minhash.Hash(shingles)
}

func (dd Deduper) Index(key string, prints []uint32) {
	dd.lsh.Insert(key, prints)
}

func (dd Deduper) Query(prints []uint32) []string {
	return dd.lsh.Query(prints)
}

func (dd Deduper) Dedupe(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.Println("--- First pass, indexing documents")
		
	reader := bufio.NewReader(file)
	doccount := 0

	doc_chan := make(chan Document)
	go dd.readfn(reader, doc_chan)

	for doc := range doc_chan {
		doccount++
		if (doccount % 10000) == 0 {
			log.Println(doccount, "docs")
		}

		sigs := Shingle(doc.Text)
		sigs = dd.Fingerprint(sigs)
		dd.Index(doc.Id, sigs)
	}

	log.Println("--- Second pass, identifying duplicates")
		
	file.Seek(0, 0)
	reader = bufio.NewReader(file)
	id2cluster := make(map[string]string, doccount)
	doccount = 0

	doc_chan = make(chan Document)
	go dd.readfn(reader, doc_chan)
	
	for doc := range doc_chan {
		doccount++
		if (doccount % 10000) == 0 {
			log.Println(doccount, "docs")
		}
		
		cluster, ok := id2cluster[doc.Id]
		if ok {
			fmt.Println(cluster, doc.Id, doc.Name)
			continue
		}
		
		id2cluster[doc.Id] = doc.Id
		fmt.Println(doc.Id, doc.Id, doc.Name)
		
		sigs := Shingle(doc.Text)
		sigs = dd.Fingerprint(sigs)
		dupes := dd.Query(sigs)
		for _, d := range dupes {
			if _, ok := id2cluster[d]; !ok {
				id2cluster[d] = doc.Id
			}
		}
	}

}
