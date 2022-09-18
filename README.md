HatParty app:

https://gist.github.com/larrycinnabar/d5f1c7a008af8068b65095d876f08e3a

Missing:
* automated test
* swagger docs
* dockerfile
* proper error handling
* ci/cd pipeline

Config.yml contains all the keys for configuration

Use /init endpoint for populating the initial hats in mongo
Use /start?hatsRequested= for creating a party with specific number of hats
Use /stop?partyId= for stoping a party and releasing hats