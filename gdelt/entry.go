package gdelt

import "fmt"

type Entry struct {
	DATE               string
	DocumentIdentifier string
	V2Tone             string
	Themes             string
	Organizations      string
}

func (e Entry) String() string {
	return fmt.Sprintf("[Date:%v, DocumentIdentifier:%v, V2Tone:%v, Themes:%v, Organizations:%v]",
		e.DATE, e.DocumentIdentifier, e.V2Tone, e.Themes, e.Organizations)
}
