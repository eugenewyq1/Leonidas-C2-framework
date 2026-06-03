#!/usr/bin/env bash
# Remove Cursor from GitHub contributors by recreating the repo.
# Closed Dependabot PRs still reference old commits with Co-authored-by: Cursor;
# deleting and recreating the repository is the only reliable fix on GitHub.
set -euo pipefail

REPO="${1:-eugenewyq1/Leonidas-C2-framework}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if ! command -v gh >/dev/null 2>&1; then
  echo "Install GitHub CLI: brew install gh && gh auth login"
  exit 1
fi

echo "==> Single orphan commit (no Cursor trailer, no git commit hook)"
python3 <<'PY'
import subprocess
import time

subprocess.run(["git", "rm", "-f", "--ignore-unmatch", ".github/dependabot.yml"], check=False)
subprocess.run(["git", "add", "-A"], check=True)
tree = subprocess.check_output(["git", "write-tree"], text=True).strip()
author = "Wong Yong Quan Eugene <87485670+eugenewyq1@users.noreply.github.com>"
tz = "+0800"
now = int(time.time())
body = f"tree {tree}\nauthor {author} {now} {tz}\ncommitter {author} {now} {tz}\n\nInitial release\n"
new_commit = subprocess.check_output(
    ["git", "hash-object", "-t", "commit", "-w", "--stdin"],
    input=body,
    text=True,
).strip()
subprocess.run(["git", "update-ref", "refs/heads/main", new_commit], check=True)
print(new_commit)
PY

echo "==> Delete GitHub repo (removes old PRs/commits from contributor index)"
gh repo delete "$REPO" --yes

echo "==> Recreate empty repo and push clean history"
gh repo create "${REPO#*/}" --public --description "A modified fork of the Sliver C2 framework"
git remote remove origin 2>/dev/null || true
git remote add origin "https://github.com/${REPO}.git"
git push -u origin main --force

echo "Done. Contributors should show only eugenewyq1 within a few minutes."
