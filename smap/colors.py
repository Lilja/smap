
class ColorCode:
    BRIGHT_WHITE = 1
    GREY = 2
    BLACK = 30
    RED = 31
    GREEN = 32
    YELLOW = 33
    DARK_BLUE = 34
    PURPLE = 35
    CYAN = 36
    WHITE = 37
    LIGHT_GREY = 90
    ORANGE = 91
    LIGHT_GREEN = 92
    LIGHT_YELLOW = 93
    BLUE = 94
    PINK = 95
    SEMI_BRIGHT_WHITE = 97


class ColorStyle:
    NO_EFFECT = 0
    BOLD = 1
    WEAK = 2
    CURSIVE = 3
    UNDERLINE = 4
    REVERSED = 9
    STRIKETHROUGH = 9


class ColorInterface:
    ESCAPE = "\033["
    END = ESCAPE + "0m"


color_codes = dict(
        (x, f"{ColorInterface.ESCAPE}0;{getattr(ColorCode, x)}m")
        for x in filter(lambda x: x[0] != "_", dir(ColorCode))
    )
color_codes.update(
    dict(
        (x, f"{ColorInterface.ESCAPE}0;{getattr(ColorStyle, x)}m")
        for x in filter(lambda x: x[0] != "_", dir(ColorStyle))
    )
)
color_codes['NC'] = ColorInterface.END


def colorize(_str):
    return _str.format(**color_codes)
