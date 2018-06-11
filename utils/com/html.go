package com

import (
	"html"
	"regexp"
	"strings"
)

// Html2JS converts []byte type of HTML content into JS format.
func Html2JS(data []byte) []byte {
	s := string(data)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\"", `\"`, -1)
	s = strings.Replace(s, "<table>", "&lt;table>", -1)
	return []byte(s)
}

// encode html chars to string
func HtmlEncode(str string) string {
	return html.EscapeString(str)
}

// decode string to html chars
func HtmlDecode(str string) string {
	return html.UnescapeString(str)
}

// strip tags in html string
func StripTags(src string) string {
	//将HTML标签全转换成小写
	re := regexp.MustCompile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//remove tag <style>
	re = regexp.MustCompile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//remove tag <script>
	re = regexp.MustCompile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//replace all html tag into \n
	re = regexp.MustCompile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//trim all spaces(2+) into \n
	re = regexp.MustCompile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

// change \n to <br/>
func Nl2br(str string) string {
	return strings.Replace(str, "\n", "<br/>", -1)
}
