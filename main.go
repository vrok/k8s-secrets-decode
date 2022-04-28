package main

import (
	"encoding/base64"
	"flag"
	"io"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

var (
	flagEncode     = flag.Bool("encode", false, "Encode instead of decode")
	flagFile       = flag.String("file", "", "Read from file instead of stdin")
	flagOutputFile = flag.String("output-file", "", "Write output to file instead of stdout")
)

func main() {
	flag.Parse()

	var r io.Reader = os.Stdin
	if *flagFile != "" {
		f, err := os.Open(*flagFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		r = f
	}

	d := yaml.NewDecoder(r)

	var document map[string]interface{}
	if err := d.Decode(&document); err != nil {
		log.Fatalf("failed decoding input: %s", err)
	}

	data, ok := document["data"]
	if !ok {
		log.Fatalf("document is in a wrong format, key 'data' is missing")
	}

	entries, ok := data.(map[interface{}]interface{})
	if !ok {
		log.Fatalf("'data' is not an object but %T", data)
	}

	filterVal := func(value string) (string, error) {
		decoded, err := base64.StdEncoding.DecodeString(value)
		return string(decoded), err
	}

	if *flagEncode {
		filterVal = func(value string) (string, error) {
			return base64.StdEncoding.EncodeToString([]byte(value)), nil
		}
	}

	for k, v := range entries {
		strVal, ok := v.(string)
		if !ok {
			continue
		}

		filtered, err := filterVal(strVal)
		if err != nil {
			log.Fatalf("base64 error for key %s: %s", k, err)
		}
		entries[k] = filtered
	}

	output, err := yaml.Marshal(document)
	if err != nil {
		log.Fatalf("error marhsaling back: %s", err)
	}

	var w io.Writer = os.Stdout
	if *flagOutputFile != "" {
		w, err = os.Create(*flagOutputFile)
		if err != nil {
			log.Fatalf("error opening output file: %s", err)
		}
	}
	w.Write(output)
}
