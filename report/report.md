# Reporte

Dada la propia naturaleza distribuida de BitTorrent, describiremos, a continuación, cada uno de los componentes con respecto a los detalles intrínsecos de BitTorrent (cliente BitTorrent) y los detalles de distribución del Tracker.

## 1. Arquitectura

### 1.1 Cliente BitTorrent

Un cliente de BitTorrent mantiene interacciones (conexiones) con un número de procesos (otros clientes) que puede crecer o disminuir a medida que estos se unen o se van; de ahí la importancia de tener una arquitectura en la que exista una fuerte separación entre procesamiento y coordinación. Con este objetivo en mente, nos apoyamos en el estilo **Publish-Subscribe**, donde la idea es ver al sistema (el propio cliente) como una colección de procesos (hilos) que operan de manera **autónoma**. En este modelo, la coordinación contempla comunicación y la cooperación entre procesos.

Para materializar Publish-Subscribe diseñamos un hilo principal de procesamiento, y múltiples hilos que mantienen la interacción con los clientes en la vecindad. La comunicación solo se produce entre los hilos asociados a clientes (hilos secundarios) con el hilo de procesamiento principal, a través de un canal de comunicación (valga la redundancia), donde los hilos secundarios publican notificaciones que serán recibidas por el hilo principal, actuando en consecuencia. Notemos como este formato se ajusta a la propia naturaleza de BitTorrent y desacopla por completo el procesamiento de la interacción, permitiendo escalar con un número arbitrario de vecinos.

![](./images/publish_subscribe.jpg)

### 1.2 Distribución del Tracker

## 2. Procesos

### 2.1 Cliente BitTorrent

Cada cliente BitTorrent representa un único proceso, formado por múltiples hilos (descritos anteriormente). 

Decidimos utilizar Go para implementar los clientes, dado que las goroutines (hilos en Go) implementan un modelo many-to-many de threading, donde combinan el uso de user-threads y kernel-threads (muchos kernel-threads por proceso y muchos user-threads por kernel-thread, estos últimos a su vez manejados por un user-level thread package), lo que:

- Permite que crear, destruir y sincronizar hilos sea relativamente barato y no involucre la intervención del kernel en absoluto
- Permite que un blocking system-call no suspenda al proceso por completo
- Permite que la interacción con los hilos a nivel de aplicación sea totalmente transparente, el programador no se entera de la existencia de los kernel-threads
- Permite que los kernel-threads sean utilizados en entornos de multiprocesamiento (a diferencia de los hilos en Python) al ejecutarlos en CPUs/cores distintos; de nuevo, de una manera totalmente transparente

![](./images/many_to_many_threading.png)

**Nota**: Imagen extraída de: *M. van Steen and A.S. Tanenbaum, Distributed Systems, 4th ed., distributed-systems.net, 2023*

### 2.2 Distribución del Tracker

## 3. Comunicación

### 3.1 Cliente BitTorrent

Para la comunicación entre clientes, utilizamos directamente la **interfaz de socket** (sobre TCP) ofrecida por el lenguaje de programación. A su vez, con el objetivo de abstraer los detalles relacionados con la interpretación y manejo de los bit enviados/recibidos, definimos una interfaz **Messenger**.

Notemos que al utilizar los sockets puros solo podemos disponer de una comunicación one-to-one entre clientes; esto, por supuesto, lo tuvimos en cuenta a la hora de realizar el diseño, sin embargo, no descartamos el uso de patrones de mensajería como **ZeroMQ** (permite comunicaciones one-to-many y many-to-many) en un futuro para refactorizar y quizás optimizar el funcionamiento de cada cliente.

### 3.2 Distribución del Tracker

## 4. Coordinación

### 4.1 Cliente BitTorrent

Como mencionamos anteriormente (sección de arquitectura), los hilos que conforman al cliente están desacoplados referencialmente, sin embargo, están acoplados temporalmente (todos se encuentran en funcionamiento), esto da lugar a una **coordinación basada en eventos** (event-based coordination). En sistemas desacoplados referencialmente los procesos (hilos) no conocen de la existencia de otros de manera explícita, lo único que pueden hacer es publicar la notificación describiendo la ocurrencia de un evento.

De nuevo, Go facilita mucho la implementación de este diseño al proveer los **canales** (channels) que actúan como bus para que los hilos secundarios publiquen sus notificaciones y el hilo principal pueda recibirlas y procesarlas. Es importante destacar que para un flujo eficiente, el procesamiento realizado por el hilo principal debe estar constituido solo por una cantidad constante de operaciones en CPU, cualquier otra operación cuya duración sea dependiente de la entrada y/o de otro proceso/hilo debe convertirse en un hilo secundario.

### 4.2 Distribución del Tracker