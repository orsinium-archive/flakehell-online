from pathlib import Path
from textwrap import dedent

from flake8.main.application import Application
from flake8.formatting.default import Default

path = Path("code.py")
path.write_text(dedent(text))  # noqa: F821


class Formatter(Default):
    def after_init(self):
        self._out = []
        return super().after_init()

    def _write(self, output):
        self._out.append(output)


class App(Application):
    def make_formatter(self, formatter_class=None):
        self.formatter = Formatter(self.options)


app = App()
code = 0
try:
    app.run(["code.py"])
    app.exit()
except SystemExit as err:
    code = int(err.args[0])
