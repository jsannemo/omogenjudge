from django import template
from django.urls import reverse
from django.utils.html import format_html

register = template.Library()


@register.simple_tag(takes_context=True)
def navitem(context: dict, title: str, route: str, **kwargs):
    return format_html(
        r"""<a class="nav-link {}" href="{}">{}</a>""",
        "active" if context['request'].resolver_match.url_name == route else "",
        reverse(route, kwargs=kwargs),
        title,
    )
