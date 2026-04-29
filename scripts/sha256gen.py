#!/usr/bin/env python3

# this is just for quick hash generation for testing in dev env
import hashlib
import sys


def main():
    string = sys.argv[1] if len(sys.argv) > 1 else ""
    hash = hashlib.sha256(string.encode()).hexdigest()
    print(f"SHA256 hash of '{string}': {hash}")


if __name__ == "__main__":
    main()
