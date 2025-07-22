# Sistema de Consulta de Temperatura por CEP

Este é um sistema em Go que recebe um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin).

## Para acessar a aplicação

Acesse a URL:  https://temperature-cep-go-359560015617.us-central1.run.app/temperature?cep=01001000
(substitua `01001000` pelo CEP desejado).

## Requisitos

- Docker
- Docker Compose
- Google Cloud SDK

## Como executar localmente

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/AugustoSandim/temperature-cep-go.git
    cd temperature-cep-go
    ```
2.  **Build e execução com Docker Compose:**
    ```bash
    docker-compose up --build
    ```
3.  **Acessando a aplicação:**
    Abra o seu navegador ou utilize uma ferramenta como o `curl` para testar:
    - **CEP válido:** `curl http://localhost:8080/temperature?cep=01001000`
    - **CEP inválido:** `curl http://localhost:8080/temperature?cep=123`
    - **CEP não encontrado:** `curl http://localhost:8080/temperature?cep=99999999`

## Deploy no Google Cloud Run pelo Terminal

Para fazer o deploy da aplicação no Google Cloud Run, siga estes passos:

1.  **Autentique-se no Google Cloud:**
    ```bash
    gcloud auth login
    gcloud config set project SEU_PROJECT_ID
    ```
2.  **Ative as APIs necessárias:**
    ```bash
    gcloud services enable run.googleapis.com
    gcloud services enable containerregistry.googleapis.com
    ```
3.  **Build e envio da imagem para o Google Container Registry (GCR):**
    ```bash
    gcloud builds submit --tag gcr.io/SEU_PROJECT_ID/temperature-cep-go
    ```
4.  **Deploy no Cloud Run:**
    ```bash
    gcloud run deploy temperature-cep-go \
        --image gcr.io/SEU_PROJECT_ID/temperature-cep-go \
        --platform managed \
        --region us-central1 \
        --allow-unauthenticated
    ```

Após o deploy, o Google Cloud Run fornecerá uma URL pública para acessar a sua aplicação.
