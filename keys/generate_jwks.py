#!/bin/python

import jwt
import json

from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import serialization

with open("private.key.pub", "r") as f:
    public_key = serialization.load_pem_public_key(
        f.read().encode("utf-8"), backend=default_backend()
    )
    print(json.dumps(jwt.algorithms.RSAAlgorithm.to_jwk(public_key)))
