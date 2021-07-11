from django.db import connection, models


class TextField(models.TextField):

    def __init__(self, *args, **kwargs):
        if 'default' not in kwargs:
            kwargs['default'] = None
        super().__init__(*args, **kwargs)


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
