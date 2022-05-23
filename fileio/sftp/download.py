#! /usr/bin/python3

import os
import sys
import datetime
import pysftp

def main(argv):
    if len(argv) < 6:
        print("Arguments not given correctly (given = %d)" % len(argv), file=sys.stderr)
        sys.exit(1)

    hostport = argv[1]
    user = argv[2]
    password = argv[3]
    path = argv[4]
    localpath = argv[5]

    host = hostport
    port = 22

    if ":" in hostport:
        hostport_vars = hostport.split(":")
        host = hostport_vars[0].strip()
        port = int(hostport_vars[1].strip())

    if not host:
        print("SFTP HOST is not given", file=sys.stderr)
        sys.exit(1)

    if port <= 0:
        port = 22

    if not user:
        print("SFTP USER is not given", file=sys.stderr)
        sys.exit(1)

    if not password:
        print("SFTP PASSWORD is not given", file=sys.stderr)
        sys.exit(1)

    if not path:
        print("SFTP PATH is not given", file=sys.stderr)
        sys.exit(1)

    with pysftp.Connection(host=host, port=port, username=user, password=password) as sftp:

        print("downloading %s to %s" % (path, localpath))
        start = datetime.datetime.now()
        sftp.get(path, localpath)
        elapsed = datetime.datetime.now() - start
        print("elapsed - %s" % elapsed)


if __name__ == "__main__":
    main(sys.argv)