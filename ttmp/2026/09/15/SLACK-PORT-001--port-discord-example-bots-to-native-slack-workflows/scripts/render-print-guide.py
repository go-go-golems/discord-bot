#!/usr/bin/env python3
"""Create a printable copy of the frozen port guide with Graphviz figures."""
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "design-doc" / "01-port-inventory-and-implementation-plan.md"
OUT = ROOT / "various" / "print"
OUT.mkdir(parents=True, exist_ok=True)

flow_dot = r'''digraph G {
  graph [rankdir=TB, bgcolor="white", pad="0.2", nodesep="0.35", ranksep="0.45"];
  node [shape=box, style="rounded,filled", fillcolor="#eef4ff", color="#355c9a", fontname="DejaVu Sans", fontsize=14, margin="0.16,0.10"];
  edge [color="#355c9a", penwidth=1.3, arrowsize=0.7];
  Source [label="Discord source handlers"]; Matrix [label="Feature parity matrix"];
  Bots [label="Native Slack scripts"]; Host [label="Go-owned JavaScript runtime"];
  Domain [label="Detached domain operations"]; Live [label="Slack transport"];
  Recorder [label="Offline recorder"]; DB [label="SQLite module"];
  Slack [label="Slack workspace"]; Tests [label="Fixture assertions"];
  Reopen [label="Close and reopen tests"];
  Source -> Matrix -> Bots -> Host -> Domain;
  Domain -> Live -> Slack; Domain -> Recorder -> Tests;
  Host -> DB -> Reopen;
}'''

sequence_dot = r'''digraph G {
  graph [rankdir=LR, bgcolor="white", pad="0.2", nodesep="0.55", ranksep="0.35"];
  node [shape=box, style="rounded,filled", fillcolor="#eef4ff", color="#355c9a", fontname="DejaVu Sans", fontsize=13, margin="0.14,0.08"];
  edge [color="#355c9a", penwidth=1.2, arrowsize=0.65, fontname="DejaVu Sans", fontsize=10];
  U [label="User"]; S [label="Slack"]; G [label="Go host"]; J [label="Bot handler"]; D [label="SQLite"];
  U -> S [label="Click Edit"];
  S -> G [label="Block action"];
  G -> S [label="Empty ACK"];
  G -> J [label="Record/action context"];
  J -> D [label="Read record"];
  J -> G [label="Open modal with metadata"];
  G -> S [label="views.open"];
  U -> S [label="Submit edit"];
  S -> G [label="View submission"];
  G -> J [label="Fields and ACK capability"];
  J -> J [label="Validate actor and fields"];
  J -> G [label="Accept or field errors"];
  G -> S [label="Submission ACK"];
  J -> D [label="Commit valid edit"];
  J -> G [label="Update visible result"];
}'''

def render(name: str, dot: str) -> None:
    dot_path = OUT / f"{name}.dot"
    png_path = OUT / f"{name}.png"
    dot_path.write_text(dot)
    subprocess.run(["dot", "-Tpng", "-Gdpi=180", str(dot_path), "-o", str(png_path)], check=True)

render("architecture-flow", flow_dot)
render("review-edit-sequence", sequence_dot)

text = SOURCE.read_text()
text, flow_count = re.subn(
    r"```mermaid\nflowchart TD\n.*?\n```",
    "![Architecture flow](architecture-flow.png)\n\n*Architecture diagram rendered with Graphviz for print readability.*",
    text,
    count=1,
    flags=re.S,
)
text, sequence_count = re.subn(
    r"```mermaid\nsequenceDiagram\n.*?\n```",
    "![Review edit sequence](review-edit-sequence.png)\n\n*Review/edit sequence rendered with Graphviz for print readability.*",
    text,
    count=1,
    flags=re.S,
)
if flow_count != 1 or sequence_count != 1:
    raise SystemExit(f"expected two Mermaid blocks, replaced {flow_count} and {sequence_count}")
(OUT / SOURCE.name).write_text(text)
print(OUT / SOURCE.name)
