# test-signer

the multi-instance version requires redis, you may consider using the startRedis.sh script
the in-memory userRepository only supports a single instance

Run by supplying a jwt secret (required) as well as configuration flags as needed, e.g.:
JWT_SECRET=very_secure_secret go run cmd/test-signer/main.go --http-addr=:9090

I've decided to go with a light DDD aproach for the project structure, as this is what I'm argueably most experienced with,
I also dislike starting off with a TDD aproach when nothing substancial exists yet, as that leads, in my experience, to 
less well structured code.

The test.http file was used for manual testing using the vscode rest client extension, you might need to adjust the data (token, signatures)
if you want to use them as well. or just use the very secure jwt secret ("asd") I did and you can reuse the token. ;)

Aaaand that's four hours
