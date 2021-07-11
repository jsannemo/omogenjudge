def ask_yes_or_no(query, default_response: bool) -> bool:
    while True:
        response = input(query)
        if response == "":
            return default_response
        if response.lower()[0] == 'y':
            return True
        if response.lower()[0] == 'n':
            return False
