## nitroshare2

This repository contains the source code for the (completely rewritten) Nitroshare client.

Development goals:

- **Complete rewrite in Go** — While C++/Qt was a very good choice for a cross-platform GUI application, the complexity of using a language and framework with no built-in build / installation system wasted an enormouse amount of development time. Go offers built-in dependency management / installation, unit testing, and incredibly simple deployment.

- **Support for running in a terminal** — The existing client was written in a way that required a graphical subsystem on all major platforms. While there was work underway to separate the core protocol code from the user interface, a complete rewrite in Go will solve this problem in a much more straightforward way.

- **Vastly improved mDNS-based discovery** — One of the most frequent source for issues being filed against the client was the IP broadcast / mDNS code. The new client will use an existing (and better tested) mDNS package, hopefully eliminating a lot of the issues.

- **Extensible with plugins written in JavaScript** — Plugins can extend the client (and protocol itself!) with new features. A simple JavaScript runtime will be included, allowing developers to use a widely popular language.

- **Well-documented and comprehensive local API** — Having nearly all of the protocol code behind a local API enables the application to be very modular and therefore easier to port across platforms (including Android/iOS).

- **More service integrations (email, MQTT, etc.)** — Integrating Nitroshare with other services (like MQTT for Home Assistant, etc.) will open up a lot of new possibilities and simplify information transfer between devices.
