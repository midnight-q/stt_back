package stt_converter

import (
	"fmt"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"sort"
	"strings"
	"stt_back/services/file_storage"
)

type Params struct {
	TimeFrame         int
	IsShowEmotion     bool
	IsShowSpeaker     bool
	IsShowTag         bool
	IsShowPunctuation bool
}

func ConvertDataToText(data []Data, params Params) (res string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyTextTemplate(d, params))
	}
	return strings.Join(resultArr, "\n")
}

func ConvertDataToHtml(data []Data, params Params) (res string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyHtmlTemplate(d, params))
	}
	return strings.Join(resultArr, "")
}

func ConvertDataToPdf(data []Data, params Params) (path string) {
	html := ConvertDataToHtml(data, params)
	html = fmt.Sprintf(`<html>
									<head>
									  <title>My First HTML</title>
									  <meta charset="UTF-8">
									</head>
									<body>
										<link href='http://fonts.googleapis.com/css?family=Jolly+Lodger' rel='stylesheet' type='text/css'>
										 <style type = "text/css">
											p { font-family: 'Roboto', sans-serif;; }
										</style>
										%s
									</body>
								</html>`, html)
	pdfg, err :=  wkhtml.NewPDFGenerator()
	if err != nil{
		fmt.Println("Error in create pdf:", err)
	}

	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(html)))
	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		fmt.Println("Error in create pdf:", err)
	}

	b := pdfg.Bytes()
	path, _ = file_storage.CreateFileInLocalStorage(b, ".pdf")
	return
}

func ConvertDataToDoc(data []Data, params Params) (path string) {

	return path
}

func applyTextTemplate(data Data, params Params) string {
	str := data.RawText
	if params.IsShowPunctuation {
		str = data.Text
	}
	if params.IsShowTag && params.IsShowPunctuation {
		res := ""
		for i, s := range strings.Split(str, "") {
			isFind := false
			for _, tag := range data.Tags {
				if i == tag.Start {
					res += "[" + tag.Name + " " + s
					isFind = true
				} else if i == tag.End {
					res += s + "]"
					isFind = true
				}
			}
			if !isFind {
				res += s
			}
		}
		str = res
	}
	if params.IsShowSpeaker || params.IsShowEmotion {
		str = ":" + str
	}
	if params.IsShowEmotion {
		str = "(" + data.Emotion + ")" + str
	}
	if params.IsShowSpeaker {
		str = data.Speaker + str
	}
	return str
}

func applyHtmlTemplate(data Data, params Params) string {
	str := data.RawText
	if params.IsShowPunctuation {
		str = data.Text
	}
	if params.IsShowTag && params.IsShowPunctuation {
		res := ""
		for i, s := range strings.Split(str, "") {
			isFind := false
			for _, tag := range data.Tags {
				if i == tag.Start {
					res += "<b><i>" + tag.Name + "</i> " + s
					isFind = true
				} else if i == tag.End {
					res += s + "</b>"
					isFind = true
				}
			}
			if !isFind {
				res += s
			}
		}
		str = res
	}
	str = fmt.Sprintf("<div>%s</div>", str)

	if params.IsShowEmotion {
		str = fmt.Sprintf("<div>%s</div>", data.Emotion) + str
	}
	if params.IsShowSpeaker {
		str = fmt.Sprintf("<div>%s</div>", data.Speaker) + str
	}
	return fmt.Sprintf("<p>%s</p>", str)
}
