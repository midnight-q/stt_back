package stt_converter

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"stt_back/services/file_storage"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gingfrederik/docx"
	"github.com/orcaman/writerseeker"
)

type Params struct {
	TimeFrame         int
	IsShowEmotion     bool
	IsShowSpeaker     bool
	IsShowTag         bool
	IsShowPunctuation bool
	NamedEntityTypes  []string
}

func ConvertDataToText(data []Data, params Params) (path, text string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyTextTemplate(d, params))
	}
	resultText := strings.Join(resultArr, "\n")
	path, _ = file_storage.CreateFileInLocalStorage([]byte(resultText), ".txt")
	return path, resultText
}

func ConvertDataToHtml(data []Data, params Params) (path string, content string) {
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyHtmlTemplate(d, params))
	}
	innerHtml := strings.Join(resultArr, "")
	html := fmt.Sprintf(`<html>
									<head>
									  <title>Result</title>
									  <meta charset="UTF-8">
									</head>
									<body>
										<link href='http://fonts.googleapis.com/css?family=Jolly+Lodger' rel='stylesheet' type='text/css'>
										 <style type = "text/css">
											p { font-family: 'Roboto', sans-serif;; }
										</style>
										%s
									</body>
								</html>`, innerHtml)

	path, _ = file_storage.CreateFileInLocalStorage([]byte(html), ".html")
	return path, html
}

func ConvertDataToPdf(data []Data, params Params) (path string) {
	_, html := ConvertDataToHtml(data, params)
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
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
	resultArr := []string{}
	sort.Slice(data, func(i, j int) bool {
		return data[i].TimeStart < data[j].TimeStart
	})
	for _, d := range data {
		resultArr = append(resultArr, applyTextTemplate(d, params))
	}
	f := docx.NewFile()
	for _, line := range resultArr {
		para := f.AddParagraph()
		para.AddText(line)
		_ = f.AddParagraph()
	}

	buf := writerseeker.WriterSeeker{}
	_ = f.Write(&buf)
	b, _ := ioutil.ReadAll(buf.Reader())

	path, _ = file_storage.CreateFileInLocalStorage(b, ".docx")

	return path
}

func applyTextTemplate(data Data, params Params) string {
	if data.IsTimeFrameLabel {
		return fmt.Sprintf("[%s]", data.Text)
	}
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
	if data.IsTimeFrameLabel {
		return fmt.Sprintf("<p>[%s]</p>", data.Text)
	}
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
					res += "<strong><em>" + tag.Name + "</em> " + s
					isFind = true
				} else if i == tag.End {
					res += s + "</strong>"
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
