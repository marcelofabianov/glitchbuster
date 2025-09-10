# Estudo de Resiliencia

## Atores

API (Produtor de Eventos): Este é o ponto de entrada. Ele recebe a requisição para criar a Order. Sua principal responsabilidade é validar a requisição, gerar o evento, persistir esse evento no EventStoreDB e, só então, publicar a mensagem no RabbitMQ. A persistência no Event Store antes de publicar na fila é crucial para a durabilidade.

EventStoreDB: Este é o banco de dados de eventos, sua fonte da verdade. A sua lista está perfeita. Ele armazena os eventos de forma imutável e ordenada. O papel dele é garantir que, em caso de falha de qualquer outro componente, o estado do sistema possa ser reconstruído a qualquer momento.

RabbitMQ (Broker de Eventos): Perfeito. É o seu broker de mensagens. Ele atua como um sistema de filas, desacoplando o produtor (a API) do consumidor (o Worker). Sua principal função é garantir que a mensagem chegue ao Worker de forma assíncrona. Ele também gerencia a DLQ para onde as mensagens que falharam serão movidas.

Worker (Consumidor da Fila): Exatamente. O Worker é o seu processador de lógica de negócio. Ele consome a mensagem do RabbitMQ e executa as ações necessárias, como a validação de anti-fraude que simulamos. Este é o componente onde os padrões de Retry, Backoff e Circuit Breaker são implementados para lidar com a falha do serviço externo.

Serviço Externo (Anti-Fraude, etc.): Este é o ator que o seu Worker consome. Ele pode ser um serviço próprio ou de terceiros. Sua instabilidade é o motivo da existência da maioria dos padrões de resiliência que aplicamos.

Banco de Dados PostgreSQL (Projeção/Estado Final): Este é um ponto de refinamento importante na sua lista. O PostgreSQL, neste cenário de Event Sourcing, não armazena registros do processamento da forma tradicional. Em vez disso, ele é uma projeção (ou banco de dados de leitura). Os dados no PostgreSQL são gerados a partir da leitura dos eventos do EventStoreDB, permitindo consultas rápidas para a sua aplicação, como "me mostre o estado atual do pedido". O EventStoreDB é a fonte da verdade; o PostgreSQL é apenas uma visão otimizada para leitura.

### Fluxo de Eventos e Padrões em Ação
Vamos seguir o caminho de um evento, como OrderCreated, através do sistema.

1. Geração e Persistência do Evento (Durabilidade e Ordem)
Um Produtor de Eventos (seu serviço de e-commerce, por exemplo) recebe um pedido e antes de fazer qualquer outra coisa, persiste o evento OrderCreated no Event Store.

Ação do Event Store: Esta é a sua base confiável e ordenada. O Event Store garante que o evento seja gravado de forma atômica e imutável. Isso é a sua fonte única de verdade. Se algo der errado depois, você sempre pode reconstruir o estado a partir daqui.

Após a persistência bem-sucedida, o produtor publica o evento em uma fila principal do RabbitMQ.

2. Processamento Assíncrono (Desacoplamento)
O Worker (o nosso cliente com o pool de workers) consome a mensagem da fila principal do RabbitMQ.

Ação do Worker: O worker agora tem a tarefa de processar a ordem, o que envolve comunicar com o Serviço Externo de anti-fraude.

3. Resiliência em Ação (Retry e Circuit Breaker)
O worker tenta se comunicar com o serviço de anti-fraude.

Cenário de Sucesso: Se o serviço externo responder com 200 OK, a tarefa é concluída. O worker envia um "ack" (acknowlegment) para o RabbitMQ, removendo a mensagem da fila.

Cenário de Falha Transitória: Se o serviço externo falhar (ex: 500 Internal Server Error ou um timeout de rede), a lógica do Retry Pattern com Backoff entra em ação. O worker tenta novamente a requisição (por exemplo, 3 vezes), esperando um tempo maior a cada nova tentativa.

Cenário de Falha Sustentada: Se, após as 3 tentativas, a requisição ainda falhar, o worker conclui que a falha é permanente. Ele reporta esta falha final ao Circuit Breaker. Ao atingir o limite configurado de falhas, o disjuntor "arma", mudando para o estado Aberto.

Proteção do Sistema: Agora, qualquer novo worker que tente processar uma ordem e chame o serviço de anti-fraude será imediatamente barrado pelo Circuit Breaker. A chamada de rede nem é feita. O worker recebe o erro do disjuntor e sabe que não adianta prosseguir.

4. Roteamento para a DLQ (Preservação de Dados)
Quando uma mensagem não pode ser processada (seja por esgotar as tentativas de retry ou por ser barrada pelo disjuntor), o worker não envia o ack. O RabbitMQ percebe que a mensagem não foi processada e, após um número pré-definido de tentativas de reentrega, a move automaticamente para uma Dead Letter Queue (DLQ).

Ação da DLQ: A DLQ é um local seguro para as mensagens falhas. Elas não são perdidas. Uma equipe de suporte ou um serviço de monitoramento pode inspecionar a DLQ para entender a causa da falha e decidir se a mensagem deve ser reprocessada.


