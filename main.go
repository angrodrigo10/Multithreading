package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/*
"cep": "34012-694",
"logradouro": "Rua Adelaide Pedrosa do Vale",
"complemento": "até 99997/99998",
"unidade": "",
"bairro": "Honório Bicalho",
"localidade": "Nova Lima",
"uf": "MG",
"estado": "Minas Gerais",
"regiao": "Sudeste",
"ibge": "3144805",
"gia": "",
"ddd": "31",
"siafi": "4895"
*/
type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"localidade"`
	Uf         string `json:"uf"`
}

/*
"cep":"34012690",
"state":"MG",
"city":"Nova Lima",
"neighborhood":"Honório Bicalho",
"street":"Rua Liberato Augusto",
"service":"open-cep"
*/
type BrasilAPIResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Cidade     string `json:"city"`
	Uf         string `json:"state"`
}

func fetchViaCEP(cep string, ch chan<- *ViaCEPResponse) {
	//time.Sleep(2 * time.Second) // Pausa por 2 segundos
	req, _ := http.NewRequest("GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- nil
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- nil
		return
	}

	var response ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ch <- nil
		return
	}

	ch <- &response
}

func fetchBrasilAPI(cep string, ch chan<- *BrasilAPIResponse) {
	//time.Sleep(2 * time.Second) // Pausa por 2 segundos
	req, _ := http.NewRequest("GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- nil
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- nil
		return
	}

	var response BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ch <- nil
		return
	}

	ch <- &response
}

func main() {
	cep := "34012694" // Exemplo de CEP

	viaCEPCh := make(chan *ViaCEPResponse)
	brasilAPICh := make(chan *BrasilAPIResponse)

	go fetchViaCEP(cep, viaCEPCh)
	go fetchBrasilAPI(cep, brasilAPICh)

	select {
	case viaCEPResponse := <-viaCEPCh:
		if viaCEPResponse != nil {
			fmt.Printf("Resposta da ViaCEP: %+v\n", viaCEPResponse)
		} else {
			fmt.Println("ViaCEP não retornou dados válidos.")
		}
	case brasilAPIResponse := <-brasilAPICh:
		if brasilAPIResponse != nil {
			fmt.Printf("Resposta da BrasilAPI: %+v\n", brasilAPIResponse)
		} else {
			fmt.Println("BrasilAPI não retornou dados válidos.")
		}
	case <-time.After(time.Second * 1):
		fmt.Println("Erro de timeout: As APIs não responderam a tempo.")
	}
}
