<!-- @format -->

# URL Shortener

A simple URL shortening service built in Go.

## Overview

This project provides a basic URL-shortening service implemented in Go. It allows users to shorten long URLs into more manageable and shareable links. The service also includes a redirect feature to redirect users from the shortened URL to the original long URL.

## Note

Whenever changes are made to the code and the server is restarted, the data stored in the mapping becomes empty. When a URL is shortened and stored in the mapping, if the server or Golang is stopped, the data in the mapping is lost.

Therefore, before calling the Redirect method, the ShortURL method or URL should be triggered so that the value is stored in the mapping. Once the value is stored, the Redirect method will run smoothly. If no data is found in the mapping, the root component page will be hit.
