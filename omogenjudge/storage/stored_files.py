import base64
import hashlib
from typing import Union

from omogenjudge.storage.models import StoredFile


def insert_file(contents: Union[bytes, memoryview]) -> StoredFile:
    file_hash = hashlib.sha3_512()
    file_hash.update(contents)
    stored_file = StoredFile(
        file_hash=base64.b64encode(file_hash.digest()),
        file_contents=contents,
    )
    stored_file.save()
    return stored_file
