#! /usr/bin/python3

import os
import sys
import getpass

from irods.session import iRODSSession
from irods.exception import DataObjectDoesNotExist

# irods_mkdir.py hostname:port zone /parentdir/targetdir
# iRODS Username and Password are passed via STDIN

def main(argv):
    if len(argv) < 5:
        print("Arguments not given correctly (given = %d)" % len(argv), file=sys.stderr)
        sys.exit(1)

    hostport = argv[1]
    user = argv[2]
    password = argv[3]
    zone = argv[4]
    localpath = argv[5]
    path = argv[6]

    host = hostport
    port = 1247

    if ":" in hostport:
        hostport_vars = hostport.split(":")
        host = hostport_vars[0].strip()
        port = int(hostport_vars[1].strip())

    if not host:
        print("iRODS HOST is not given", file=sys.stderr)
        sys.exit(1)

    if not zone:
        print("iRODS ZONE is not given", file=sys.stderr)
        sys.exit(1)

    if port <= 0:
        port = 1247

    if not user:
        print("iRODS USER is not given", file=sys.stderr)
        sys.exit(1)

    if not password:
        print("iRODS PASSWORD is not given", file=sys.stderr)
        sys.exit(1)

    if not path:
        print("iRODS PATH is not given", file=sys.stderr)
        sys.exit(1)

    zonepath = path
    if not path.startswith("/" + zone + "/"):
        zonepath = "/" + zone + "/" + path.lstrip("/")

    with iRODSSession(host=host, port=port, user=user, password=password, zone=zone, client_user=user) as session:
        try:
            session.data_objects.put(localpath, zonepath)
        except DataObjectDoesNotExist:
            print("Could not upload a file %s" % localpath, file=sys.stderr)
            sys.exit(1)


if __name__ == "__main__":
    main(sys.argv)