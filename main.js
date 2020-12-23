const { Resolver } = require("dns").promises;
const https = require("https");

// Comando para ejecutar el CRON:
// node cdmon_dyndns_updater.js

/**
 *
 * PROTOCOLO
 *
 * Para poder actualizar su IP tiene que hacer una llamada a la siguiente URL:
 *    https://dinamico.cdmon.org/onlineService.php
 * con los argumentos via GET siguientes:
 *    enctype=MD5
 *    n=nombre_de_usuario
 *    p=contrasea_codificada_con_md5
 * si la IP que quiere actualizar es diferente a la IP que le asigna el sistema
 * puede definir una IP propia con el argumento "cip"
 *    cip=x.x.x.x
 * de modo que tendremos:
 *    https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=usuario&p=1bc29b36f623ba82aaf6724fd3b16718&cip=x.x.x.x
 * donde cip es opcional ya que al hacer la peticin via URL el servidor devuelve
 * un resultado.
 *
 *    RESULTADOS:
 * la peticion https nos devuelve una variable con el formato siguiente:
 *    &resultat=resultado de la petición del servidor&
 * donde tenemos las siguientes opciones:
 *
 * Cuando se hace una petición sin la variable cip y la autentificación ha sido
 * correcta nos devuelve la IP actual que detecta el servidor.
 *    &resultat=guardatok&newip=x.x.x.x&
 *
 * Cuando hemos mandado nuestra IP mediante la variable cip y la autentificación
 * ha sido satisfactoria.
 *    &resultat=customok&
 *
 * Nos devuelve este resultado cuando la autentificación ha sido
 * pero la IP es erronea.
 *    &resultat=badip&
 *
 * Nos devuelve este resultado cuando la autentificación no ha sido satisfactoria.
 *    &resultat=errorlogin&
 *
 * Nos devuelve este resultado en raras ocasiones, solo cuando modificamos el
 * archivo que procesa todas las peticiones para obligar a todos los usuarios a
 * actualizar a una nueva version de la aplicacion. En su caso solo tendra
 * que ponerse en contacto con nosotros para obtener la nueva URL para hacer la petición.
 *    &resultat=novaversio&
 *
 */

//           AQUI COMIENZA

// Comenzamos con los datos de usuario de CDMON.COM
// Le debes dar valores a las variables.
// USUARIO = es el nombre de usuario para entrar en CDMON.COM
// PASSWORDMD5 = Es la contraseña para entrar en CDMON.COM encriptada con el algoritmo MD5.
// EMAIL = es donde queremos que lleguen los mensajes del CRON.
// HOST = el dominio/subdominio que se desea actualizar

const updates = [
  {
    user: "dyndnsupdater",
    passmd5: "cc2c3715b70c11f62c9ac6c70389e957",
    email: "jllopis@gimlab.net",
    host: "gimlab.net",
  },
  {
    user: "srv1gimlabupdater",
    passmd5: "56f1c60558111852a08654d49b04d3ac",
    email: "jllopis@gimlab.net",
    host: "mx.gimlab.net",
  },
  {
    user: "mxgimlabupdater",
    passmd5: "48834d0b40252d0960e35b624efac8c7",
    email: "jllopis@gimlab.net",
    host: "srv1.gimlab.net",
  },
];

// use the resolver form cdmon to be sure we get the correct registered IP
const resolver = new Resolver();
resolver.setServers(["46.16.60.166", "46.16.60.159", "35.156.85.88"]);

async function getExternalIp(domain) {
  let address = await resolver.resolve4(domain);
  return address;
}

function getRegisteredIp(updaterObj) {
  let user = updaterObj.user;
  let passmd5 = updaterObj.passmd5;
  const cdmonUrl = `https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=${user}&p=${passmd5}`;

  return new Promise((resolve, reject) => {
    https.get(cdmonUrl, (res) => {
      res.setEncoding("utf8");
      let body = "";

      res.on("data", (data) => {
        body += data;
      });

      res.on("end", () => {
        let exp = /&newip=\b(?:\d{1,3}\.){3}\d{1,3}\b&resultat=.*/i;
        if (body.match(exp)) {
          let ip = body.split("&")[1].split("=")[1];
          resolve(ip);
        }
      });

      res.on("error", (error) => {
        reject(error);
      });
    });
  });
}

function updateRegisteredIp(updaterObj) {
  let user = updaterObj.user;
  let passmd5 = updaterObj.passmd5;
  let newIp = updaterObj.newIp;

  const cdmonUrl = `https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=${user}&p=${passmd5}&cip=${newIp}`;

  return new Promise((resolve, reject) => {
    https.get(cdmonUrl, (res) => {
      res.setEncoding("utf8");
      let body = "";

      res.on("data", (data) => {
        body += data;
      });

      res.on("end", () => {
        if (body.includes("customok")) {
          resolve(true);
        }
      });

      res.on("error", (error) => {
        reject(error);
      });
    });
  });
}

async function run() {
  for await (updater of updates) {
    // Averiguar la ip que tenemos registrada en cdmon
    let currentIP = await getExternalIp(updater.host);
    console.log(`${updater.host} Current IP = ${currentIP}`);

    // Averiguar la ip actual a través de whatismyip.com
    let registeredIP = await getRegisteredIp(updater);
    console.log(`${updater.host} Registered IP = ${registeredIP}`);

    // Comparar si no son iguales
    if (currentIP != registeredIP) {
      // en cuyo caso, actualizar cdmon dyn dns
      console.log(
        `${updater.host} RegisteredIP=${registeredIP} currentIP=${currentIP} => updating`
      );
      updater.newIp = currentIP;
      let res = await updateRegisteredIp(updater);
      if (res) {
        console.log(
          `Updated ${updater.host} IP: old=${registeredIP} new=${currentIP}`
        );
      }
    } else {
      console.log(`${updater.host} IP not changed: ${registeredIP}`);
    }
  }
}

run();

/*
#establece una variable con el GET que tiene que hacer, con todos los datos
    CHANGE_IP="https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=$USUARIO&p=$PASSMD5&cip=$IP_ACTUAL"
	# luego al establecer la variable RESULTADO, hace el GET y la variable se queda con la respuesta que le da
	# si es satisfactorio, la respuesta deberÃ­a ser &resultat=customok&
	RESULTADO=`wget $CHANGE_IP -o /dev/null -O /dev/stdout --no-check-certificate`
	#Ponemos que es lo que queremos que salga en el email
	MENSAJE="Ha habido un cambio en la IP de los nombres de dominio.\n"
	MENSAJE=$MENSAJE"Se han actualizado los servidores DNS dinamicos de CDMON.\n"
	MENSAJE=$MENSAJE"El resultado devuelto ha sido el siguiente:\n"

	#Finalmente envia un email con los resultados
	echo -e $MENSAJE $RESULTADO IP DEL SITIO era: $IP_DNS_ONLINE por lo tanto fue modificada por la IP ACTUAL:$IP_ACTUAL | mail $EMAIL -s "cambio de IP"

*/
