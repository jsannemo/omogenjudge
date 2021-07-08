import unicodedata


def normalize_username(username: str) -> str:
    return unicodedata.normalize('NFKC', username)
