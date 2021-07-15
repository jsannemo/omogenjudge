import enum

from omogenjudge.util.enums import EnumChoices


class Language(EnumChoices['Language'], enum.Enum):
    CPP = 'cpp'
    PYTHON3 = 'python3'

    def display(self):
        return LANGUAGE_NAMES[self]


LANGUAGE_NAMES = {
    Language.CPP: 'C++',
    Language.PYTHON3: 'Python 3'
}
