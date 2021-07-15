from django.db import migrations


class Migration(migrations.Migration):
    initial = True

    dependencies = [
        ('storage', '0001_initial'),
    ]

    operations = [
        migrations.RunSQL(
            """
            CREATE FUNCTION notify_run() RETURNS TRIGGER AS $$
            BEGIN
                PERFORM pg_notify('new_run', (NEW.submission_run_id)::text);
                RETURN NULL;
            END;
            $$ LANGUAGE plpgsql;

            CREATE TRIGGER "new_run"
                AFTER INSERT ON submission_run
                FOR EACH ROW EXECUTE PROCEDURE notify_run();
            """,
            """
            DROP TRIGGER "new_run"
            DROP FUNCTION notify_run;
            """
        )
    ]
