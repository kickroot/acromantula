# acromantula
###Acromantula: Curl with a REPL


### Features
- A request REPL
- Persistent and selectable configurations for request parameters, headers, and base URIs
- Support for GET/HEAD/PUT/POST/DELETE methods
- Automatic JSON formatting
- Easy file uploads for PUT/POST
- Automatic content-type detection for uploads

### License
Apache 2.0

### Upcoming Features
- Proxy support
- Cookie support

### Installation

Acromantula is currently available via source installation.  If you've got Go installed it's as easy as:

```
git clone https://github.com/kickroot/acromantula.git
cd acromantula
go install
```
With this, acromantula will now reside in $GOPATH/bin

### Usage

#### Starting acromantula

```
$> acromantula
Acromantula 0.1.0-alpha
Hit Ctrl+D to quit
acro >>
```

#### Simple Browsing
Just type in the method (GET/POST/PUT/DELETE/HEAD) and the full path to the URL:
```
acro >> head https://api.github.com

<<  HEAD https://api.github.com
 >  User-Agent : [Acromantula 0.1.0-alpha]

<<  HTTP 200 OK
 <  Access-Control-Allow-Origin : [*]
 <  Access-Control-Expose-Headers : [ETag, Link, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval]
 <  Cache-Control : [public, max-age=60, s-maxage=60]
 <  Content-Length : [2039]
 <  Content-Security-Policy : [default-src 'none']
(full headers snipped)

<<  Content:

acro >>
```

#### Posting data
To upload data via POST/PUT, just prepend the filename with an '@'
```
acro >> post http://myserver.com/upload @/path/to/my/file
```
Acromantula will automatically guess the content type from the file extension (if available).
