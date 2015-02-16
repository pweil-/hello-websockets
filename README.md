# hello-websockets

# Build the server executable
This will install the necessary libraries, build the server executable, and also
build the docker file.

```
cd <your directory where you cloned the repo>/openshift
./build.sh
```

# Running the docker container manually
```
$ docker run -p 9999:9999 -p 9443:9443 pweil/hello-websocket
Serving at :9999
Serving TLS at :9443
```

# Running in OpenShift

This assumes you are running in the vagrant environment and that OpenShift is running.  If you are not, please adjust
the ip addresses accordingly.  Your vagrant machine must be set up with an IP address that can be reached from the 
*host* machine.  This can be accomplished by using the private or public network options.

Example:

`override.vm.network "private_network", ip: "10.245.2.2"`

# Install the router and demo pod
```
$ vagrant ssh master

$ hack/install-router.sh router https://10.0.2.15:8443
Creating router file and starting pod...
router

$ cd
$ git clone https://github.com/pweil-/hello-websockets.git
$ cd hello-websockets/openshift

$ osc create -f pod.json
hello-websocket

$ osc get pods
POD                 IP                  CONTAINER(S)                   IMAGE(S)                          HOST                           LABELS                 STATUS
hello-websocket     172.17.0.3          hello-websocket                pweil/hello-websocket             openshiftdev.local/127.0.0.1   name=hello-websocket   Running
router              172.17.0.2          origin-haproxy-router-router   openshift/origin-haproxy-router   openshiftdev.local/127.0.0.1   <none>                 Running
```

# Testing the routes

All of these use cases build upon each other.  If you are not following them in order please add the appropriate steps
to create any missing services or routes.


##  Unsecure websocket

```
$ osc create -f service_unsecure.json 
hello-ws

$ osc create -f route_unsecure.json 
ws-unsecure

```

Now, on your *host* machine you'll want to run the client.html in a browser.  The client uses a url of www.example.com.
You can add an entry into your *host* machine's `/etc/hosts` that points to your vagrant machine's IP Address

```
$ cat /etc/hosts
127.0.0.1	localhost.localdomain localhost hello-external.v3.rhcloud.com 
::1		localhost6.localdomain6 localhost6
10.245.2.2 www.example.com
```

To test this route click the Test Unsecure link.  You should see output as depicted below.

![Testing unsecure route](https://github.com/pweil-/hello-websockets/blob/master/openshift/test_images/unsecure_route.png)

##  Edge terminated websocket

```
$ osc create -f route_edge.json 
ws-edge
```

Now, to test a secure route we must first accept the self signed certificate.  Otherwise, the websocket connection will
fail.  To do this, you can navigate to [https://www.example.com/echo] and you should be prompted by your browser.  Follow
your browser's instructions to import/accept the certificate.

To test this route click the Test Secure link. You should see output as depicted below.

![Testing edge route](https://github.com/pweil-/hello-websockets/blob/master/openshift/test_images/edge_route.png)

##  Pass through route

```
$ osc delete route ws-edge
ws-edge
$ osc create -f route_passthrough.json
ws-passthrough
$ osc create -f service_secure.json
hello-ws-secure
```

To test this route click the Test Secure link. You should see output as depicted below.

![Testing pass through route](https://github.com/pweil-/hello-websockets/blob/master/openshift/test_images/route_passthrough.png)

##  Re-encrypt route

```
$ osc delete route ws-passthrough
ws-passthrough

$ osc create -f route_reencrypt.json 
ws-reencrypt

# check the cert being served by the pod is for www.example.com
$ openssl s_client -servername www.example.com -connect $(osc get -t "{{.portalIP}}:{{.port}}" service hello-ws-secure) | grep subject
depth=0 CN = www.example.com, ST = SC, C = US, emailAddress = example@example.com, O = Example, OU = Example
verify error:num=20:unable to get local issuer certificate
verify return:1
depth=0 CN = www.example.com, ST = SC, C = US, emailAddress = example@example.com, O = Example, OU = Example
verify error:num=27:certificate not trusted
verify return:1
depth=0 CN = www.example.com, ST = SC, C = US, emailAddress = example@example.com, O = Example, OU = Example
verify error:num=21:unable to verify the first certificate
verify return:1
subject=/CN=www.example.com/ST=SC/C=US/emailAddress=example@example.com/O=Example/OU=Example
^C

# check the cert being served by the router is www.example2.com (which matches our route config)
$ openssl s_client -servername www.example2.com -connect 10.0.2.15:443 | grep subject
depth=1 C = US, ST = SC, L = Default City, O = Default Company Ltd, OU = Test CA, CN = www.exampleca.com, emailAddress = example@example.com
verify error:num=19:self signed certificate in certificate chain
verify return:0
subject=/CN=www.example2.com/ST=SC/C=SU/emailAddress=example@example.com/O=Example2/OU=Example2
^C
```

To test this route, we need to adjust `/etc/hosts` to now point the www.example2.com domain at our router ip address.

```
$ cat /etc/hosts
127.0.0.1	localhost.localdomain localhost hello-external.v3.rhcloud.com 
::1		localhost6.localdomain6 localhost6
10.245.2.2 www.example2.com
```

Now, test by clicking the Test Reencrypt link in the browser. 


![Testing reencrypt route](https://github.com/pweil-/hello-websockets/blob/master/openshift/test_images/route_reencrypt.png)
