import enum
import typing

T = typing.TypeVar('T', bound=enum.Enum)


class EnumChoices(typing.Generic[T]):

    @classmethod
    def as_choices(cls) -> typing.List[typing.Tuple[T, str]]:
        return [(choice.value, choice.name) for choice in cls]
