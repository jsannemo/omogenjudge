from django.utils.translation import get_language


def preferred_languages():
    lang = get_language()
    split = lang.split("-")
    langs = []
    for i in range(len(split)):
        langs.append("-".join(split[0:i + 1]))
    return langs[::-1] + ["en", "sv"]
