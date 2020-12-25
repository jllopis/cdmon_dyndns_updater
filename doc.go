package main

/*
 * Setup command in CRON:
 * node cdmon_dyndns_updater.js
 *
 * CDMON UPDATE PROTOCOL
 *
 * In order to update your IP you have to make a call to the following URL:
 *
 * https://dinamico.cdmon.org/onlineService.php
 *
 * with the following arguments via GET:
 *
 * enctype = MD5
 * n = username
 * p = password_encoded_with_md5
 *
 * if the IP you want to update is different from the IP assigned by the system
 * you can define your own IP with the argument "cip"
 *
 * cip = x.x.x.x
 *
 * so that we will have:
 *
 * https://dinamico.cdmon.org/onlineService.php?enctype=MD5&n=user&p=1bc29b36f623ba82aaf6724fd3b16718&cip=x.x.x.x
 *
 * where cip is optional since when making the request via URL the server returns a result.
 *
 * RESULTS:
 *
 * the https request returns a variable with the following format:
 *
 * &resultat = result of server request
 *
 * where we have the following options:
 *
 * When a request is made without the cip variable and the authentication has been
 * Correct, it returns the current IP detected by the server:
 *
 * &resultat=guardatok&newip=XXX.XXX.XXX.XXX&
 *
 * When we have sent our IP through the cip variable and authentication
 * has been satisfactory:
 *
 * &resultat=customok&
 *
 * When authenticated but the IP is wrong:
 *
 * &resultat=badip&
 *
 * When authentication has not been successful:
 *
 * &resultat=errorlogin&
 *
 * It exceptional occasions, we can get the following result only when the file that
 * processes all requests has been modified. It is done in order to force all users
 * to update to a new version of the application. We have to contact cdmon to
 * obtain the new URL to make the request;
 *
 * &resultat=novaversio&
 *
 */
