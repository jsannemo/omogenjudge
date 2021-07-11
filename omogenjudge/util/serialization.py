import dataclasses
import json
import typing


class IsDataclass(typing.Protocol):
    __dataclass_fields__: typing.Dict


class DataclassJsonEncoder(json.JSONEncoder):
    def default(self, o):
        if dataclasses.is_dataclass(o):
            return dataclasses.asdict(o)
        return json.JSONEncoder.default(self, o)


T = typing.TypeVar('T', bound=IsDataclass)


def dict_to_dataclass(cls: typing.ClassVar[IsDataclass], obj: dict):
    try:
        fields = {field.name: field.type for field in dataclasses.fields(cls)}
        return cls(**{key: dict_to_dataclass(fields[key], value) for key, value in obj})
    except:
        return obj


class DataclassJsonDecoder(json.JSONDecoder, typing.Generic[T]):
    def __init__(self, cls: typing.ClassVar[T], **kwargs):
        self._cls = cls
        super().__init__(**kwargs)

    def __call__(self, *args) -> T:
        obj = super().decode(*args)
        return dict_to_dataclass(self._cls, obj)
