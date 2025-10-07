package compress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"log"
)

func CompressBytes(value []byte) ([]byte, error) {
	// As far as I can tell, errors in these zlib calls are impossible when
	// writing to a bytes.Buffer.
	var buf bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&buf, zlib.BestSpeed)
	defer zw.Close()
	zw.Write(value)
	zw.Flush()
	zw.Close()
	return buf.Bytes(), nil
}

func CompressBytesLevel(value []byte, level int) ([]byte, error) {
	// As far as I can tell, errors in these zlib calls are impossible when
	// writing to a bytes.Buffer.
	var buf bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&buf, level)
	defer zw.Close()
	zw.Write(value)
	zw.Flush()
	zw.Close()
	return buf.Bytes(), nil
}

func UncompressBytes(zValue []byte) ([]byte, error) {
	b := bytes.NewBuffer(zValue)
	r, err := zlib.NewReader(b)
	if err != nil {
		log.Printf("zlib.NewReader somehow returned an error: %+v", err)
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
