# UDP syslog shim

This is a lightweight UDP syslog server. It listens for syslog
messages on a socket and prints them to stderr as they're received.

It also spawns a child process using the given command line arguments
and passes signals down to it.

This is intended as a simple [haproxy](http://haproxy.1wt.eu/)
wrapper. Haproxy _requires_ a syslog server. There's no option to log
to a file, or even stderr.  Usually this is fine since you already
have [syslog-ng](http://www.balabit.com/network-security/syslog-ng) or
[rsyslog](http://www.rsyslog.com/) available on the system. But
sometimes you want to keep things simple, and this shim lets you
pretend that haproxy is logging to stderr.

## Credits

Thanks for [Michał Derkacz](https://github.com/ziutek) for the awesome
[syslog](https://github.com/ziutek/syslog server) Go package which
does all the heavy lifting.
