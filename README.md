# GlitchBuster: Order API

## Visão Geral do Projeto

O GlitchBuster é um projeto de estudo focado na construção de sistemas distribuídos e resilientes. A Order API é o primeiro microsserviço deste ecossistema, responsável por gerenciar o ciclo de vida de pedidos. Sua arquitetura é baseada em princípios como Domain-Driven Design (DDD), Event Sourcing, Arquitetura Hexagonal (Ports & Adapters) e Saga Pattern (Orquestração).

O objetivo principal é demonstrar como padrões de resiliência, como Circuit Breaker, Retry Pattern e Dead Letter Queue (DLQ), se combinam para garantir a consistência e a robustez do sistema, mesmo em face de falhas em serviços externos.

## Arquitetura e Padrões de Design

**Fluxo de Eventos (Event-Driven Architecture)**

1. Criação da Ordem: A API recebe uma requisição para criar uma nova ordem.
2. Persistência e Publicação: A API valida a requisição, persiste o evento OrderCreated no EventStoreDB (nossa fonte da verdade) e, em seguida, publica este evento no RabbitMQ.
3. Início do Saga: Um serviço Orquestrador consome o evento OrderCreated do RabbitMQ e inicia uma Saga.
4. Comandos e Eventos: O Orquestrador envia comandos para os serviços participantes (ex: serviço de anti-fraude), que, após executarem sua lógica, publicam eventos de resposta.

**Padrões de Resiliência em Ação**

- Circuit Breaker: Protege o Worker de continuar a se comunicar com serviços externos que estão com falha ou alta latência, evitando que os recursos do sistema sejam esgotados.

- Retry Pattern com Backoff: O Worker tenta novamente as operações que falharam, mas de forma inteligente, com um intervalo crescente e aleatório entre as tentativas, para não sobrecarregar o serviço externo.

- Dead Letter Queue (DLQ): Caso o Worker falhe em processar uma mensagem após todas as tentativas de Retry esgotarem, a mensagem é movida para uma DLQ, garantindo que não haja perda de dados e permitindo uma análise posterior.

**Estrutura de Diretórios (Clean Architecture)**
A estrutura do projeto segue a Arquitetura Limpa, desacoplando a lógica de negócio (domínio) da infraestrutura (banco de dados, frameworks, etc.).

```sh
.
├── cmd/order-api          # Ponto de entrada da aplicação
├── internal/application   # Serviços de aplicação e lógica de negócio
├── internal/domain        # Modelos de domínio e interfaces (Ports)
├── internal/http          # Handlers HTTP
└── pkg                    # Pacotes compartilhados (config, modelos, etc.)
```

### Como Executar

**Pré-requisitos**
- Go 1.25+
- Docker

**Configuração**
Crie o arquivo .env na raiz do projeto com as configurações essenciais.

Construa e execute os contêineres Docker para os serviços externos necessários (EventStoreDB, RabbitMQ, etc.).

## Contribuição

Contribuições são bem-vindas! Se você tiver alguma ideia ou encontrar um bug, sinta-se à vontade para abrir uma issue ou um pull request.
