from pathlib import Path
import re

ticket = Path(__file__).resolve().parents[1]
source = ticket / "design-doc/01-slack-surfaces-and-ui-dsl-intern-guide.md"
reading = ticket / "various/print/01-slack-surfaces-and-ui-dsl-intern-guide.md"
text = source.read_text()

replacements = [
    '![Message flow from bot JavaScript through Slack UI builders and the Go host](assets/message-flow.png)',
    '![Modal submission sequence showing validation and acknowledgment paths](assets/modal-sequence.png)',
]
text, count = re.subn(r"```mermaid\nflowchart TD\n.*?```", replacements[0], text, count=1, flags=re.S)
text, count2 = re.subn(r"```mermaid\nsequenceDiagram\n.*?```", replacements[1], text, count=1, flags=re.S)
if count != 1 or count2 != 1:
    raise SystemExit(f"expected two Mermaid replacements, got flow={count}, sequence={count2}")
text = text.replace("../sources/", "../../sources/")
reading.write_text(text)
print(f"refreshed {reading}")
