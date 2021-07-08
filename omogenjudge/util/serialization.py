import dataclasses
import json
import typing


class IsDataclassClass(typing.Protocol):
    __dataclass_fields__: dict


class IsDictable(typing.Protocol):
    __dict__: dict


class DataclassJsonEncoder(json.JSONEncoder):
    def default(self, o):
        if dataclasses.is_dataclass(o):
            return dataclasses.asdict(o)
        return json.JSONEncoder.default(self, o)


T = typing.TypeVar('T', bound=IsDataclassClass)


def dict_to_dataclass(cls: typing.Type[T], obj: dict) -> T:
    fields = {field.name: field.type for field in dataclasses.fields(cls)}
    return cls(**{key: dict_to_dataclass(fields[key], value) if dataclasses.is_dataclass(fields[key]) else value for
                  key, value in obj})


class DataclassJsonDecoder(json.JSONDecoder, typing.Generic[T]):
    def __init__(self, cls: typing.Type[T], **kwargs):
        self._cls = cls
        super().__init__(**kwargs)

    def __call__(self, *args) -> T:
        obj = super().decode(*args)
        return dict_to_dataclass(self._cls, obj)
