package checks

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func HaveViewport(siteURL string) (bool, error) {
	resp, err := http.Get(siteURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("o site não está acessível, status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return false, fmt.Errorf("o site não parece ser uma página HTML")
	}

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return false, fmt.Errorf("não foi encontrada a meta tag <meta name=\"viewport\">")
		case html.SelfClosingTagToken, html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "meta" {
				for _, attr := range token.Attr {
					if attr.Key == "name" && attr.Val == "viewport" {
						return true, nil // A meta tag responsiva foi encontrada.
					}
				}
			}
		}
	}
}

func ExamineMediaQueries(siteURL string) (bool, error) {
	resp, err := http.Get(siteURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("o site não está acessível, status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, fmt.Errorf("falha ao analisar o HTML da página: %v", err)
	}

	doc.Find("link[href$='.css']").Each(func(index int, linkHtml *goquery.Selection) {
		cssURL, exists := linkHtml.Attr("href")
		if !exists {
			return
		}

		cssResp, err := http.Get(cssURL)
		if err != nil {
			fmt.Printf("Erro ao recuperar o arquivo CSS (%s): %v\n", cssURL, err)
			return
		}
		defer cssResp.Body.Close()

		cssContent := ""
		if cssResp.StatusCode == http.StatusOK {
			cssContentBytes, _ := io.ReadAll(cssResp.Body)
			cssContent = string(cssContentBytes)
		}

		if strings.Contains(cssContent, "@media") {
			fmt.Printf("O arquivo CSS (%s) contém media queries.\n", cssURL)
		} else {
			fmt.Printf("O arquivo CSS (%s) não contém media queries.\n", cssURL)
		}
	})

	return true, nil
}

func IsResponsive(siteURL string) (bool, error) {
	haveViewport, err := HaveViewport(siteURL)
	if err != nil {
		return false, err
	}

	checkMediaQueries, err := ExamineMediaQueries(siteURL)
	if err != nil {
		return false, err
	}

	if haveViewport || checkMediaQueries {

		return true, nil
	}

	return false, nil

}
