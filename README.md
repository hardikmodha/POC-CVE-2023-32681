# POC for CVE-2023-32681

This is a Python 3 implementation of CVE-2023-32681, affecting clients using the [requests](https://github.com/psf/requests) version <= 2.30.0. 

As per https://www.rfc-editor.org/rfc/rfc9110.html#section-15.4 whenever redirection 3xx is issued by the server then user-agent should modify the request 
to remove certain headers and fields from the request. 

Details of those fields can be found at https://www.rfc-editor.org/rfc/rfc9110.html#section-15.4-5. Specific to this CVE I've highlighted the relevant part here. 
 
![image](https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/f0ead2a3-ac4e-4a60-a7a8-f66d2ab573c5)


requests version <= 2.30.0 didn't remove the "Proxy-Authorization" headers while handling the redirection to "https" and thus leaking the proxy credentials to the redirected server.

More details for the same can be found at

- https://www.cve.org/CVERecord?id=CVE-2023-32681
- https://security.snyk.io/vuln/SNYK-PYTHON-REQUESTS-5595532
  
## Pre-requisites

- Valid set of certificate for the HTTP(S) server. You can use [mkcert](https://github.com/FiloSottile/mkcert) to generate the certificates locally. 
- Go v1.19
- Python 3.x
- virtualenv (Optional)

## Overview

 1. Proxy server implementation in Go. It starts a proxy server on Port 8080 and requires a basic authentication. For the POC purpose, it has the hardcoded username and password defined in the same file.  
 2. HTTP(S) server implementation in Go. It starts 2 servers. 
	  - Redirection server listening on Port 443 and another echo server listening on Port 4431.
     - Redirection server defines a route `/redirect` that redirects the traffic to second server running on 4431.
     -  Echo server defines a route `/echoHeaders`. This route returns the headers it received in the request to the client. 
 3. POC Python script. It issues a proxied request to redirection server running on port 443 and prints the response that server returns.
 4. If the response contains the `Proxy-Authorization` headers then the `requests` version is said to be vulnerable to "CVE-2023-32681". 

## Steps

1. Clone this repository.
2. Update the `certFile` and `keyFile` variables in `server/main.go`.
3. Start the HTTP(S) servers by running command `go run server/main.go`. This will start two servers listening on port 443 and 4431.
4. Start the Proxy server by running command `go run proxy/main.go`. This will start the proxy server on port 8080.
6. Install dependencies to run the Python script `pip3 install -r script/requirements_request_2_30_0.txt`
7. Execute the POC script by running `python3 script/poc.py` and observe the output. You can see that "Proxy-Authorization" header is returned by the server in the response. 
Server also prints the received headers in the logs. 
8. Now, upgrade the requests version to >= 2.31.0 or do `pip3 install -r script/requirements_request_2_31_0.txt`
9. Again execute the POC script by running `python3 script/poc.py` and observe the output. You can see that "Proxy-Authorization" header is now not present in the server response.
Same can also be verified in the server logs. 


## Screenshots

### Proxy server

<img width="620" alt="proxy_server_logs" src="https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/60038d79-6591-4a0a-9eab-8736f728fcc8">


### requests==2.30.0

<img width="562" alt="requests_2_30_0_script_output" src="https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/0790195c-69c2-4e8e-8b6f-b6722629e76a">

<img width="555" alt="requests_2_30_0_server_output" src="https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/aa5b8fef-038e-443a-b266-21fe2ce1a2e5">


### requests==2.31.0

<img width="588" alt="requests_2_31_0_script_output" src="https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/3d2fff74-a3f7-4702-bb14-8ed0982b2172">

<img width="542" alt="requests_2_31_0_server_output" src="https://github.com/hardikmodha/POC-CVE-2023-32681/assets/22439276/cdc9b4be-ee04-4e25-a773-391aba1a63d1">



