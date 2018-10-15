package xcharset

import (
	iconv "github.com/djimenez/iconv-go"
	"github.com/saintfish/chardet"
	"strings"
)

func DetectCharset(str string, isHtml bool) (string, error) {
	var d *chardet.Detector
	if isHtml {
		d = chardet.NewHtmlDetector()
	} else {
		d = chardet.NewTextDetector()
	}
	r, e := d.DetectBest([]byte(str))
	if e != nil {
		return "", e
	} else {
		return r.Charset, nil
	}
}

func ConvCharset(src string, fromCharset string, toCharset string) (string, error) {
	fromCharset = strings.ToLower(fromCharset)
	toCharset = strings.ToLower(toCharset)
	if fromCharset == "gb-18030" {
		fromCharset = "gb18030"
	}
	if toCharset == "gb-18030" {
		toCharset = "gb18030"
	}
	if fromCharset == toCharset {
		return src, nil
	}
	rst, e2 := iconv.ConvertString(src, fromCharset, toCharset)
	return rst, e2
}

func ConvCharsetAuto(src string, isHtml bool, toCharset string) (string, error) {
	fromCharset, e := DetectCharset(src, isHtml)
	if e != nil {
		return "", e
	}
	return ConvCharset(src, fromCharset, toCharset)
}
