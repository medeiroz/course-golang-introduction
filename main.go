package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5

func main() {
	exibirIntroducao()

	for {
		exibirMenu()
		comando := leOpcao()
		executaComando(comando)
	}
}

func exibirIntroducao() {
	nome := "Flavio"
	versao := 1.0

	fmt.Println("Olá, Sr.", nome)
	fmt.Println("Este programa esta na versão", versao)
}

func exibirMenu() {
	fmt.Println("[1] - Iniciar monitoramento")
	fmt.Println("[2] - Exibir logs")
	fmt.Println("[0] - Sair")

}

func leOpcao() int {
	var comando int
	fmt.Scan(&comando)
	return comando
}

func executaComando(comando int) {
	switch comando {
	case 1:
		monitorar()
	case 2:
		lerLogs()
	default:
		sair()
	}
}

func monitorar() {
	fmt.Println("Monitorando...")
	sites := lerSitesDoArquivo()

	for i := 0; i < monitoramentos; i++ {
		for _, site := range sites {
			resultado := testarSite(site)
			registraLog(site, resultado)
			imprimeResultado(site, resultado)
		}
		time.Sleep(delay * time.Second)
	}

}

func lerLogs() {
	fmt.Println("Carregando Logs...")

	arquivo, err := os.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de log")
		panic(err)
	}

	fmt.Println(string(arquivo))
}

func sair() {
	os.Exit(0)
}

func testarSite(site string) bool {
	response, err := http.Get(site)

	return err == nil && response.StatusCode == 200
}

func lerSitesDoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de sites")
		panic(err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')

		sites = append(sites, strings.TrimSpace(linha))

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return sites
}

func converteResultado(status bool) string {
	if status {
		return "Online"
	}

	return "Offline"
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		fmt.Println("Erro ao abrir o arquivo de log")
		panic(err)
	}

	dataAtual := time.Now().Format("02/01/2006 15:04:05")
	resultado := converteResultado(status)

	log := dataAtual + " - " + site + " - " + resultado + "\n"

	arquivo.WriteString(log)

	arquivo.Close()
}

func imprimeResultado(site string, status bool) {
	fmt.Println("O site", site, "esta", converteResultado(status))
}
