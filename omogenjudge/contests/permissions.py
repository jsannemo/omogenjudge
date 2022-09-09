from omogenjudge.util.request_global import current_contest


def can_view_contest_problems() -> bool:
    contest = current_contest()
    if not contest:
        return False
    if not contest.published:
        return False
    if contest.only_virtual_contest:
        return True
    else:
        return contest.has_started
