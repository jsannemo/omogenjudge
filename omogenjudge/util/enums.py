import enum
import typing

T = typing.TypeVar('T', bound=enum.Enum)


class EnumChoices(typing.Generic[T]):

    @classmethod
    def as_choices(cls) -> list[tuple[T, str]]:
        return [(choice, choice.display) for choice in cls]

    def display(self):
        return self.name

    def __str__(self):
        return self.value
