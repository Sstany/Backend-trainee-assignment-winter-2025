#/bin/sh

hey -n 100000 -q 1000 -m POST -D ./post_auth.txt http://localhost:8080/api/auth

# hey -n 100000 -q 1000 -m GET -H "Authorization: Bearer `<jwt-token>`" http://localhost:8080/api/info