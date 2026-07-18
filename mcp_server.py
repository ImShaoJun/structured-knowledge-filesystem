#!/usr/bin/env python3
from __future__ import annotations

import os
from pathlib import Path

from skfs.server import run_stdio_server


def main() -> int:
    root = Path(os.environ.get("SKFS_ROOT", os.getcwd()))
    return run_stdio_server(root=root)


if __name__ == "__main__":
    raise SystemExit(main())

