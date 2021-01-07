# admin

Para poder utilizar el código se deben seguir los siguientes pasos:

- Ingrese a las máquinas virtuales:
1. DNS_1: 10.10.28.17
2. DNS_2: 10.10.28.18
3. DNS_3: 10.10.28.19 
4. Broker: 10.10.28.20

- Se deben activar los servidores en cada máquina virtual, para ello, en una terminal de las máquinas virtuales correspondientes al DNS_1, DNS_2 y DNS_3 se debe ingresar el siguiente comando:

>`` make dns ``

- Para la máquina virtual del broker se debe ingresar:

>``make broker``

- Finalmente, para utilizar el código debemos abrir otra terminal donde se encuentra la máquina virtual del DNS_1 donde debemos ingresar lo siguiente, para actuar como administrador:

>``make admin``

- Y lo siguiente, para actuar como cliente:
  
>``make client``


Los comandos deben ser ingresados como se explica en el enunciado, es decir, estos deben ser {Create, Update, Delete, Get} y deben ser seguidos por un espacio y el nombre.dominio, en otras palabras, comandos del tipo "www.google.com" no funcionarán.

Para terminar debemos mencionar que no se estableció cómo se le ingresaría la IP al momento de crear una nueva ruta, por lo que optamos por poner una por default, que debe ser cambiada con el comando "Update".