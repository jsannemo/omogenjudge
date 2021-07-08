import enum
import typing

from django.db import connection, models


class TextField(models.TextField):

    def __init__(self, *args, **kwargs):
        # Always set None as the default; otherwise this may be the empty string, which we don't want
        if 'default' not in kwargs:
            kwargs['default'] = None
        super().__init__(*args, **kwargs)

    def deconstruct(self):
        name, path, args, kwargs = super().deconstruct()
        if kwargs['default'] is None:
            del kwargs['default']
        return name, path, args, kwargs


class StrEnum(str, enum.Enum):
    def __str__(self) -> str:
        return self.value


E = typing.TypeVar('E', bound=enum.Enum)


class EnumField(models.Field, typing.Generic[E]):
    def __init__(self, *args, enum_type: typing.Type[E], **kwargs):
        if 'default' not in kwargs:
            kwargs['default'] = None
        super().__init__(*args, **kwargs)
        self._enum_type: typing.Type[E] = enum_type
        self._values_to_enum: typing.Dict[str, E] = {
            e.value: e for e in enum_type
        }
        self.choices = [
            (e, e.value)
            for e in enum_type
        ]

    def deconstruct(self):
        name, path, args, kwargs = super().deconstruct()
        kwargs['enum_type'] = self._enum_type
        if kwargs['default'] is None:
            del kwargs['default']
        return name, path, args, kwargs

    def to_python(self, value: typing.Union[E, str, None]) -> typing.Optional[E]:
        if isinstance(value, self._enum_type):
            enum_value: E = value
            return enum_value
        if value is None:
            return value
        assert isinstance(value, str)
        return self._values_to_enum[value]

    def from_db_value(self, value: typing.Optional[str], _expression, _connection) -> typing.Optional[E]:
        if value is None:
            return value
        assert isinstance(value, str)
        return self._values_to_enum[value]

    def get_prep_value(self, value: typing.Optional[E]) -> typing.Optional[str]:
        if value is None:
            return None
        if not isinstance(value, self._enum_type):
            raise ValueError(f"Expected field value of enum {self._enum_type}; was {type(value)}")
        return value.value

    def value_to_string(self, obj) -> str:
        value = self.value_from_object(obj)
        str_value = self.get_prep_value(value)
        assert isinstance(str_value, str)
        return str_value

    def get_internal_type(self) -> str:
        return 'TextField'


# https://djangosnippets.org/snippets/10474/
class PrefetchIDMixin(object):

    def prefetch_id(self):
        cursor = connection.cursor()
        table_name = self._meta.db_table
        cursor.execute(
            "SELECT nextval('{0}_{1}_seq'::regclass)".format(
                table_name,
                self._meta.pk.name,
            )
        )
        row = cursor.fetchone()
        cursor.close()
        self.pk = row[0]
