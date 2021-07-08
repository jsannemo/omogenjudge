from django.db import models


class StoredFile(models.Model):
    file_hash = models.CharField(primary_key=True, max_length=256)
    file_contents = models.BinaryField()

    class Meta:
        db_table = 'stored_file'
