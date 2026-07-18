import subprocess
import tempfile
import unittest
from pathlib import Path

from skfs.server import McpError, StructuredKnowledgeServer


class StructuredKnowledgeServerTests(unittest.TestCase):
    def test_list_nodes_and_read_markdown(self) -> None:
        with tempfile.TemporaryDirectory() as tempdir:
            root = Path(tempdir)
            (root / "notes").mkdir()
            (root / "notes" / "intro.md").write_text("# Intro\nhello world\n", encoding="utf-8")
            server = StructuredKnowledgeServer(root=root)

            nodes = server.list_nodes({"path": ".", "depth": 2})
            self.assertIn({"path": "notes", "type": "directory"}, nodes)
            self.assertIn({"path": "notes/intro.md", "type": "file"}, nodes)

            content = server.read_markdown({"path": "notes/intro.md"})
            self.assertIn("hello world", content)

    def test_path_escape_is_rejected(self) -> None:
        with tempfile.TemporaryDirectory() as tempdir:
            server = StructuredKnowledgeServer(root=Path(tempdir))
            with self.assertRaises(McpError):
                server.read_markdown({"path": "../outside.md"})

    def test_git_show_file(self) -> None:
        with tempfile.TemporaryDirectory() as tempdir:
            root = Path(tempdir)
            subprocess.run(["git", "-C", str(root), "init"], check=True, capture_output=True)
            subprocess.run(
                ["git", "-C", str(root), "config", "user.email", "test@example.com"],
                check=True,
                capture_output=True,
            )
            subprocess.run(
                ["git", "-C", str(root), "config", "user.name", "Test User"],
                check=True,
                capture_output=True,
            )
            (root / "doc.md").write_text("hello from git\n", encoding="utf-8")
            subprocess.run(["git", "-C", str(root), "add", "doc.md"], check=True, capture_output=True)
            subprocess.run(
                ["git", "-C", str(root), "commit", "-m", "add doc"],
                check=True,
                capture_output=True,
            )

            server = StructuredKnowledgeServer(root=root)
            content = server.git_show_file({"path": "doc.md", "ref": "HEAD"})
            self.assertEqual("hello from git\n", content)


if __name__ == "__main__":
    unittest.main()

