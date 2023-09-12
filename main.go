package main

import (
	"fmt"
	"project-root/checks"
)

func main() {
	siteURL := "https://steamcommunity.com/"
	isResponsive, err := checks.IsResponsive(siteURL)
	if err != nil {
		fmt.Printf("Erro ao verificar a responsividade do site: %v\n", err)
		return
	}

	if isResponsive {
		fmt.Println("O site é responsivo.")
	} else {
		fmt.Println("O site não é responsivo.")
	}

	isSecure, err := checks.ExamineMediaQueries(siteURL)
	if err != nil {
		fmt.Printf("Erro ao verificar a segurança SSL/TLS: %v\n", err)
		return
	}

	if isSecure {
		fmt.Println("A conexão SSL/TLS é segura.")
	} else {
		fmt.Println("A conexão SSL/TLS não é segura.")
	}

	checkPrivacyPolicy, err := checks.ExaminePrivacyPolicy(siteURL)
	if err != nil {
		fmt.Println("Erro ao verificar a politica de privacidade do site.", err)
	}

	if checkPrivacyPolicy {
		fmt.Println("O site contém política de privacidade.")
	} else {
		fmt.Println("O site não contém política de privacidade .")
	}

}
