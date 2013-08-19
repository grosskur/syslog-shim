"""
Fake server to test signal handling
"""
import signal
import sys
import time


def handler_term(signum, frame):
    print >> sys.stderr, 'fake_server: received signal', signum
    time.sleep(5)
    sys.exit(2)


def handler_test(signum, frame):
    print >> sys.stderr, 'fake_server: received signal', signum


signal.signal(signal.SIGINT, handler_term)
signal.signal(signal.SIGTERM, handler_term)
signal.signal(signal.SIGQUIT, handler_term)
signal.signal(signal.SIGHUP, handler_test)
signal.signal(signal.SIGUSR1, handler_test)
signal.signal(signal.SIGUSR2, handler_test)

try:
    print >> sys.stderr, 'fake_server: starting'
    while True:
        time.sleep(10)
except KeyboardInterrupt:
    print >> sys.stderr, 'fake_server: interrupted'
