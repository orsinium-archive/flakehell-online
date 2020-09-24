from pathlib import Path
from zipfile import ZipFile
from base64 import b64decode
from io import BytesIO

stream = BytesIO(b64decode(archive))  # noqa: F821
with ZipFile(stream) as zip:
    zip.extractall('.')
del archive  # noqa: F821
