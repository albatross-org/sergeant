package sergeant

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// magicTable is used to detect the filetype of images attached to a Card.
// Courtest of https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

// mimeFromIncipit returns the mime type of an image file from its first few bytes or the empty string if the
// file does not look like a known file type.
// Courtesy of https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
func mimeFromIncipit(incipit []byte) string {
	incipitStr := string(incipit)
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}

// encodeAsDataURI returns the contents of an image as a data URI. This is used to save two additional requests.
func encodeAsDataURI(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	mime := mimeFromIncipit(content)
	if mime != "image/png" && mime != "image/jpeg" && mime != "image/gif" {
		return "", fmt.Errorf("unrecognised mime type %q", mime)
	}

	prefix := fmt.Sprintf("data:%s;base64,", mime)
	encoded := base64.StdEncoding.EncodeToString(content)

	return prefix + encoded, nil
}
