#!/bin/python

import sys
import jwt
import json

from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import serialization

with open("./private.key", "r") as f:
    private_key = serialization.load_pem_private_key(
        f.read().encode("utf-8"), backend=default_backend(), password=None
    )
    print(jwt.encode({"iss": sys.argv[1]}, private_key, algorithm="RS512"))
