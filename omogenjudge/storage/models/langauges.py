import enum

from omogenjudge.util.django_fields import StrEnum


class Language(StrEnum):
    CPP = 'cpp'
    PYTHON3 = 'python3'
    RUBY = 'ruby'
    RUST = 'rust'
    JAVA = 'java'
    CSHARP = 'csharp'

    # TODO: unsupported
    # C = 'c'
    # GO = 'go'
    # JS = 'js'

    def display(self):
        return LANGUAGE_NAMES[self]

    @classmethod
    def as_choices(cls):
        return [(v, v.display()) for v in Language]


LANGUAGE_NAMES = {
    Language.CPP: 'C++',
    Language.PYTHON3: 'Python 3',
    # Language.C: 'C',
    Language.CSHARP: 'C#',
    Language.JAVA: 'Java',
    Language.RUBY: 'Ruby',
    Language.RUST: 'Rust',
    # Language.GO: 'go',
    # Language.JS: 'js',
}
