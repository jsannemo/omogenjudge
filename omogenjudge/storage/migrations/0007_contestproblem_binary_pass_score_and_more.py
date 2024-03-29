# Generated by Django 4.1.6 on 2023-03-19 09:31

from django.db import migrations, models
import django.db.models.deletion
import omogenjudge.util.django_fields


class Migration(migrations.Migration):

    dependencies = [
        ('storage', '0006_team_practice'),
    ]

    operations = [
        migrations.AddField(
            model_name='contestproblem',
            name='binary_pass_score',
            field=models.IntegerField(blank=True, default=None, null=True),
        ),
        migrations.AlterField(
            model_name='problemstatementfile',
            name='problem',
            field=models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='statement_files', to='storage.problem'),
        ),
        migrations.AlterField(
            model_name='team',
            name='team_name',
            field=omogenjudge.util.django_fields.TextField(blank=True, null=True),
        ),
    ]
