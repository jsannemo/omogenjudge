from typing import Optional

from django import template

register = template.Library()


@register.filter
def format_duration_ms(time_ms: Optional[int]):
    if time_ms is None:
        return ""
    return "{:.2f}".format(time_ms / 1000)
