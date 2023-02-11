from datetime import timezone, timedelta
from typing import Optional

from crispy_forms.utils import render_crispy_form
from django.contrib.staticfiles.storage import staticfiles_storage
from django.template.context_processors import csrf
from django.urls import reverse
from jinja2 import Environment, pass_context

from omogenjudge.storage.models import Status, Verdict
from omogenjudge.util.contest_urls import reverse_contest
from omogenjudge.util.i18n import preferred_languages


def route_name(request):
    return request.resolver_match.url_name


def url_tag(view_name, *args, **kwargs):
    return reverse(view_name, args=args, kwargs=kwargs)


def contest_url_tag(view_name, *args, **kwargs):
    return reverse_contest(view_name, *args, **kwargs)


def format_duration_ms(time_ms: Optional[int], time_limit_ms: Optional[int] = None):
    if time_ms is None:
        return ""
    if time_limit_ms and time_ms > time_limit_ms:
        return ">{:.2f}".format(time_limit_ms / 1000)
    return "{:.2f}".format(time_ms / 1000)


def format_timedelta(delta: timedelta):
    seconds = int(delta.total_seconds())
    minutes = seconds // 60
    seconds %= 60

    hours = minutes // 60
    minutes %= 60

    days = hours // 24
    hours %= 24

    res = ""
    if days == 1:
        res += "1 day"
    elif days:
        res += f"{days} days"

    if (hours or minutes or seconds) or not res:
        if res:
            res += ", "
        res += f"{hours:02d}:{minutes:02d}:{seconds:02d}"
    return res


@pass_context
def render_crispy(ctx, form):
    context = csrf(ctx.get('request'))
    return render_crispy_form(form, context=context)


language_to_flag = {
    "da": "🇩🇰",
    "de": "🇩🇪",
    "en": "🇬🇧",
    "et": "🇪🇪",
    "is": "🇮🇸",
    "lt": "🇱🇹",
    "lv": "🇱🇻",
    "no": "🇳🇴",
    "pl": "🇵🇱",
    "ru": "🇷🇺",
    "sv": "🇸🇪",
}


def emoji_flag(language):
    return language_to_flag.get(language, language)


def environment(**options):
    env = Environment(**options)
    env.globals.update({
        "static": staticfiles_storage.url,
        "crispy": render_crispy,
        "url": url_tag,
        "contest_url": contest_url_tag,
        "route_name": route_name,
        "timezone": timezone,
        "Status": Status,
        "Verdict": Verdict,
        "format_duration_ms": format_duration_ms,
        "preferred_languages": preferred_languages,
        "emoji_flag": emoji_flag,
    })
    env.filters.update({
        "format_duration_ms": format_duration_ms,
        "format_timedelta": format_timedelta,
    })
    return env
