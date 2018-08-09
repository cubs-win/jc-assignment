# You have to check the return code from ListenandServe and only wait on the channel if it was shutdown. 
# A non-graceful shutdown in that case will still work properly
# Run a go linter
# Consider ptuting the (global) http server in the context so that you have no globals
# Add tests. For example, test that the hash function returns the expected value for "angryMonkey"
# Add test that spawns many connections and verify that all get a response before shutdown
# search for TODO in code & resolve all
# Add some kind of in/out/return documentation for all functions
# Remove all commented out code that you don't need
# Add production-ish logging and remove unused prints
