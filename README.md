# Overview


Rod is a High-level Devtools driver directly based on [DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/).
It's designed for web automation and scraping. Rod also tries to expose low-level interfaces to users, so that whenever a function is missing users can easily send control requests to the browser directly.

Wayang is created to be a controller of Rod, making it possible to write programs with Rod in a more language neutral way. While it does not cover all of Rod's API at the moment, we plan to make it feature complete.
